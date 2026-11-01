package node

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"relay/cluster"
	"relay/config"
)

type testNode struct {
	*Node
	url string
}

// pair starts two nodes that know each other, with the same seed backends.
func pair(t *testing.T, backends string, edit func(*config.Config)) (testNode, testNode) {
	t.Helper()
	var a, b *Node
	sa := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { a.Handler().ServeHTTP(w, r) }))
	sb := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { b.Handler().ServeHTTP(w, r) }))
	t.Cleanup(sa.Close)
	t.Cleanup(sb.Close)
	build := func(id, peer string) *Node {
		c := config.Default()
		c.ID, c.Peers = id, []string{strings.TrimPrefix(peer, "http://")}
		bs, err := config.ParseBackends(backends)
		if err != nil {
			t.Fatal(err)
		}
		c.Backends = bs
		if edit != nil {
			edit(&c)
		}
		n, err := New(c)
		if err != nil {
			t.Fatal(err)
		}
		return n
	}
	a, b = build("a", sb.URL), build("b", sa.URL)
	return testNode{a, sa.URL}, testNode{b, sb.URL}
}

func call(t *testing.T, method, url, body, token string) int {
	t.Helper()
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	return resp.StatusCode
}

func status(t *testing.T, n testNode) statusReply {
	t.Helper()
	resp, err := http.Get(n.url + "/relay/status")
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

func onRing(s statusReply) []string {
	var out []string
	for _, m := range s.Ring {
		out = append(out, m.Backend)
	}
	return out
}

func TestDrainSpreadsByGossip(t *testing.T) {
	a, b := pair(t, "10.0.0.1:80,10.0.0.2:80", nil)
	if code := call(t, "PUT", a.url+"/relay/backends/10.0.0.2:80", `{"status":"draining"}`, ""); code != 200 {
		t.Fatalf("PUT: %d", code)
	}
	a.Gossiper().Round(context.Background())
	if got := onRing(status(t, b)); !slices.Equal(got, []string{"10.0.0.1:80"}) {
		t.Fatalf("ring on b: %v", got)
	}
}

func TestAddWithWeight(t *testing.T) {
	a, b := pair(t, "10.0.0.1:80", nil)
	call(t, "PUT", a.url+"/relay/backends/10.0.0.3:80", `{"status":"active","weight":2}`, "")
	b.Gossiper().Round(context.Background()) // pull
	s := status(t, b)
	if len(s.Ring) != 2 || s.Ring[1].Backend != "10.0.0.3:80" || s.Ring[1].Weight != 2 {
		t.Fatalf("ring on b: %+v", s.Ring)
	}
}

func TestRemoveStaysRemoved(t *testing.T) {
	a, b := pair(t, "10.0.0.1:80,10.0.0.2:80", nil)
	if code := call(t, "DELETE", a.url+"/relay/backends/10.0.0.2:80", "", ""); code != 200 {
		t.Fatalf("DELETE: %d", code)
	}
	for range 3 {
		a.Gossiper().Round(context.Background())
		b.Gossiper().Round(context.Background())
	}
	for _, n := range []testNode{a, b} {
		if got := onRing(status(t, n)); !slices.Equal(got, []string{"10.0.0.1:80"}) {
			t.Fatalf("ring on %s: %v", n.cfg.ID, got)
		}
	}
	if code := call(t, "DELETE", a.url+"/relay/backends/10.0.0.2:80", "", ""); code != 404 {
		t.Fatalf("second DELETE: %d", code)
	}
}

func TestAdminToken(t *testing.T) {
	a, _ := pair(t, "10.0.0.1:80", func(c *config.Config) { c.AdminToken = "t0k" })
	url := a.url + "/relay/backends/10.0.0.1:80"
	if code := call(t, "PUT", url, `{"status":"draining"}`, ""); code != 401 {
		t.Fatalf("without token: %d", code)
	}
	if code := call(t, "PUT", url, `{"status":"draining"}`, "wrong"); code != 401 {
		t.Fatalf("wrong token: %d", code)
	}
	if code := call(t, "PUT", url, `{"status":"draining"}`, "t0k"); code != 200 {
		t.Fatalf("right token: %d", code)
	}
}

func TestBadAdminRequests(t *testing.T) {
	a, _ := pair(t, "10.0.0.1:80", nil)
	for _, tc := range []struct {
		method, path, body string
		want               int
	}{
		{"PUT", "/relay/backends/nohost", `{"status":"active"}`, 400},
		{"PUT", "/relay/backends/10.0.0.1:80", `{"status":"asleep"}`, 400},
		{"PUT", "/relay/backends/10.0.0.1:80", `{"weight":0}`, 400},
		{"PUT", "/relay/backends/10.0.0.1:80", `{"colour":"red"}`, 400},
		{"PUT", "/relay/backends/10.0.0.1:80", `not json`, 400},
		{"DELETE", "/relay/backends/10.0.0.9:80", ``, 404},
	} {
		if code := call(t, tc.method, a.url+tc.path, tc.body, ""); code != tc.want {
			t.Errorf("%s %s %s: %d, want %d", tc.method, tc.path, tc.body, code, tc.want)
		}
	}
}

func TestPinnedKeyStaysWhenABackendJoins(t *testing.T) {
	a, _ := pair(t, "10.0.0.1:80,10.0.0.2:80", nil)
	always := func(string) bool { return true }
	before := map[string]string{}
	for i := range 200 {
		k := "user-" + string(rune('A'+i%26)) + string(rune('a'+i/26))
		before[k], _ = a.pick(k, always)
	}
	a.State().Apply("10.0.0.3:80", changeActive(1))
	for k, b := range before {
		if got, _ := a.pick(k, always); got != b {
			t.Fatalf("%s moved from %s to %s", k, b, got)
		}
	}
}

func TestMetricsAndPprof(t *testing.T) {
	a, _ := pair(t, "10.0.0.1:80", nil)
	for _, path := range []string{"/relay/metrics", "/relay/debug/pprof/", "/relay/debug/pprof/goroutine?debug=1"} {
		resp, err := http.Get(a.url + path)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Errorf("%s: %d", path, resp.StatusCode)
		}
	}
	if s := status(t, a); s.Runtime.Goroutines == 0 || s.Node != "a" {
		t.Fatalf("status %+v", s)
	}
}

func changeActive(w int) (c cluster.Change) {
	s := cluster.Active
	c.Status, c.Weight = &s, &w
	return c
}
