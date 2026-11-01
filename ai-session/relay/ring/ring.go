// Package ring places backends on a consistent-hash ring.
//
// Every backend is hashed onto the ring at Replicas × weight points. A key
// belongs to the first point clockwise from the key's own hash. Adding or
// removing a backend only moves the keys next to its points, and two nodes
// that build a ring from the same members route every key the same way.
package ring

import (
	"cmp"
	"hash/fnv"
	"slices"
	"sort"
	"strconv"
)

// Replicas is the number of points a backend of weight 1 has on the ring.
const Replicas = 64

// Member is a backend and its weight.
type Member struct {
	Backend string
	Weight  int
}

type point struct {
	hash    uint32
	backend string
}

// Ring is immutable once built; build a new one when the members change.
type Ring struct {
	points  []point
	members []Member
}

// New builds a ring. The order of members does not matter; a backend listed
// twice keeps its first weight.
func New(members []Member) *Ring {
	r := &Ring{}
	seen := map[string]bool{}
	for _, m := range members {
		if seen[m.Backend] {
			continue
		}
		seen[m.Backend] = true
		r.members = append(r.members, m)
	}
	slices.SortFunc(r.members, func(a, b Member) int { return cmp.Compare(a.Backend, b.Backend) })

	for _, m := range r.members {
		// Point names do not depend on the weight, so raising a weight only
		// adds points and only the keys next to the new points move.
		for i := range Replicas * m.Weight {
			r.points = append(r.points, point{hash: hash(m.Backend + "-" + strconv.Itoa(i)), backend: m.Backend})
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

// Members returns the members the ring was built from, sorted by backend.
func (r *Ring) Members() []Member {
	return slices.Clone(r.members)
}

// Lookup returns the backend that owns key, skipping backends that usable
// rejects: it walks clockwise from the key's hash to the first point whose
// backend is usable. It reports false when no backend is.
func (r *Ring) Lookup(key string, usable func(backend string) bool) (string, bool) {
	if len(r.points) == 0 {
		return "", false
	}
	h := hash(key)
	start := sort.Search(len(r.points), func(i int) bool { return r.points[i].hash >= h })
	rejected := map[string]bool{}
	for n := range len(r.points) {
		p := r.points[(start+n)%len(r.points)]
		if rejected[p.backend] {
			continue
		}
		if usable == nil || usable(p.backend) {
			return p.backend, true
		}
		rejected[p.backend] = true
		if len(rejected) == len(r.members) {
			break
		}
	}
	return "", false
}

// hash is 32-bit FNV-1a: fast, and it spreads short strings such as
// addresses well enough for placing points.
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
