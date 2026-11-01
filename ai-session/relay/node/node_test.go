package node

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"relay/cluster"
)

// pair starts two nodes that know each other, with the same seed backends.
func pair(t *testing.T, backends ...string) (*Node, *Node, string, string) {
	t.Helper()
	var a, b *Node
	sa := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { a.Handler().ServeHTTP(w, r) }))
	sb := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { b.Handler().ServeHTTP(w, r) }))
	t.Cleanup(sa.Close)
	t.Cleanup(sb.Close)
	addrA := strings.TrimPrefix(sa.URL, "http://")
	addrB := strings.TrimPrefix(sb.URL, "http://")
	cfg := func(id, peer string) Config {
		return Config{ID: id, Peers: []string{peer}, Backends: backends, CheckInterval: time.Hour, GossipInterval: time.Hour}
	}
	a = New(cfg("a", addrB))
	b = New(cfg("b", addrA))
	return a, b, sa.URL, sb.URL
}

func put(t *testing.T, base, backend string, status cluster.Status) {
	t.Helper()
	req, _ := http.NewRequest("PUT", base+"/relay/backends/"+backend, strings.NewReader(`{"status":"`+string(status)+`"}`))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("PUT %s: %s", backend, resp.Status)
	}
}

func status(t *testing.T, base string) statusReply {
	t.Helper()
	resp, err := http.Get(base + "/relay/status")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var s statusReply
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestDrainSpreadsByGossip(t *testing.T) {
	a, b, urlA, urlB := pair(t, "10.0.0.1:80", "10.0.0.2:80")

	put(t, urlA, "10.0.0.2:80", cluster.Draining)
	a.GossipOnce(context.Background())

	if got := status(t, urlB).Ring; !slices.Equal(got, []string{"10.0.0.1:80"}) {
		t.Fatalf("ring on b: %v", got)
	}
	if !slices.Equal(a.State().Entries(), b.State().Entries()) {
		t.Fatalf("a %v, b %v", a.State().Entries(), b.State().Entries())
	}
}

func TestAddSpreadsByGossip(t *testing.T) {
	_, b, urlA, urlB := pair(t, "10.0.0.1:80")

	put(t, urlA, "10.0.0.3:80", cluster.Active)
	b.GossipOnce(context.Background()) // pull: b asks a

	if got := status(t, urlB).Ring; !slices.Equal(got, []string{"10.0.0.1:80", "10.0.0.3:80"}) {
		t.Fatalf("ring on b: %v", got)
	}
}

func TestBadAdminRequests(t *testing.T) {
	_, _, urlA, _ := pair(t, "10.0.0.1:80")
	for _, tc := range []struct {
		method, path, body string
		want               int
	}{
		{"PUT", "/relay/backends/nohost", `{"status":"active"}`, 400},
		{"PUT", "/relay/backends/10.0.0.1:80", `{"status":"asleep"}`, 400},
		{"PUT", "/relay/backends/10.0.0.1:80", `not json`, 400},
		{"DELETE", "/relay/backends/10.0.0.9:80", ``, 404},
	} {
		req, _ := http.NewRequest(tc.method, urlA+tc.path, strings.NewReader(tc.body))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != tc.want {
			t.Errorf("%s %s %s: %d, want %d", tc.method, tc.path, tc.body, resp.StatusCode, tc.want)
		}
	}
}
