package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

var errProbe = errors.New("probe failed")

const b = "10.0.0.1:80"

func TestUnknownBackendIsHealthy(t *testing.T) {
	if !New(time.Second, 3, 2).Healthy(b) {
		t.Fatal("a backend that was never checked should take traffic")
	}
}

func TestDownAfterFallFailuresInARow(t *testing.T) {
	c := New(time.Second, 3, 2)
	c.Record(b, errProbe)
	c.Record(b, errProbe)
	if !c.Healthy(b) {
		t.Fatal("down after two failures")
	}
	c.Record(b, errProbe)
	if c.Healthy(b) {
		t.Fatal("up after three failures")
	}
}

func TestScatteredFailuresKeepItUp(t *testing.T) {
	c := New(time.Second, 3, 2)
	for range 10 {
		c.Record(b, errProbe)
		c.Record(b, errProbe)
		c.Record(b, nil)
	}
	if !c.Healthy(b) {
		t.Fatal("down although failures were never three in a row")
	}
}

func TestUpAfterRiseSuccessesInARow(t *testing.T) {
	c := New(time.Second, 3, 2)
	for range 3 {
		c.Record(b, errProbe)
	}
	c.Record(b, nil)
	c.Record(b, errProbe)
	c.Record(b, nil)
	if c.Healthy(b) {
		t.Fatal("up after successes that were not in a row")
	}
	c.Record(b, nil)
	if !c.Healthy(b) {
		t.Fatal("still down after two successes in a row")
	}
}

func TestOnFlip(t *testing.T) {
	c := New(time.Second, 1, 1)
	var flips []bool
	c.OnFlip(func(_ string, up bool) { flips = append(flips, up) })
	c.Record(b, nil)
	c.Record(b, errProbe)
	c.Record(b, errProbe)
	c.Record(b, nil)
	if len(flips) != 2 || flips[0] || !flips[1] {
		t.Fatalf("flips %v", flips)
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
	addr := strings.TrimPrefix(srv.URL, "http://")

	c := New(time.Second, 3, 2)
	c.CheckAll(context.Background(), []string{addr})
	if !c.Healthy(addr) {
		t.Fatal("down although /health answers 200")
	}
	up.Store(false)
	for range 3 {
		c.CheckAll(context.Background(), []string{addr})
	}
	if c.Healthy(addr) {
		t.Fatal("up although /health answered 503 three times")
	}
	if st := c.Statuses(); st[0].LastErr == "" {
		t.Fatalf("no error recorded: %+v", st)
	}
}

func TestCheckAllForgetsRemovedBackends(t *testing.T) {
	c := New(time.Second, 3, 2)
	c.Record(b, nil)
	c.CheckAll(context.Background(), nil)
	if len(c.Statuses()) != 0 {
		t.Fatalf("still tracking %v", c.Statuses())
	}
}
