// Package proxy forwards client requests to a backend.
//
// Requests are sticky: every request with the same affinity key goes to the
// same backend. The key is the X-Relay-Key header if the client sends one,
// otherwise the client's IP address (see ClientIP). Choosing the backend is
// the caller's job (Picker); the proxy forwards, retries when that is safe
// (see shouldRetry), and reports what happened.
package proxy

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

// KeyHeader lets a client choose its affinity key, usually a session or
// user id.
const KeyHeader = "X-Relay-Key"

// RequestIDHeader identifies a request in the logs of Relay and the backend.
const RequestIDHeader = "X-Request-Id"

// Picker chooses the backend for a key among the usable backends.
type Picker func(key string, usable func(backend string) bool) (string, bool)

// Outcome is reported once per attempt.
type Outcome struct {
	Backend string
	Status  int // 0 if no response
	Err     error
	Retried bool // another attempt followed
}

// Proxy is an http.Handler that forwards to backends.
type Proxy struct {
	Node       string
	Pick       Picker
	Healthy    func(backend string) bool
	Report     func(Outcome)
	Client     *http.Client
	Retries    int
	TryTimeout time.Duration
	MaxBody    int64
	Trusted    []netip.Prefix
}

// NewClient returns the client used for forwarding: redirects are passed to
// the client, not followed.
func NewClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConnsPerHost = 64
	return &http.Client{
		Transport: t,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// hop-by-hop headers belong to one connection and are never forwarded.
var hopHeaders = []string{
	"Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization",
	"Te", "Trailer", "Transfer-Encoding", "Upgrade",
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id := r.Header.Get(RequestIDHeader)
	if id == "" {
		id = newID()
		r.Header.Set(RequestIDHeader, id)
	}
	w.Header().Set(RequestIDHeader, id)

	// The body is read once and kept, so that a retry can send it again.
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, p.MaxBody))
	if err != nil {
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "reading request body: "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	key := r.Header.Get(KeyHeader)
	if key == "" {
		key = ClientIP(r, p.Trusted)
	}

	tried := make(map[string]bool)
	usable := func(b string) bool { return !tried[b] && p.Healthy(b) }

	var (
		lastErr     error
		lastBackend string
	)
	for attempt := 1; attempt <= 1+p.Retries; attempt++ {
		backend, ok := p.Pick(key, usable)
		if !ok {
			break
		}
		tried[backend] = true
		lastBackend = backend

		ctx, cancel := context.WithTimeout(r.Context(), p.TryTimeout)
		resp, err := p.forward(ctx, r, backend, body)
		retry := attempt <= p.Retries && shouldRetry(r, resp, err)
		p.report(Outcome{Backend: backend, Status: status(resp), Err: err, Retried: retry})

		if err == nil && !retry {
			p.reply(w, r, resp, backend)
			resp.Body.Close()
			cancel()
			slog.Info("request", "id", id, "method", r.Method, "path", r.URL.Path, "key", key,
				"backend", backend, "status", resp.StatusCode, "attempts", attempt, "took", time.Since(start).Round(time.Millisecond))
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		cancel()
		if err != nil {
			lastErr = err
		} else {
			lastErr = errors.New("backend answered " + resp.Status)
		}
		if !retry {
			break
		}
	}

	code, msg := http.StatusServiceUnavailable, "no backend available"
	if lastErr != nil {
		code, msg = http.StatusBadGateway, "backend "+lastBackend+" failed: "+lastErr.Error()
		if timedOut(lastErr) {
			code = http.StatusGatewayTimeout
		}
	}
	w.Header().Set("X-Relay-Node", p.Node)
	http.Error(w, msg, code)
	slog.Info("request", "id", id, "method", r.Method, "path", r.URL.Path, "key", key,
		"backend", lastBackend, "status", code, "attempts", len(tried), "err", lastErr)
}

func (p *Proxy) report(o Outcome) {
	if p.Report != nil {
		p.Report(o)
	}
}

func (p *Proxy) forward(ctx context.Context, r *http.Request, backend string, body []byte) (*http.Response, error) {
	out, err := http.NewRequestWithContext(ctx, r.Method, "http://"+backend+r.URL.RequestURI(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	out.Header = r.Header.Clone()
	removeHop(out.Header)
	out.Host = r.Host
	out.Header.Set("X-Forwarded-For", forwardedFor(r))
	out.Header.Set("X-Forwarded-Host", r.Host)
	out.Header.Set("X-Relay-Node", p.Node)
	resp, err := p.Client.Do(out)
	if err != nil {
		return nil, err
	}
	// Read the body while the attempt's context is alive; responses through
	// Relay are small, and a retry decision needs the whole answer anyway.
	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(b))
	resp.ContentLength = int64(len(b))
	return resp, nil
}

func (p *Proxy) reply(w http.ResponseWriter, r *http.Request, resp *http.Response, backend string) {
	removeHop(resp.Header)
	for k, vs := range resp.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	// The body was read in full, so its length is known; a HEAD answer or a
	// 204 or 304 keeps the backend's own header.
	if r.Method != http.MethodHead && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		w.Header().Set("Content-Length", strconv.FormatInt(resp.ContentLength, 10))
	}
	w.Header().Set("X-Relay-Node", p.Node)
	w.Header().Set("X-Relay-Backend", backend)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func status(resp *http.Response) int {
	if resp == nil {
		return 0
	}
	return resp.StatusCode
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

func newID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
