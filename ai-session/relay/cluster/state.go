// Package cluster keeps the list of backends that every node agrees on, and
// spreads it by gossip.
//
// An operator can add, change or remove a backend through the admin API of
// any node. The change is stored with a version from the node's logical clock
// and the node's id as its origin. Every gossip round, nodes exchange their
// whole lists; for each backend the entry with the higher version wins, the
// origin breaking ties, so all nodes end up with the same list. Receiving an
// entry moves the local clock past its version, so a change made anywhere
// always gets a higher version than anything its node has already seen.
//
// A removed backend is kept as a tombstone, an entry with status "removed",
// so that the removal spreads like any other change. Tombstones are forgotten
// after a while (see Sweep).
package cluster

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"relay/config"
)

// Status says whether a backend takes traffic.
type Status string

const (
	Active   Status = "active"   // on the ring
	Draining Status = "draining" // known and checked, but off the ring
	Removed  Status = "removed"  // a tombstone
)

// Entry is what the cluster knows about one backend.
type Entry struct {
	Backend string    `json:"backend"`
	Status  Status    `json:"status"`
	Weight  int       `json:"weight"`
	Version uint64    `json:"version"`
	Origin  string    `json:"origin"` // the node that made this change
	Since   time.Time `json:"since"`  // when the origin made it
}

// newer reports whether a should replace b.
func newer(a, b Entry) bool {
	if a.Version != b.Version {
		return a.Version > b.Version
	}
	return a.Origin > b.Origin
}

// Change is a partial update from the admin API; nil fields keep their value.
type Change struct {
	Status *Status `json:"status"`
	Weight *int    `json:"weight"`
}

// State is one node's copy of the backend list.
type State struct {
	id  string
	ttl time.Duration // how long tombstones are kept
	now func() time.Time

	mu       sync.Mutex
	clock    uint64
	entries  map[string]Entry
	onChange []func()
}

// NewState starts a list from the seed backends, all active at version 0.
// Every node is started with the same seed, so seed entries agree without any
// gossip. Tombstones are kept for ttl.
func NewState(id string, seed []config.Backend, ttl time.Duration) *State {
	s := &State{id: id, ttl: ttl, now: time.Now, entries: make(map[string]Entry)}
	for _, b := range seed {
		s.entries[b.Addr] = Entry{Backend: b.Addr, Status: Active, Weight: b.Weight}
	}
	return s
}

// OnChange registers a function that is called, without the lock held, after
// every change to the list.
func (s *State) OnChange(f func()) {
	s.mu.Lock()
	s.onChange = append(s.onChange, f)
	s.mu.Unlock()
}

// Apply makes a local change to one backend, adding it if it is new or was
// removed.
func (s *State) Apply(backend string, c Change) (Entry, error) {
	if c.Status != nil && *c.Status != Active && *c.Status != Draining {
		return Entry{}, fmt.Errorf("status must be %q or %q", Active, Draining)
	}
	if c.Weight != nil && (*c.Weight < 1 || *c.Weight > config.MaxWeight) {
		return Entry{}, fmt.Errorf("weight must be between 1 and %d", config.MaxWeight)
	}

	s.mu.Lock()
	e, ok := s.entries[backend]
	if !ok || e.Status == Removed {
		e = Entry{Backend: backend, Status: Active}
	}
	if c.Status != nil {
		e.Status = *c.Status
	}
	if c.Weight != nil {
		e.Weight = *c.Weight
	}
	e = s.stamp(e)
	s.mu.Unlock()

	s.changed()
	return e, nil
}

// ErrUnknown is returned when removing a backend the node does not know.
var ErrUnknown = errors.New("unknown backend")

// Remove replaces a backend's entry with a tombstone.
func (s *State) Remove(backend string) (Entry, error) {
	s.mu.Lock()
	e, ok := s.entries[backend]
	if !ok || e.Status == Removed {
		s.mu.Unlock()
		return Entry{}, ErrUnknown
	}
	e.Status = Removed
	e = s.stamp(e)
	s.mu.Unlock()

	s.changed()
	return e, nil
}

// stamp gives e the next version from this node and stores it. s.mu is held.
func (s *State) stamp(e Entry) Entry {
	s.clock++
	e.Version = s.clock
	e.Origin = s.id
	e.Since = s.now()
	s.entries[e.Backend] = e
	return e
}

// Merge applies entries received from a peer and returns how many replaced
// the local ones.
func (s *State) Merge(remote []Entry) int {
	s.mu.Lock()
	n := 0
	for _, e := range remote {
		s.clock = max(s.clock, e.Version)
		if s.expired(e) {
			continue // a peer that has not swept yet; it will
		}
		cur, ok := s.entries[e.Backend]
		if ok && !newer(e, cur) {
			continue
		}
		s.entries[e.Backend] = e
		n++
	}
	s.mu.Unlock()
	if n > 0 {
		s.changed()
	}
	return n
}

// Sweep forgets tombstones older than the state's ttl. By then every node has received
// them, and keeping them for ever would let the list grow without bound as
// backends come and go. The age is counted from the removal itself, which
// travels with the tombstone, so every node forgets it at the same time.
func (s *State) Sweep() int {
	s.mu.Lock()
	n := 0
	for b, e := range s.entries {
		if s.expired(e) {
			delete(s.entries, b)
			n++
		}
	}
	s.mu.Unlock()
	return n
}

// Entries returns a copy of the list, tombstones included, sorted by backend.
func (s *State) Entries() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	slices.SortFunc(out, func(a, b Entry) int { return cmp.Compare(a.Backend, b.Backend) })
	return out
}

// Get returns the entry for one backend, which may be a tombstone.
func (s *State) Get(backend string) (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[backend]
	return e, ok
}

// Known returns the backends that are not removed, sorted.
func (s *State) Known() []string {
	var out []string
	for _, e := range s.Entries() {
		if e.Status != Removed {
			out = append(out, e.Backend)
		}
	}
	return out
}

// Active returns the active entries, sorted by backend.
func (s *State) Active() []Entry {
	var out []Entry
	for _, e := range s.Entries() {
		if e.Status == Active {
			out = append(out, e)
		}
	}
	return out
}

// Clock returns the node's logical clock.
func (s *State) Clock() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clock
}

func (s *State) expired(e Entry) bool {
	return e.Status == Removed && s.now().Sub(e.Since) > s.ttl
}

func (s *State) changed() {
	s.mu.Lock()
	fs := slices.Clone(s.onChange)
	s.mu.Unlock()
	for _, f := range fs {
		f()
	}
}
