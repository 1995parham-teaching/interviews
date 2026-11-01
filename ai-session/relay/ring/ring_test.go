package ring

import (
	"strconv"
	"testing"
)

var three = []string{"10.0.0.1:80", "10.0.0.2:80", "10.0.0.3:80"}

func TestSameKeySameBackend(t *testing.T) {
	a := New(three)
	b := New([]string{three[2], three[0], three[1]})
	for i := range 100 {
		key := "user-" + strconv.Itoa(i)
		x, _ := a.Lookup(key, nil)
		y, _ := b.Lookup(key, nil)
		if x != y {
			t.Fatalf("key %s: %s on one ring, %s on the other", key, x, y)
		}
	}
}

func TestEveryBackendGetsAFairShare(t *testing.T) {
	r := New(three)
	count := map[string]int{}
	const n = 30000
	for i := range n {
		b, ok := r.Lookup("user-"+strconv.Itoa(i), nil)
		if !ok {
			t.Fatal("lookup failed on a full ring")
		}
		count[b]++
	}
	for _, b := range three {
		if share := float64(count[b]) / n; share < 0.25 || share > 0.42 {
			t.Errorf("%s gets %.2f of the keys", b, share)
		}
	}
}

func TestDownBackendIsSkipped(t *testing.T) {
	r := New(three)
	for _, key := range []string{"alice", "bob", "dave"} {
		owner, _ := r.Lookup(key, nil)
		got, ok := r.Lookup(key, func(b string) bool { return b != owner })
		if !ok {
			t.Fatalf("%s: no backend although two are up", key)
		}
		if got == owner {
			t.Fatalf("%s: got the down backend %s", key, owner)
		}
	}
}

func TestRemovingABackendOnlyMovesItsKeys(t *testing.T) {
	before := New(three)
	after := New(three[:2])
	for i := range 1000 {
		key := "user-" + strconv.Itoa(i)
		x, _ := before.Lookup(key, nil)
		y, _ := after.Lookup(key, nil)
		if x != three[2] && x != y {
			t.Fatalf("%s moved from %s to %s", key, x, y)
		}
	}
}

func TestEmptyRing(t *testing.T) {
	if _, ok := New(nil).Lookup("user-1", nil); ok {
		t.Fatal("an empty ring found a backend")
	}
}
