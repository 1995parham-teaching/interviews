package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestUnknownBackendIsHealthy(t *testing.T) {
	if !New().Healthy("10.0.0.1:80") {
		t.Fatal("a backend that was never checked should take traffic")
	}
}

func TestDownAfterFallFailures(t *testing.T) {
	c := New()
	b := "10.0.0.1:80"
	c.Record(b, false)
	c.Record(b, false)
	if !c.Healthy(b) {
		t.Fatal("down after two failures; Fall is three")
	}
	c.Record(b, false)
	if c.Healthy(b) {
		t.Fatal("still up after three failures")
	}
}

func TestUpAfterRiseSuccesses(t *testing.T) {
	c := New()
	b := "10.0.0.1:80"
	for range 3 {
		c.Record(b, false)
	}
	c.Record(b, true)
	if c.Healthy(b) {
		t.Fatal("up after one success; Rise is two")
	}
	c.Record(b, true)
	if !c.Healthy(b) {
		t.Fatal("still down after two successes")
	}
}

func TestAFailureInterruptsRecovery(t *testing.T) {
	c := New()
	b := "10.0.0.1:80"
	for range 3 {
		c.Record(b, false)
	}
	c.Record(b, true)
	c.Record(b, false)
	c.Record(b, true)
	if c.Healthy(b) {
		t.Fatal("up after two successes that were not in a row")
	}
}

func TestCheckAllProbesHealth(t *testing.T) {
	var up atomic.Bool
	up.Store(true)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" || !up.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	defer srv.Close()
	b := strings.TrimPrefix(srv.URL, "http://")

	c := New()
	c.CheckAll(context.Background(), []string{b})
	if !c.Healthy(b) {
		t.Fatal("down although /health answers 200")
	}
	up.Store(false)
	for range 3 {
		c.CheckAll(context.Background(), []string{b})
	}
	if c.Healthy(b) {
		t.Fatal("up although /health answered 503 three times")
	}
}

func TestCheckAllForgetsRemovedBackends(t *testing.T) {
	c := New()
	c.Record("10.0.0.1:80", true)
	c.CheckAll(context.Background(), nil)
	if len(c.Statuses()) != 0 {
		t.Fatalf("still tracking %v", c.Statuses())
	}
}
