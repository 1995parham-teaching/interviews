package ring

import (
	"strconv"
	"testing"
)

func members(weights ...int) []Member {
	var out []Member
	for i, w := range weights {
		out = append(out, Member{Backend: "10.0.0." + strconv.Itoa(i+1) + ":80", Weight: w})
	}
	return out
}

func share(r *Ring, n int) map[string]float64 {
	count := map[string]int{}
	for i := range n {
		b, _ := r.Lookup("user-"+strconv.Itoa(i), nil)
		count[b]++
	}
	out := map[string]float64{}
	for b, c := range count {
		out[b] = float64(c) / float64(n)
	}
	return out
}

func TestSameMembersSameRouting(t *testing.T) {
	a := New(members(1, 1, 1))
	m := members(1, 1, 1)
	b := New([]Member{m[2], m[0], m[1]})
	for i := range 200 {
		key := "user-" + strconv.Itoa(i)
		x, _ := a.Lookup(key, nil)
		y, _ := b.Lookup(key, nil)
		if x != y {
			t.Fatalf("key %s: %s on one ring, %s on the other", key, x, y)
		}
	}
}

func TestEveryBackendGetsKeys(t *testing.T) {
	r := New(members(1, 1, 1))
	s := share(r, 20000)
	for _, m := range r.Members() {
		if s[m.Backend] < 0.1 {
			t.Errorf("%s gets %.2f of the keys", m.Backend, s[m.Backend])
		}
	}
}

func TestHeavierBackendGetsMore(t *testing.T) {
	r := New(members(1, 4))
	s := share(r, 20000)
	if s["10.0.0.2:80"] <= s["10.0.0.1:80"] {
		t.Fatalf("weight 4 gets %.2f, weight 1 gets %.2f", s["10.0.0.2:80"], s["10.0.0.1:80"])
	}
}

func TestUnusableBackendsAreSkipped(t *testing.T) {
	r := New(members(1, 1, 1))
	for i := range 2000 {
		key := "user-" + strconv.Itoa(i)
		owner, _ := r.Lookup(key, nil)
		got, ok := r.Lookup(key, func(b string) bool { return b != owner })
		if !ok || got == owner {
			t.Fatalf("%s: got %q, %v with %s down", key, got, ok, owner)
		}
	}
	if _, ok := r.Lookup("user-1", func(string) bool { return false }); ok {
		t.Fatal("found a backend although none is usable")
	}
}

func TestRemovingABackendOnlyMovesItsKeys(t *testing.T) {
	m := members(1, 1, 1)
	before, after := New(m), New(m[:2])
	for i := range 2000 {
		key := "user-" + strconv.Itoa(i)
		x, _ := before.Lookup(key, nil)
		y, _ := after.Lookup(key, nil)
		if x != m[2].Backend && x != y {
			t.Fatalf("%s moved from %s to %s", key, x, y)
		}
	}
}

func TestRaisingAWeightOnlyMovesKeysToThatBackend(t *testing.T) {
	before := New(members(1, 1, 1))
	after := New(members(1, 3, 1))
	for i := range 2000 {
		key := "user-" + strconv.Itoa(i)
		x, _ := before.Lookup(key, nil)
		y, _ := after.Lookup(key, nil)
		if x != y && y != "10.0.0.2:80" {
			t.Fatalf("%s moved from %s to %s", key, x, y)
		}
	}
}

func TestEmptyRing(t *testing.T) {
	if _, ok := New(nil).Lookup("user-1", nil); ok {
		t.Fatal("an empty ring found a backend")
	}
}
