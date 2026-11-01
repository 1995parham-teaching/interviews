// Package node puts the pieces of one Relay node together: the shared backend
// list, the health checker, the ring, the proxy, the admin API and gossip.
package node

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"relay/cluster"
	"relay/health"
	"relay/proxy"
	"relay/ring"
)

// Config is everything a node needs to start.
type Config struct {
	ID             string
	Peers          []string // other nodes, host:port
	Backends       []string // the seed list, the same on every node
	CheckInterval  time.Duration
	GossipInterval time.Duration
}

// Node is one load balancer in the cluster.
type Node struct {
	cfg     Config
	state   *cluster.State
	checker *health.Checker
	proxy   *proxy.Proxy
	client  *http.Client

	mu   sync.RWMutex
	ring *ring.Ring
}

// New builds a node. Nothing runs until Run is called.
func New(cfg Config) *Node {
	n := &Node{
		cfg:     cfg,
		state:   cluster.NewState(cfg.Backends),
		checker: health.New(),
		client:  &http.Client{Timeout: 2 * time.Second},
	}
	n.proxy = proxy.New(cfg.ID, n.pick, n.checker.Healthy)
	n.rebuild()
	n.state.OnChange(n.rebuild)
	return n
}

// State is the node's copy of the backend list.
func (n *Node) State() *cluster.State { return n.state }

// Checker is the node's health checker.
func (n *Node) Checker() *health.Checker { return n.checker }

// Run checks backends and gossips with peers until ctx is done.
func (n *Node) Run(ctx context.Context) {
	go n.checker.Run(ctx, n.cfg.CheckInterval, n.state.Backends)
	t := time.NewTicker(n.cfg.GossipInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			n.GossipOnce(ctx)
		}
	}
}

// Handler serves client traffic on / and the node's own API under /relay/.
func (n *Node) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /relay/status", n.status)
	mux.HandleFunc("PUT /relay/backends/{backend}", n.putBackend)
	mux.HandleFunc("DELETE /relay/backends/{backend}", n.deleteBackend)
	mux.HandleFunc("POST /relay/gossip", n.gossip)
	mux.Handle("/", n.proxy)
	return mux
}

func (n *Node) rebuild() {
	r := ring.New(n.state.ActiveBackends())
	n.mu.Lock()
	n.ring = r
	n.mu.Unlock()
}

func (n *Node) pick(key string, usable func(string) bool) (string, bool) {
	n.mu.RLock()
	r := n.ring
	n.mu.RUnlock()
	return r.Lookup(key, usable)
}

// ---------------------------------------------------------------- admin API

type backendStatus struct {
	cluster.Entry
	Healthy bool `json:"healthy"`
}

type statusReply struct {
	Node     string          `json:"node"`
	Peers    []string        `json:"peers"`
	Ring     []string        `json:"ring"`
	Backends []backendStatus `json:"backends"`
}

func (n *Node) status(w http.ResponseWriter, _ *http.Request) {
	n.mu.RLock()
	onRing := n.ring.Backends()
	n.mu.RUnlock()

	reply := statusReply{Node: n.cfg.ID, Peers: n.cfg.Peers, Ring: onRing}
	for _, e := range n.state.Entries() {
		reply.Backends = append(reply.Backends, backendStatus{Entry: e, Healthy: n.checker.Healthy(e.Backend)})
	}
	writeJSON(w, http.StatusOK, reply)
}

// PUT /relay/backends/{backend} with {"status": "active"} or {"status": "draining"}
// adds the backend or changes its status.
func (n *Node) putBackend(w http.ResponseWriter, r *http.Request) {
	backend := r.PathValue("backend")
	if !strings.Contains(backend, ":") {
		http.Error(w, "backend must be host:port", http.StatusBadRequest)
		return
	}
	var req struct {
		Status cluster.Status `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad body: "+err.Error(), http.StatusBadRequest)
		return
	}
	e, err := n.state.Set(backend, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("backend set", "node", n.cfg.ID, "backend", e.Backend, "status", e.Status, "version", e.Version)
	writeJSON(w, http.StatusOK, e)
}

func (n *Node) deleteBackend(w http.ResponseWriter, r *http.Request) {
	backend := r.PathValue("backend")
	if !n.state.Remove(backend) {
		http.Error(w, "unknown backend "+backend, http.StatusNotFound)
		return
	}
	slog.Info("backend removed", "node", n.cfg.ID, "backend", backend)
	w.WriteHeader(http.StatusNoContent)
}

// ---------------------------------------------------------------- gossip

// gossip is the receiving side: merge what the peer sent, answer with ours.
func (n *Node) gossip(w http.ResponseWriter, r *http.Request) {
	var remote []cluster.Entry
	if err := json.NewDecoder(r.Body).Decode(&remote); err != nil {
		http.Error(w, "bad gossip: "+err.Error(), http.StatusBadRequest)
		return
	}
	if c := n.state.Merge(remote); c > 0 {
		slog.Info("gossip merged", "node", n.cfg.ID, "changed", c)
	}
	writeJSON(w, http.StatusOK, n.state.Entries())
}

// GossipOnce sends the list to every peer and merges each answer.
func (n *Node) GossipOnce(ctx context.Context) {
	var wg sync.WaitGroup
	for _, peer := range n.cfg.Peers {
		wg.Go(func() {
			if err := n.exchange(ctx, peer); err != nil {
				slog.Warn("gossip failed", "node", n.cfg.ID, "peer", peer, "err", err)
			}
		})
	}
	wg.Wait()
}

func (n *Node) exchange(ctx context.Context, peer string) error {
	body, err := json.Marshal(n.state.Entries())
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+peer+"/relay/gossip", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("peer answered " + resp.Status)
	}
	var remote []cluster.Entry
	if err := json.NewDecoder(resp.Body).Decode(&remote); err != nil {
		return err
	}
	if c := n.state.Merge(remote); c > 0 {
		slog.Info("gossip merged", "node", n.cfg.ID, "peer", peer, "changed", c)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
