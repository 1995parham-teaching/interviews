package cluster

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// peer serves one gossiper's Handler and returns its address.
func peer(t *testing.T, g *Gossiper) string {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(g.Handler))
	t.Cleanup(srv.Close)
	return strings.TrimPrefix(srv.URL, "http://")
}

func TestExchangeGoesBothWays(t *testing.T) {
	a := &Gossiper{State: newState("a", nil), Timeout: time.Second}
	b := &Gossiper{State: newState("b", nil), Timeout: time.Second}
	a.Peers = []string{peer(t, b)}

	a.State.Apply("10.0.0.1:80", Change{Status: status(Draining)})
	b.State.Apply("10.0.0.3:80", Change{Status: status(Active), Weight: weight(1)})
	a.Round(context.Background())

	if ea, _ := b.State.Get("10.0.0.1:80"); ea.Status != Draining {
		t.Fatalf("push: b has %+v", ea)
	}
	if _, ok := a.State.Get("10.0.0.3:80"); !ok {
		t.Fatal("pull: a did not learn b's backend")
	}
	if st := a.Statuses(); len(st) != 1 || st[0].LastOK.IsZero() {
		t.Fatalf("peer status %+v", st)
	}
}

func TestSecretIsRequired(t *testing.T) {
	b := &Gossiper{State: newState("b", nil), Timeout: time.Second, Secret: "s3"}
	a := &Gossiper{State: newState("a", nil), Timeout: time.Second, Peers: []string{peer(t, b)}}
	a.State.Apply("10.0.0.1:80", Change{Status: status(Draining)})

	a.Round(context.Background())
	if e, _ := b.State.Get("10.0.0.1:80"); e.Status != Active {
		t.Fatal("b merged gossip without the secret")
	}
	if st := a.Statuses(); st[0].Failures != 1 {
		t.Fatalf("peer status %+v", st)
	}

	a.Secret = "s3"
	a.Round(context.Background())
	if e, _ := b.State.Get("10.0.0.1:80"); e.Status != Draining {
		t.Fatal("b refused gossip with the secret")
	}
}

func TestUnreachablePeerIsRecorded(t *testing.T) {
	a := &Gossiper{State: newState("a", nil), Timeout: 200 * time.Millisecond, Peers: []string{"127.0.0.1:1"}}
	a.Round(context.Background())
	if st := a.Statuses(); st[0].Failures != 1 || st[0].LastError == "" {
		t.Fatalf("peer status %+v", st)
	}
}

func TestInvalidEntriesAreDropped(t *testing.T) {
	b := &Gossiper{State: newState("b", nil), Timeout: time.Second}
	addr := peer(t, b)
	body := `[{"backend":"not-an-address","status":"active","weight":1,"version":9},
	          {"backend":"10.0.0.1:80","status":"sleeping","weight":1,"version":9}]`
	resp, err := http.Post("http://"+addr+"/relay/gossip", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if e, _ := b.State.Get("10.0.0.1:80"); e.Status != Active || len(b.State.Entries()) != 2 {
		t.Fatalf("b merged invalid entries: %+v", b.State.Entries())
	}
}
