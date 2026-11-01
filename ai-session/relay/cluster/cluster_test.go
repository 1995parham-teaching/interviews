package cluster

import (
	"slices"
	"testing"
)

func TestSeedIsActiveAtVersionZero(t *testing.T) {
	s := NewState([]string{"b:2", "a:1"})
	want := []Entry{{"a:1", Active, 0}, {"b:2", Active, 0}}
	if got := s.Entries(); !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSetBumpsTheVersion(t *testing.T) {
	s := NewState([]string{"a:1"})
	e1, _ := s.Set("a:1", Draining)
	e2, _ := s.Set("a:1", Active)
	if e1.Version != 1 || e2.Version != 2 {
		t.Fatalf("versions %d and %d, want 1 and 2", e1.Version, e2.Version)
	}
	if _, err := s.Set("a:1", "sleeping"); err == nil {
		t.Fatal("accepted an unknown status")
	}
}

func TestMergeKeepsTheHigherVersion(t *testing.T) {
	s := NewState([]string{"a:1"})
	s.Set("a:1", Draining) // version 1

	if n := s.Merge([]Entry{{"a:1", Active, 0}}); n != 0 {
		t.Fatal("an older entry replaced a newer one")
	}
	if n := s.Merge([]Entry{{"a:1", Active, 5}, {"c:3", Active, 2}}); n != 2 {
		t.Fatalf("changed %d entries, want 2", n)
	}
	if e, _ := s.Get("a:1"); e.Status != Active || e.Version != 5 {
		t.Fatalf("a:1 is %v", e)
	}
}

func TestTwoStatesConverge(t *testing.T) {
	a := NewState([]string{"x:1", "y:2"})
	b := NewState([]string{"x:1", "y:2"})
	a.Set("y:2", Draining)
	a.Set("z:3", Active)

	b.Merge(a.Entries())
	a.Merge(b.Entries())
	if !slices.Equal(a.Entries(), b.Entries()) {
		t.Fatalf("a %v, b %v", a.Entries(), b.Entries())
	}
	if got := b.ActiveBackends(); !slices.Equal(got, []string{"x:1", "z:3"}) {
		t.Fatalf("active on b: %v", got)
	}
}

func TestRemove(t *testing.T) {
	s := NewState([]string{"a:1", "b:2"})
	if !s.Remove("a:1") {
		t.Fatal("could not remove a known backend")
	}
	if s.Remove("a:1") {
		t.Fatal("removed a backend twice")
	}
	if got := s.Backends(); !slices.Equal(got, []string{"b:2"}) {
		t.Fatalf("backends %v", got)
	}
}

func TestOnChange(t *testing.T) {
	s := NewState(nil)
	calls := 0
	s.OnChange(func() { calls++ })
	s.Set("a:1", Active)
	s.Merge([]Entry{{"a:1", Active, 0}}) // older, no change
	s.Merge([]Entry{{"b:2", Active, 1}})
	s.Remove("a:1")
	if calls != 3 {
		t.Fatalf("OnChange called %d times, want 3", calls)
	}
}
