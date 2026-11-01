package proxy

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// backend starts a test server that answers with its name, or with the
// status code returned by fail, if it is not zero.
func backend(t *testing.T, name string, fail func() int) (string, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		body, _ := io.ReadAll(r.Body)
		if fail != nil {
			if code := fail(); code != 0 {
				w.WriteHeader(code)
				return
			}
		}
		w.Header().Set("X-Seen-Key", r.Header.Get(KeyHeader))
		w.Header().Set("X-Seen-For", r.Header.Get("X-Forwarded-For"))
		io.WriteString(w, name+" "+r.Method+" "+r.URL.RequestURI()+" "+string(body))
	}))
	t.Cleanup(srv.Close)
	return strings.TrimPrefix(srv.URL, "http://"), &calls
}

// closedAddr is an address nothing listens on.
func closedAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

// inOrder picks the first usable backend of a fixed list, whatever the key.
func inOrder(backends ...string) Picker {
	return func(_ string, usable func(string) bool) (string, bool) {
		for _, b := range backends {
			if usable(b) {
				return b, true
			}
		}
		return "", false
	}
}

func newProxy(pick Picker) *Proxy {
	return &Proxy{
		Node: "test", Pick: pick, Healthy: func(string) bool { return true },
		Client: NewClient(), Retries: 1, TryTimeout: 2 * time.Second, MaxBody: 1 << 20,
	}
}

func do(t *testing.T, h http.Handler, method, path, key, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if key != "" {
		req.Header.Set(KeyHeader, key)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestForwardsTheRequest(t *testing.T) {
	one, _ := backend(t, "one", nil)
	rec := do(t, newProxy(inOrder(one)), "POST", "/orders?x=1", "alice", "hello")
	if rec.Code != 200 || rec.Body.String() != "one POST /orders?x=1 hello" {
		t.Fatalf("%d %q", rec.Code, rec.Body)
	}
	if rec.Header().Get("X-Relay-Backend") != one || rec.Header().Get("X-Seen-Key") != "alice" {
		t.Fatalf("headers %v", rec.Header())
	}
}

func TestGetIsRetriedOn503(t *testing.T) {
	sick, sickCalls := backend(t, "sick", func() int { return 503 })
	well, _ := backend(t, "well", nil)
	rec := do(t, newProxy(inOrder(sick, well)), "GET", "/", "alice", "")
	if rec.Code != 200 || rec.Header().Get("X-Relay-Backend") != well || sickCalls.Load() != 1 {
		t.Fatalf("%d from %q", rec.Code, rec.Header().Get("X-Relay-Backend"))
	}
}

func TestPostIsNotRetriedAfterTheBackendAnswered(t *testing.T) {
	for _, code := range []int{500, 502, 503} {
		sick, _ := backend(t, "sick", func() int { return code })
		well, wellCalls := backend(t, "well", nil)
		rec := do(t, newProxy(inOrder(sick, well)), "POST", "/orders", "alice", "{}")
		if rec.Code != code || wellCalls.Load() != 0 {
			t.Fatalf("%d: client got %d, second backend called %d times", code, rec.Code, wellCalls.Load())
		}
	}
}

func TestPostIsRetriedWhenTheConnectionIsRefused(t *testing.T) {
	well, _ := backend(t, "well", nil)
	rec := do(t, newProxy(inOrder(closedAddr(t), well)), "POST", "/orders", "alice", "{}")
	if rec.Code != 200 || rec.Header().Get("X-Relay-Backend") != well {
		t.Fatalf("%d: %s", rec.Code, rec.Body)
	}
}

func TestNoBackend(t *testing.T) {
	if rec := do(t, newProxy(inOrder()), "GET", "/", "", ""); rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestTheRealFailureIsReported(t *testing.T) {
	rec := do(t, newProxy(inOrder(closedAddr(t))), "GET", "/", "", "")
	if rec.Code != http.StatusBadGateway || !strings.Contains(rec.Body.String(), "refused") {
		t.Fatalf("%d %q", rec.Code, rec.Body)
	}
}

func TestSlowBackendTimesOut(t *testing.T) {
	slow, _ := backend(t, "slow", func() int { time.Sleep(300 * time.Millisecond); return 0 })
	p := newProxy(inOrder(slow))
	p.TryTimeout = 50 * time.Millisecond
	if rec := do(t, p, "GET", "/", "", ""); rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestBodyLimit(t *testing.T) {
	one, _ := backend(t, "one", nil)
	p := newProxy(inOrder(one))
	p.MaxBody = 4
	if rec := do(t, p, "POST", "/", "", "too long"); rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestHopHeadersAreDropped(t *testing.T) {
	var seen http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { seen = r.Header.Clone() }))
	defer srv.Close()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Connection", "X-Secret")
	req.Header.Set("X-Secret", "1")
	req.Header.Set("Keep-Alive", "timeout=5")
	newProxy(inOrder(strings.TrimPrefix(srv.URL, "http://"))).ServeHTTP(httptest.NewRecorder(), req)
	if seen.Get("X-Secret") != "" || seen.Get("Keep-Alive") != "" || seen.Get(RequestIDHeader) == "" {
		t.Fatalf("headers %v", seen)
	}
}

func request(remote, xff string) *http.Request {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = remote
	if xff != "" {
		r.Header.Set("X-Forwarded-For", xff)
	}
	return r
}

func TestClientIP(t *testing.T) {
	trusted := []netip.Prefix{netip.MustParsePrefix("10.0.0.0/8")}
	for _, tc := range []struct{ remote, xff, want string }{
		{"203.0.113.7:50001", "", "203.0.113.7"},
		{"203.0.113.7:50002", "198.51.100.1", "203.0.113.7"}, // not a trusted proxy: header ignored
		{"10.1.2.3:40000", "198.51.100.1", "198.51.100.1"},   // the corporate proxy
		{"10.1.2.3:40000", "", "10.1.2.3"},
		{"10.1.2.3:40000", "garbage", "10.1.2.3"},
		{"[::ffff:203.0.113.7]:5", "", "203.0.113.7"},
	} {
		if got := ClientIP(request(tc.remote, tc.xff), trusted); got != tc.want {
			t.Errorf("%s with %q: got %s, want %s", tc.remote, tc.xff, got, tc.want)
		}
	}
}

func TestForwardedForIsAppended(t *testing.T) {
	if got := forwardedFor(request("203.0.113.7:1", "198.51.100.1")); got != "198.51.100.1, 203.0.113.7" {
		t.Fatalf("got %q", got)
	}
	if got := forwardedFor(request("203.0.113.7:1", "")); got != "203.0.113.7" {
		t.Fatalf("got %q", got)
	}
}
