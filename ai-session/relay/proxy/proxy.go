// Package proxy forwards client requests to a backend.
//
// Requests are sticky: every request with the same affinity key goes to the
// same backend, on every node, for as long as that backend is healthy. The
// key is the X-Relay-Key header if the client sends one, otherwise the
// client's IP address. If a backend cannot be reached or answers with a
// server error, the request is tried on the next backend of the ring.
package proxy

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"
)

// KeyHeader lets a client choose its affinity key, usually a session or user
// id.
const KeyHeader = "X-Relay-Key"

// MaxBody is the largest request body the proxy accepts.
const MaxBody = 1 << 20

// Picker chooses the backend for an affinity key among the usable backends.
type Picker func(key string, usable func(backend string) bool) (string, bool)

// Proxy is an http.Handler that forwards to backends.
type Proxy struct {
	Node    string // the node's id, sent back in X-Relay-Node
	Pick    Picker
	Healthy func(backend string) bool
	Client  *http.Client
	Retries int // extra attempts after the first
}

// New returns a proxy with a five second timeout per attempt and one retry.
func New(node string, pick Picker, healthy func(string) bool) *Proxy {
	return &Proxy{
		Node:    node,
		Pick:    pick,
		Healthy: healthy,
		Client: &http.Client{
			Timeout: 5 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Retries: 1,
	}
}

// hop-by-hop headers belong to one connection and are never forwarded.
var hopHeaders = []string{
	"Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization",
	"Te", "Trailer", "Transfer-Encoding", "Upgrade",
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The body is read once and kept, so that a retry can send it again.
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, MaxBody))
	if err != nil {
		http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Without a key from the client, the client's address keeps it on one
	// backend.
	key := r.Header.Get(KeyHeader)
	if key == "" {
		key = r.RemoteAddr
	}

	tried := make(map[string]bool)
	usable := func(b string) bool { return !tried[b] && p.Healthy(b) }

	for attempt := 0; attempt <= p.Retries; attempt++ {
		backend, ok := p.Pick(key, usable)
		if !ok {
			break
		}
		tried[backend] = true

		resp, err := p.forward(r, backend, body)
		last := attempt == p.Retries
		if err != nil {
			if last {
				http.Error(w, "backend "+backend+" failed: "+err.Error(), http.StatusBadGateway)
				return
			}
			continue
		}
		if resp.StatusCode >= 500 && !last {
			resp.Body.Close()
			continue
		}
		p.reply(w, resp, backend)
		return
	}
	http.Error(w, "no backend available", http.StatusServiceUnavailable)
}

func (p *Proxy) forward(r *http.Request, backend string, body []byte) (*http.Response, error) {
	url := "http://" + backend + r.URL.RequestURI()
	out, err := http.NewRequestWithContext(r.Context(), r.Method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	out.Header = r.Header.Clone()
	removeHop(out.Header)
	out.Host = r.Host
	out.Header.Set("X-Forwarded-For", r.RemoteAddr)
	out.Header.Set("X-Forwarded-Host", r.Host)
	out.Header.Set("X-Relay-Node", p.Node)
	return p.Client.Do(out)
}

func (p *Proxy) reply(w http.ResponseWriter, resp *http.Response, backend string) {
	defer resp.Body.Close()
	removeHop(resp.Header)
	for k, vs := range resp.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	w.Header().Set("X-Relay-Node", p.Node)
	w.Header().Set("X-Relay-Backend", backend)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func removeHop(h http.Header) {
	for f := range strings.SplitSeq(h.Get("Connection"), ",") {
		if f = strings.TrimSpace(f); f != "" {
			h.Del(f)
		}
	}
	for _, k := range hopHeaders {
		h.Del(k)
	}
}
