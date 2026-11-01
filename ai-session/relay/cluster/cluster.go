// Package cluster keeps the list of backends that every node agrees on.
//
// An operator can add, drain or remove a backend through the admin API of any
// node. The change is recorded locally with a new version and spreads to the
// other nodes by gossip: every node regularly sends its whole list to each
// peer, and the peer answers with its own. Both sides keep, for each backend,
// the entry with the higher version, so all nodes end up with the same list.
//
// Versions come from a logical clock that each node keeps. Wall-clock time is
// never compared between nodes.
package cluster

import (
	"cmp"
	"fmt"
	"slices"
	"sync"
)

// Status says whether a backend takes new traffic.
type Status string

const (
	// Active backends are on the ring.
	Active Status = "active"
	// Draining backends stay known and health-checked but take no traffic.
	Draining Status = "draining"
)

// Entry is what the cluster knows about one backend.
type Entry struct {
	Backend string `json:"backend"`
	Status  Status `json:"status"`
	Version uint64 `json:"version"`
}

// State is one node's copy of the backend list.
type State struct {
	mu       sync.Mutex
	clock    uint64
	entries  map[string]Entry
	onChange func()
}

// NewState starts a list with the given backends, all active, at version 0.
// Every node is started with the same seed list, so seed entries agree
// without any gossip.
func NewState(seed []string) *State {
	s := &State{entries: make(map[string]Entry)}
	for _, b := range seed {
		s.entries[b] = Entry{Backend: b, Status: Active}
	}
	return s
}

// OnChange registers a function that is called, without the lock held, after
// every change to the list.
func (s *State) OnChange(f func()) {
	s.mu.Lock()
	s.onChange = f
	s.mu.Unlock()
}

// Set adds a backend or changes its status, as a new local change.
func (s *State) Set(backend string, status Status) (Entry, error) {
	if status != Active && status != Draining {
		return Entry{}, fmt.Errorf("unknown status %q", status)
	}
	s.mu.Lock()
	s.clock++
	e := Entry{Backend: backend, Status: status, Version: s.clock}
	s.entries[backend] = e
	s.mu.Unlock()
	s.changed()
	return e, nil
}

// Remove forgets a backend. It reports false if the backend was not known.
func (s *State) Remove(backend string) bool {
	s.mu.Lock()
	_, ok := s.entries[backend]
	delete(s.entries, backend)
	s.mu.Unlock()
	if ok {
		s.changed()
	}
	return ok
}

// Merge applies entries received from a peer and returns how many changed
// the local list. For each backend the higher version wins.
func (s *State) Merge(remote []Entry) int {
	s.mu.Lock()
	n := 0
	for _, e := range remote {
		cur, ok := s.entries[e.Backend]
		if !ok || e.Version > cur.Version {
			s.entries[e.Backend] = e
			n++
		}
	}
	s.mu.Unlock()
	if n > 0 {
		s.changed()
	}
	return n
}

// Entries returns a copy of the list, sorted by backend.
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

// Get returns the entry for one backend.
func (s *State) Get(backend string) (Entry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[backend]
	return e, ok
}

// Backends returns every known backend, whatever its status, sorted.
func (s *State) Backends() []string {
	var out []string
	for _, e := range s.Entries() {
		out = append(out, e.Backend)
	}
	return out
}

// ActiveBackends returns the backends that take traffic, sorted.
func (s *State) ActiveBackends() []string {
	var out []string
	for _, e := range s.Entries() {
		if e.Status == Active {
			out = append(out, e.Backend)
		}
	}
	return out
}

func (s *State) changed() {
	s.mu.Lock()
	f := s.onChange
	s.mu.Unlock()
	if f != nil {
		f()
	}
}
