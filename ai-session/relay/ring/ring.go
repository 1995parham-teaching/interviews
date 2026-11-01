// Package ring places backends on a consistent-hash ring.
//
// Every backend is hashed onto the ring at several points (its replicas). A
// key belongs to the first point clockwise from the key's own hash. Adding or
// removing a backend only moves the keys next to its points; every other key
// keeps its backend. Two nodes that build a ring from the same backends route
// every key the same way.
package ring

import (
	"crypto/sha256"
	"encoding/binary"
	"slices"
	"sort"
	"strconv"
)

// Replicas is the number of points each backend has on the ring.
const Replicas = 64

type point struct {
	hash    uint64
	backend string
}

// Ring is immutable once built; build a new one when the backends change.
type Ring struct {
	points   []point
	backends []string
}

// New builds a ring over the given backends. The order of the slice does not
// matter.
func New(backends []string) *Ring {
	r := &Ring{backends: slices.Clone(backends)}
	slices.Sort(r.backends)
	r.backends = slices.Compact(r.backends)

	for _, b := range r.backends {
		for i := range Replicas {
			r.points = append(r.points, point{hash: hash(b + "#" + strconv.Itoa(i)), backend: b})
		}
	}
	sort.Slice(r.points, func(i, j int) bool {
		if r.points[i].hash != r.points[j].hash {
			return r.points[i].hash < r.points[j].hash
		}
		return r.points[i].backend < r.points[j].backend
	})
	return r
}

// Backends returns the backends on the ring, sorted.
func (r *Ring) Backends() []string {
	return slices.Clone(r.backends)
}

// Lookup returns the backend that owns key.
//
// usable says whether a backend may take the request (it is healthy, it was
// not tried already). Lookup walks clockwise from the key's hash and returns
// the first usable backend, so a key whose backend is down moves to the next
// backend on the ring and comes back when it recovers. Each backend is tried
// at most once. It reports false when no backend is usable.
func (r *Ring) Lookup(key string, usable func(backend string) bool) (string, bool) {
	if len(r.points) == 0 {
		return "", false
	}
	h := hash(key)
	i := sort.Search(len(r.points), func(i int) bool { return r.points[i].hash >= h })
	for n := range len(r.backends) {
		p := r.points[(i+n)%len(r.points)]
		if usable == nil || usable(p.backend) {
			return p.backend, true
		}
	}
	return "", false
}

// hash is the first eight bytes of SHA-256: slow next to FNV, but it spreads
// similar names such as 10.0.0.1:80#1 and 10.0.0.1:80#2 evenly.
func hash(s string) uint64 {
	sum := sha256.Sum256([]byte(s))
	return binary.BigEndian.Uint64(sum[:8])
}
