package metrics

import (
	"strings"
	"sync"
	"testing"
)

func TestCounterText(t *testing.T) {
	var r Registry
	c := r.NewCounter("hits_total", "Hits.", "backend", "status")
	c.Inc("a:1", "200")
	c.Add(2, "a:1", "200")
	c.Inc("b:1", `5"0"0`)
	r.NewGauge("up", "Up.", func() float64 { return 1 })

	var sb strings.Builder
	r.Write(&sb)
	for _, want := range []string{
		`hits_total{backend="a:1",status="200"} 3`,
		`hits_total{backend="b:1",status="5\"0\"0"} 1`,
		"# TYPE up gauge\nup 1",
	} {
		if !strings.Contains(sb.String(), want) {
			t.Errorf("missing %q in\n%s", want, sb.String())
		}
	}
}

func TestCounterIsSafeForConcurrentUse(t *testing.T) {
	var r Registry
	c := r.NewCounter("hits_total", "Hits.", "backend")
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Go(func() {
			for range 100 {
				c.Inc(string(rune('a' + i%5)))
			}
		})
	}
	wg.Wait()
	total := int64(0)
	for _, l := range []string{"a", "b", "c", "d", "e"} {
		total += c.Value(l)
	}
	if total != 5000 {
		t.Fatalf("total %d, want 5000", total)
	}
}
