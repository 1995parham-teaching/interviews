package cluster

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"

	"relay/config"
)

// SecretHeader carries the cluster secret between nodes.
const SecretHeader = "X-Relay-Secret"

// maxGossip bounds the size of a list a peer may send.
const maxGossip = 4 << 20

// Gossiper exchanges the list with every peer, every interval.
type Gossiper struct {
	State   *State
	Peers   []string
	Timeout time.Duration
	Secret  string

	mu    sync.Mutex
	peers map[string]*PeerStatus
}

// PeerStatus is what a node knows about one peer.
type PeerStatus struct {
	Peer      string    `json:"peer"`
	LastOK    time.Time `json:"last_ok"`
	LastError string    `json:"last_error,omitempty"`
	Failures  int       `json:"failures"`
}

// Run gossips every interval until ctx is done.
func (g *Gossiper) Run(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			g.Round(ctx)
			if n := g.State.Sweep(); n > 0 {
				slog.Info("tombstones swept", "node", g.State.id, "count", n)
			}
		}
	}
}

// Round exchanges the list with every peer once, in parallel.
func (g *Gossiper) Round(ctx context.Context) {
	var wg sync.WaitGroup
	for _, peer := range g.Peers {
		wg.Go(func() {
			err := g.exchange(ctx, peer)
			g.record(peer, err)
		})
	}
	wg.Wait()
}

// Statuses returns what this node knows about its peers.
func (g *Gossiper) Statuses() []PeerStatus {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]PeerStatus, 0, len(g.Peers))
	for _, p := range g.Peers {
		if st, ok := g.peers[p]; ok {
			out = append(out, *st)
		} else {
			out = append(out, PeerStatus{Peer: p})
		}
	}
	return out
}

func (g *Gossiper) record(peer string, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.peers == nil {
		g.peers = map[string]*PeerStatus{}
	}
	st, ok := g.peers[peer]
	if !ok {
		st = &PeerStatus{Peer: peer}
		g.peers[peer] = st
	}
	if err == nil {
		st.LastOK, st.LastError, st.Failures = time.Now(), "", 0
		return
	}
	st.LastError = err.Error()
	st.Failures++
	if st.Failures == 1 || st.Failures%30 == 0 {
		slog.Warn("gossip failed", "node", g.State.id, "peer", peer, "failures", st.Failures, "err", err)
	}
}

// client returns an HTTP client for one exchange. The dial and the whole
// exchange each get the gossip timeout, so a peer that is down or stuck costs
// at most one timeout per round and never holds up the next round.
func (g *Gossiper) client() *http.Client {
	return &http.Client{
		Timeout: g.Timeout,
		Transport: &http.Transport{
			DialContext:         (&net.Dialer{Timeout: g.Timeout}).DialContext,
			TLSHandshakeTimeout: g.Timeout,
		},
	}
}

func (g *Gossiper) exchange(ctx context.Context, peer string) error {
	body, err := json.Marshal(g.State.Entries())
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+peer+"/relay/gossip", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if g.Secret != "" {
		req.Header.Set(SecretHeader, g.Secret)
	}
	resp, err := g.client().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("peer answered %s: %s", resp.Status, bytes.TrimSpace(msg))
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, maxGossip))
	if err != nil {
		return err
	}
	var remote []Entry
	if err := json.Unmarshal(b, &remote); err != nil {
		return err
	}
	remote = slices.DeleteFunc(remote, func(e Entry) bool { return !valid(e) })
	if n := g.State.Merge(remote); n > 0 {
		slog.Info("gossip merged", "node", g.State.id, "peer", peer, "changed", n)
	}
	return nil
}

// Handler is the receiving side of an exchange: merge what the peer sent and
// answer with this node's list.
func (g *Gossiper) Handler(w http.ResponseWriter, r *http.Request) {
	if g.Secret != "" && subtle.ConstantTimeCompare([]byte(r.Header.Get(SecretHeader)), []byte(g.Secret)) != 1 {
		http.Error(w, "bad cluster secret", http.StatusForbidden)
		return
	}
	var remote []Entry
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxGossip)).Decode(&remote); err != nil {
		http.Error(w, "bad gossip: "+err.Error(), http.StatusBadRequest)
		return
	}
	remote = slices.DeleteFunc(remote, func(e Entry) bool { return !valid(e) })
	if n := g.State.Merge(remote); n > 0 {
		slog.Info("gossip merged", "node", g.State.id, "from", r.RemoteAddr, "changed", n)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g.State.Entries())
}

// valid drops entries no node could have produced.
func valid(e Entry) bool {
	if _, _, err := net.SplitHostPort(e.Backend); err != nil {
		return false
	}
	switch e.Status {
	case Active, Draining:
		return e.Weight >= 0 && e.Weight <= config.MaxWeight
	case Removed:
		return true
	}
	return false
}
