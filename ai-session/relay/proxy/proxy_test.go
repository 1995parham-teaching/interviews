package proxy

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"relay/ring"
)

// backend starts a test server that answers with its own name.
func backend(t *testing.T, name string) string {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("X-Seen-Key", r.Header.Get(KeyHeader))
		io.WriteString(w, name+" "+r.Method+" "+r.URL.RequestURI()+" "+string(body))
	}))
	t.Cleanup(srv.Close)
	return strings.TrimPrefix(srv.URL, "http://")
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

func always(string) bool { return true }

func newProxy(backends ...string) *Proxy {
	r := ring.New(backends)
	return New("test", r.Lookup, always)
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
	p := newProxy(backend(t, "one"))
	rec := do(t, p, "POST", "/orders?x=1", "alice", "hello")
	if rec.Code != 200 {
		t.Fatalf("status %d: %s", rec.Code, rec.Body)
	}
	if got := rec.Body.String(); got != "one POST /orders?x=1 hello" {
		t.Fatalf("body %q", got)
	}
	if rec.Header().Get("X-Relay-Node") != "test" || rec.Header().Get("X-Seen-Key") != "alice" {
		t.Fatalf("headers %v", rec.Header())
	}
}

func TestSameKeySameBackend(t *testing.T) {
	p := newProxy(backend(t, "one"), backend(t, "two"), backend(t, "three"))
	for _, key := range []string{"alice", "bob", "carol"} {
		first := do(t, p, "GET", "/", key, "").Header().Get("X-Relay-Backend")
		for range 5 {
			if got := do(t, p, "GET", "/", key, "").Header().Get("X-Relay-Backend"); got != first {
				t.Fatalf("%s went to %s, then %s", key, first, got)
			}
		}
	}
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

func TestRetriesOnAnotherBackend(t *testing.T) {
	alive := backend(t, "alive")
	p := New("test", inOrder(closedAddr(t), alive), always)

	rec := do(t, p, "GET", "/", "alice", "")
	if rec.Code != 200 || rec.Header().Get("X-Relay-Backend") != alive {
		t.Fatalf("status %d from %q", rec.Code, rec.Header().Get("X-Relay-Backend"))
	}
}

func TestUnhealthyBackendIsSkipped(t *testing.T) {
	sick := backend(t, "sick")
	well := backend(t, "well")
	p := New("test", inOrder(sick, well), func(b string) bool { return b != sick })

	if got := do(t, p, "GET", "/", "alice", "").Header().Get("X-Relay-Backend"); got != well {
		t.Fatalf("went to %q", got)
	}
}

func TestNoBackend(t *testing.T) {
	p := New("test", ring.New(nil).Lookup, always)
	if rec := do(t, p, "GET", "/", "", ""); rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestHopHeadersAreDropped(t *testing.T) {
	var seen http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Clone()
	}))
	defer srv.Close()
	p := newProxy(strings.TrimPrefix(srv.URL, "http://"))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Connection", "X-Secret")
	req.Header.Set("X-Secret", "1")
	req.Header.Set("Keep-Alive", "timeout=5")
	p.ServeHTTP(httptest.NewRecorder(), req)
	if seen.Get("X-Secret") != "" || seen.Get("Keep-Alive") != "" {
		t.Fatalf("hop headers forwarded: %v", seen)
	}
}
