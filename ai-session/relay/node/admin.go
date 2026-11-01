package node

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"net/http/pprof"
	"net/netip"
	"runtime"
	"strings"
	"time"

	"relay/cluster"
	"relay/health"
	"relay/ring"
)

// Handler serves client traffic on / and the node's own API under /relay/.
func (n *Node) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /relay/status", n.status)
	mux.HandleFunc("PUT /relay/backends/{backend}", n.admin(n.putBackend))
	mux.HandleFunc("DELETE /relay/backends/{backend}", n.admin(n.deleteBackend))
	mux.HandleFunc("POST /relay/gossip", n.gossiper.Handler)
	mux.Handle("GET /relay/metrics", n.metrics.Handler())

	mux.HandleFunc("GET /relay/debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /relay/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /relay/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /relay/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /relay/debug/pprof/trace", pprof.Trace)

	mux.HandleFunc("/relay/", http.NotFound) // the node's own space; never proxied
	mux.Handle("/", n.proxy)
	return pprofPrefix(mux)
}

// pprofPrefix lets pprof.Index find profiles under /relay/debug/pprof/: it
// only knows the path /debug/pprof/.
func pprofPrefix(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rest, ok := strings.CutPrefix(r.URL.Path, "/relay/debug/pprof/"); ok && rest != "" &&
			rest != "cmdline" && rest != "profile" && rest != "symbol" && rest != "trace" {
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/debug/pprof/" + rest
			pprof.Index(w, r2)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// admin requires the admin token, when one is configured.
func (n *Node) admin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if n.cfg.AdminToken != "" {
			got, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			if subtle.ConstantTimeCompare([]byte(got), []byte(n.cfg.AdminToken)) != 1 {
				http.Error(w, "admin token required", http.StatusUnauthorized)
				return
			}
		}
		h(w, r)
	}
}

// PUT /relay/backends/{backend} with {"status": "active"|"draining", "weight": 1..10};
// either field may be left out to keep its current value. An unknown backend
// is added.
func (n *Node) putBackend(w http.ResponseWriter, r *http.Request) {
	backend, err := backendParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var c cluster.Change
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		http.Error(w, "bad body: "+err.Error(), http.StatusBadRequest)
		return
	}
	e, err := n.state.Apply(backend, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("backend set", "node", n.cfg.ID, "backend", e.Backend, "status", e.Status, "weight", e.Weight, "version", e.Version)
	writeJSON(w, http.StatusOK, e)
}

func (n *Node) deleteBackend(w http.ResponseWriter, r *http.Request) {
	backend, err := backendParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	e, err := n.state.Remove(backend)
	if errors.Is(err, cluster.ErrUnknown) {
		http.Error(w, "unknown backend "+backend, http.StatusNotFound)
		return
	}
	slog.Info("backend removed", "node", n.cfg.ID, "backend", backend, "version", e.Version)
	writeJSON(w, http.StatusOK, e)
}

func backendParam(r *http.Request) (string, error) {
	b := r.PathValue("backend")
	ap, err := netip.ParseAddrPort(b)
	if err != nil {
		return "", errors.New("backend must be ip:port")
	}
	return ap.String(), nil
}

// ---------------------------------------------------------------- status

type backendStatus struct {
	cluster.Entry
	Healthy bool   `json:"healthy"`
	LastErr string `json:"last_error,omitempty"`
}

type statusReply struct {
	Node     string               `json:"node"`
	Clock    uint64               `json:"clock"`
	Uptime   string               `json:"uptime"`
	Ring     []ring.Member        `json:"ring"`
	Backends []backendStatus      `json:"backends"`
	Peers    []cluster.PeerStatus `json:"peers"`
	Pins     int                  `json:"pins"`
	Runtime  runtimeStatus        `json:"runtime"`
}

type runtimeStatus struct {
	Goroutines int     `json:"goroutines"`
	HeapMB     float64 `json:"heap_mb"`
}

func (n *Node) status(w http.ResponseWriter, _ *http.Request) {
	health := map[string]health.Status{}
	for _, s := range n.checker.Statuses() {
		health[s.Backend] = s
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	reply := statusReply{
		Node:   n.cfg.ID,
		Clock:  n.state.Clock(),
		Uptime: time.Since(n.started).Round(time.Second).String(),
		Ring:   n.currentRing().Members(),
		Peers:  n.gossiper.Statuses(),
		Pins:   n.pins.Len(),
		Runtime: runtimeStatus{
			Goroutines: runtime.NumGoroutine(),
			HeapMB:     math.Round(float64(ms.HeapAlloc)/(1<<20)*10) / 10,
		},
	}
	for _, e := range n.state.Entries() {
		s := backendStatus{Entry: e, Healthy: n.checker.Healthy(e.Backend)}
		if h, ok := health[e.Backend]; ok {
			s.LastErr = h.LastErr
		}
		reply.Backends = append(reply.Backends, s)
	}
	writeJSON(w, http.StatusOK, reply)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
