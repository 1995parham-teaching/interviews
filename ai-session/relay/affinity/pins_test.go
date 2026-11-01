package affinity

import (
	"testing"
	"time"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func newPins(c *clock) *Pins {
	p := New(time.Minute)
	p.now = c.now
	return p
}

func TestPutGet(t *testing.T) {
	p := newPins(&clock{t: time.Unix(0, 0)})
	if _, ok := p.Get("alice"); ok {
		t.Fatal("pin out of nowhere")
	}
	p.Put("alice", "a:1")
	if b, ok := p.Get("alice"); !ok || b != "a:1" {
		t.Fatalf("got %q, %v", b, ok)
	}
}

func TestPinsExpireWhenIdle(t *testing.T) {
	c := &clock{t: time.Unix(0, 0)}
	p := newPins(c)
	p.Put("alice", "a:1")
	p.Put("bob", "a:1")
	c.t = c.t.Add(50 * time.Second)
	p.Put("alice", "a:1") // alice is still active
	c.t = c.t.Add(20 * time.Second)
	if _, ok := p.Get("bob"); ok {
		t.Fatal("bob's pin outlived its TTL")
	}
	if _, ok := p.Get("alice"); !ok {
		t.Fatal("alice's pin expired although she was seen 20s ago")
	}
	if n := p.Sweep(); n != 1 {
		t.Fatalf("%d pins after the sweep, want 1", n)
	}
}

func TestDrop(t *testing.T) {
	p := newPins(&clock{t: time.Unix(0, 0)})
	p.Put("alice", "a:1")
	p.Put("bob", "b:1")
	if n := p.Drop("a:1"); n != 1 {
		t.Fatalf("dropped %d", n)
	}
	if _, ok := p.Get("alice"); ok {
		t.Fatal("alice is still pinned to a dropped backend")
	}
	if p.Len() != 1 {
		t.Fatalf("%d pins left", p.Len())
	}
}
