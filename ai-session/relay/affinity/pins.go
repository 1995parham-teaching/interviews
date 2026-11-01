// Package affinity keeps sessions where their state is.
//
// The ring already sends a key to the same backend every time, but only for
// as long as the ring does not change. When a backend is added, it takes over
// about 1/n of the keys, and every session among them lands on a backend that
// has never seen it. A pin remembers which backend served a key; while that
// backend is usable, the key stays on it, whatever the ring says. Pins expire
// when a key has not been seen for a while, so new sessions follow the ring.
package affinity

import (
	"sync"
	"time"
)

type pin struct {
	backend string
	seen    time.Time
}

// Pins maps affinity keys to backends.
type Pins struct {
	TTL time.Duration
	now func() time.Time

	mu   sync.Mutex
	pins map[string]pin
}

// New returns an empty set of pins that expire after ttl without use.
func New(ttl time.Duration) *Pins {
	return &Pins{TTL: ttl, now: time.Now, pins: make(map[string]pin)}
}

// Get returns the backend pinned for key, if the pin has not expired.
func (p *Pins) Get(key string) (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	pn, ok := p.pins[key]
	if !ok || p.now().Sub(pn.seen) > p.TTL {
		return "", false
	}
	return pn.backend, true
}

// Put pins key to backend, or refreshes the pin.
func (p *Pins) Put(key, backend string) {
	p.mu.Lock()
	p.pins[key] = pin{backend: backend, seen: p.now()}
	p.mu.Unlock()
}

// Drop removes every pin to backend, for a backend that left the ring.
func (p *Pins) Drop(backend string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	n := 0
	for k, pn := range p.pins {
		if pn.backend == backend {
			delete(p.pins, k)
			n++
		}
	}
	return n
}

// Sweep deletes expired pins and returns how many are left.
func (p *Pins) Sweep() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	for k, pn := range p.pins {
		if now.Sub(pn.seen) > p.TTL {
			delete(p.pins, k)
		}
	}
	return len(p.pins)
}

// Len returns the number of pins, expired ones included.
func (p *Pins) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.pins)
}
