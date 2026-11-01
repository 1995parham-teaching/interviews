// Package node puts the pieces of one Relay node together: the shared backend
// list and its gossip, the health checker, the ring, session pins, the proxy,
// metrics and the admin API.
package node

import (
	"context"
	"log/slog"
	"runtime"
	"strconv"
	"sync"
	"time"

	"relay/affinity"
	"relay/cluster"
	"relay/config"
	"relay/health"
	"relay/metrics"
	"relay/proxy"
	"relay/ring"
)

// Node is one load balancer in the cluster.
type Node struct {
	cfg      config.Config
	state    *cluster.State
	gossiper *cluster.Gossiper
	checker  *health.Checker
	pins     *affinity.Pins
	proxy    *proxy.Proxy
	metrics  *metrics.Registry
	started  time.Time

	requests *metrics.Counter
	retries  *metrics.Counter
	flips    *metrics.Counter

	buildMu sync.Mutex // serialises rebuilds, so the last one wins
	mu      sync.RWMutex
	ring    *ring.Ring
}

// New builds a node from a validated configuration. Nothing runs until Run.
func New(cfg config.Config) (*Node, error) {
	trusted, err := cfg.TrustedPrefixes()
	if err != nil {
		return nil, err
	}
	n := &Node{
		cfg:     cfg,
		state:   cluster.NewState(cfg.ID, cfg.Backends, cfg.TombstoneTTL.D()),
		checker: health.New(cfg.CheckTimeout.D(), cfg.Fall, cfg.Rise),
		pins:    affinity.New(cfg.PinTTL.D()),
		metrics: &metrics.Registry{},
		started: time.Now(),
	}
	n.gossiper = &cluster.Gossiper{
		State:   n.state,
		Peers:   cfg.Peers,
		Timeout: cfg.GossipTimeout.D(),
		Secret:  cfg.ClusterSecret,
	}
	n.proxy = &proxy.Proxy{
		Node:       cfg.ID,
		Pick:       n.pick,
		Healthy:    n.checker.Healthy,
		Report:     n.report,
		Client:     proxy.NewClient(),
		Retries:    cfg.Retries,
		TryTimeout: cfg.TryTimeout.D(),
		MaxBody:    cfg.MaxBody,
		Trusted:    trusted,
	}

	n.requests = n.metrics.NewCounter("relay_attempts_total", "Attempts sent to backends, by backend and status (0: no response).", "backend", "status")
	n.retries = n.metrics.NewCounter("relay_retries_total", "Attempts that were followed by another attempt, by backend.", "backend")
	n.flips = n.metrics.NewCounter("relay_health_flips_total", "Backends going down or up, as seen by this node.", "backend", "to")
	n.metrics.NewGauge("relay_pins", "Session pins held by this node.", func() float64 { return float64(n.pins.Len()) })
	n.metrics.NewGauge("relay_clock", "This node's logical clock.", func() float64 { return float64(n.state.Clock()) })
	n.metrics.NewGauge("go_goroutines", "Goroutines in this process.", func() float64 { return float64(runtime.NumGoroutine()) })

	n.rebuild()
	n.state.OnChange(n.rebuild)
	n.checker.OnFlip(func(b string, up bool) {
		to := "down"
		if up {
			to = "up"
		}
		n.flips.Inc(b, to)
		slog.Info("health flipped", "node", cfg.ID, "backend", b, "to", to)
	})
	return n, nil
}

// State is the node's copy of the backend list.
func (n *Node) State() *cluster.State { return n.state }

// Checker is the node's health checker.
func (n *Node) Checker() *health.Checker { return n.checker }

// Gossiper is the node's gossip.
func (n *Node) Gossiper() *cluster.Gossiper { return n.gossiper }

// Run checks backends, gossips with peers and sweeps pins until ctx is done.
func (n *Node) Run(ctx context.Context) {
	go n.checker.Run(ctx, n.cfg.CheckInterval.D(), n.state.Known)
	go n.gossiper.Run(ctx, n.cfg.GossipInterval.D())
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			n.pins.Sweep()
		}
	}
}

// rebuild makes a new ring from the active backends and drops the pins of
// backends that left it.
func (n *Node) rebuild() {
	n.buildMu.Lock()
	defer n.buildMu.Unlock()

	var members []ring.Member
	for _, e := range n.state.Active() {
		members = append(members, ring.Member{Backend: e.Backend, Weight: e.Weight})
	}
	next := ring.New(members)

	n.mu.Lock()
	prev := n.ring
	n.ring = next
	n.mu.Unlock()

	if prev == nil {
		return
	}
	still := map[string]bool{}
	for _, m := range next.Members() {
		still[m.Backend] = true
	}
	for _, m := range prev.Members() {
		if !still[m.Backend] {
			if c := n.pins.Drop(m.Backend); c > 0 {
				slog.Info("pins dropped", "node", n.cfg.ID, "backend", m.Backend, "count", c)
			}
		}
	}
}

func (n *Node) currentRing() *ring.Ring {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.ring
}

// pick keeps a key on its pinned backend while that backend is usable, and
// otherwise asks the ring and pins the answer.
func (n *Node) pick(key string, usable func(string) bool) (string, bool) {
	if b, ok := n.pins.Get(key); ok && usable(b) {
		n.pins.Put(key, b)
		return b, true
	}
	b, ok := n.currentRing().Lookup(key, usable)
	if ok {
		n.pins.Put(key, b)
	}
	return b, ok
}

func (n *Node) report(o proxy.Outcome) {
	n.requests.Inc(o.Backend, strconv.Itoa(o.Status))
	if o.Retried {
		n.retries.Inc(o.Backend)
	}
}
