package cluster

import (
	"slices"
	"testing"
	"time"

	"relay/config"
)

var seed = []config.Backend{{Addr: "10.0.0.1:80", Weight: 1}, {Addr: "10.0.0.2:80", Weight: 1}}

func status(s Status) *Status { return &s }
func weight(w int) *int       { return &w }

// fakeClock lets tests move time forward.
type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time          { return c.t }
func (c *fakeClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func newState(id string, clock *fakeClock) *State {
	s := NewState(id, seed, time.Minute)
	if clock != nil {
		s.now = clock.now
	}
	return s
}

func sync2(a, b *State) {
	b.Merge(a.Entries())
	a.Merge(b.Entries())
}

func TestSeedIsActiveAtVersionZero(t *testing.T) {
	s := newState("a", nil)
	for _, e := range s.Entries() {
		if e.Status != Active || e.Version != 0 || e.Weight != 1 {
			t.Fatalf("seed entry %+v", e)
		}
	}
}

func TestApplyStampsVersionAndOrigin(t *testing.T) {
	s := newState("a", nil)
	e1, _ := s.Apply("10.0.0.1:80", Change{Status: status(Draining)})
	e2, _ := s.Apply("10.0.0.1:80", Change{Weight: weight(3)})
	if e1.Version != 1 || e2.Version != 2 || e2.Origin != "a" {
		t.Fatalf("got %+v then %+v", e1, e2)
	}
	if e2.Status != Draining || e2.Weight != 3 {
		t.Fatalf("a partial change lost a field: %+v", e2)
	}
}

func TestApplyAddsABackend(t *testing.T) {
	s := newState("a", nil)
	e, err := s.Apply("10.0.0.3:80", Change{Status: status(Active), Weight: weight(2)})
	if err != nil || e.Weight != 2 {
		t.Fatalf("got %+v, %v", e, err)
	}
	if got := len(s.Active()); got != 3 {
		t.Fatalf("%d active backends, want 3", got)
	}
}

func TestApplyRejectsBadChanges(t *testing.T) {
	s := newState("a", nil)
	for _, c := range []Change{
		{Status: status(Removed)},
		{Status: status("sleeping")},
		{Weight: weight(0)},
		{Weight: weight(11)},
	} {
		if _, err := s.Apply("10.0.0.1:80", c); err == nil {
			t.Errorf("accepted %+v", c)
		}
	}
}

func TestMergeKeepsTheNewerEntry(t *testing.T) {
	s := newState("a", nil)
	s.Apply("10.0.0.1:80", Change{Status: status(Draining)}) // version 1
	if n := s.Merge([]Entry{{Backend: "10.0.0.1:80", Status: Active, Weight: 1, Version: 0}}); n != 0 {
		t.Fatal("an older entry replaced a newer one")
	}
	if n := s.Merge([]Entry{{Backend: "10.0.0.1:80", Status: Active, Weight: 1, Version: 5, Origin: "b"}}); n != 1 {
		t.Fatal("a newer entry was ignored")
	}
}

func TestEqualVersionsAgreeByOrigin(t *testing.T) {
	a, b := newState("a", nil), newState("b", nil)
	a.Apply("10.0.0.1:80", Change{Status: status(Draining)}) // version 1 from a
	b.Apply("10.0.0.1:80", Change{Status: status(Active)})   // version 1 from b
	sync2(a, b)
	ea, _ := a.Get("10.0.0.1:80")
	eb, _ := b.Get("10.0.0.1:80")
	if ea != eb || ea.Origin != "b" {
		t.Fatalf("a has %+v, b has %+v", ea, eb)
	}
}

func TestAChangeAfterAMergeWins(t *testing.T) {
	a, b := newState("a", nil), newState("b", nil)
	a.Apply("10.0.0.1:80", Change{Status: status(Draining)})
	a.Apply("10.0.0.1:80", Change{Status: status(Active)})
	sync2(a, b)
	b.Apply("10.0.0.1:80", Change{Status: status(Draining)})
	sync2(a, b)
	if e, _ := a.Get("10.0.0.1:80"); e.Status != Draining {
		t.Fatalf("a has %+v", e)
	}
}

func TestRemoveSpreadsAndStays(t *testing.T) {
	a, b := newState("a", nil), newState("b", nil)
	a.Remove("10.0.0.2:80")
	for range 3 {
		sync2(a, b)
	}
	for _, s := range []*State{a, b} {
		if slices.Contains(s.Known(), "10.0.0.2:80") {
			t.Fatalf("%s still knows the removed backend", s.id)
		}
	}
	if _, err := a.Remove("10.0.0.2:80"); err != ErrUnknown {
		t.Fatalf("removing twice: %v", err)
	}
}

func TestReAddThroughAnotherNode(t *testing.T) {
	a, b := newState("a", nil), newState("b", nil)
	a.Remove("10.0.0.2:80")
	sync2(a, b)
	b.Apply("10.0.0.2:80", Change{Status: status(Active), Weight: weight(1)})
	sync2(a, b)
	if e, _ := a.Get("10.0.0.2:80"); e.Status != Active {
		t.Fatalf("a has %+v", e)
	}
}

func TestSweepForgetsOldTombstones(t *testing.T) {
	clock := &fakeClock{t: time.Unix(1000, 0)}
	s := newState("a", clock)
	s.Remove("10.0.0.2:80")
	if n := s.Sweep(); n != 0 {
		t.Fatal("swept a fresh tombstone")
	}
	old := s.Entries()
	clock.advance(2 * time.Minute)
	if n := s.Sweep(); n != 1 {
		t.Fatal("kept an old tombstone")
	}
	if _, ok := s.Get("10.0.0.2:80"); ok {
		t.Fatal("tombstone still there")
	}
	if n := s.Sweep(); n != 0 {
		t.Fatal("swept something that was not a tombstone")
	}
	if n := s.Merge(old); n != 0 {
		t.Fatal("took back an expired tombstone from a peer that had not swept")
	}
}

func TestOnChange(t *testing.T) {
	s := newState("a", nil)
	calls := 0
	s.OnChange(func() { calls++ })
	s.Apply("10.0.0.1:80", Change{Status: status(Draining)})
	s.Merge([]Entry{{Backend: "10.0.0.1:80", Status: Active, Weight: 1}}) // older
	s.Merge([]Entry{{Backend: "10.0.0.3:80", Status: Active, Weight: 1, Version: 4, Origin: "b"}})
	s.Remove("10.0.0.3:80")
	if calls != 3 {
		t.Fatalf("OnChange called %d times, want 3", calls)
	}
}
