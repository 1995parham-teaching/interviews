// Command backend is a small HTTP service to put behind Relay.
//
// It answers every request with a JSON description of what it received, and
// keeps a list of orders so that side effects can be counted. Its behaviour
// can be changed while it runs, to play a slow or failing service.
//
//	GET  /health          200, or 503 while health is set to fail
//	POST /orders          records an order and returns it
//	GET  /orders          lists the orders this backend recorded
//	DELETE /orders        forgets them
//	PUT  /control         {"fail_rate": 0.5, "crash_rate": 0, "fail_health": false, "delay": "200ms"}
//	anything else         echoes the request
package main

import (
	"encoding/json"
	"flag"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"sync"
	"time"
)

type control struct {
	FailRate   float64  `json:"fail_rate"`   // share of requests answered with 500
	CrashRate  float64  `json:"crash_rate"`  // share of requests whose connection is dropped without an answer
	FailHealth bool     `json:"fail_health"` // answer /health with 503
	Delay      duration `json:"delay"`       // added to every request but /health
}

type order struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
	Key  string `json:"key"`
}

type server struct {
	name string

	mu     sync.Mutex
	ctl    control
	orders []order
}

func main() {
	listen := flag.String("listen", "127.0.0.1:9001", "address to serve on")
	name := flag.String("name", "", "name to answer with (default: the listen address)")
	flag.Parse()
	if *name == "" {
		*name = *listen
	}

	s := &server{name: *name}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("PUT /control", s.control)
	mux.HandleFunc("POST /orders", s.wrap(s.createOrder))
	mux.HandleFunc("GET /orders", s.listOrders)
	mux.HandleFunc("DELETE /orders", s.clearOrders)
	mux.HandleFunc("/", s.wrap(s.echo))

	slog.Info("backend up", "name", *name, "listen", *listen)
	if err := http.ListenAndServe(*listen, mux); err != nil {
		slog.Error("backend stopped", "err", err)
		os.Exit(1)
	}
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	fail := s.ctl.FailHealth
	s.mu.Unlock()
	if fail {
		http.Error(w, "unhealthy", http.StatusServiceUnavailable)
		return
	}
	io.WriteString(w, "ok\n")
}

func (s *server) control(w http.ResponseWriter, r *http.Request) {
	var c control
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	s.ctl = c
	s.mu.Unlock()
	slog.Info("control", "name", s.name, "fail_rate", c.FailRate, "crash_rate", c.CrashRate, "fail_health", c.FailHealth, "delay", time.Duration(c.Delay))
	writeJSON(w, http.StatusOK, c)
}

// wrap applies the delay and the failure rate. A failed request is still
// handled first, the way a real service fails after doing the work: the
// database write succeeded, the response did not.
func (s *server) wrap(h func(http.ResponseWriter, *http.Request) any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		c := s.ctl
		s.mu.Unlock()
		time.Sleep(time.Duration(c.Delay))
		reply := h(w, r)
		if rand.Float64() < c.CrashRate {
			// The work is done; the process dies before it can answer.
			slog.Info("request", "name", s.name, "method", r.Method, "path", r.URL.Path, "status", "dropped")
			if conn, _, err := http.NewResponseController(w).Hijack(); err == nil {
				conn.Close()
				return
			}
			panic(http.ErrAbortHandler)
		}
		if rand.Float64() < c.FailRate {
			slog.Info("request", "name", s.name, "method", r.Method, "path", r.URL.Path, "status", 500)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		slog.Info("request", "name", s.name, "method", r.Method, "path", r.URL.Path, "status", 200)
		writeJSON(w, http.StatusOK, reply)
	}
}

func (s *server) createOrder(_ http.ResponseWriter, r *http.Request) any {
	body, _ := io.ReadAll(r.Body)
	s.mu.Lock()
	defer s.mu.Unlock()
	o := order{ID: len(s.orders) + 1, Body: string(body), Key: r.Header.Get("X-Relay-Key")}
	s.orders = append(s.orders, o)
	return map[string]any{"backend": s.name, "order": o}
}

func (s *server) listOrders(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"backend": s.name, "orders": s.orders})
}

func (s *server) clearOrders(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	s.orders = nil
	s.mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) echo(_ http.ResponseWriter, r *http.Request) any {
	body, _ := io.ReadAll(r.Body)
	return map[string]any{
		"backend":         s.name,
		"method":          r.Method,
		"path":            r.URL.RequestURI(),
		"body":            string(body),
		"relay_key":       r.Header.Get("X-Relay-Key"),
		"x_forwarded_for": r.Header.Get("X-Forwarded-For"),
		"relay_node":      r.Header.Get("X-Relay-Node"),
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

// duration reads "200ms" from JSON.
type duration time.Duration

func (d *duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*d = 0
		return nil
	}
	v, err := time.ParseDuration(s)
	*d = duration(v)
	return err
}

func (d duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}
