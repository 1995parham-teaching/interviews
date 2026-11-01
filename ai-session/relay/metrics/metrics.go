// Package metrics counts what a node does and serves it in the Prometheus
// text format.
package metrics

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

// Counter is a family of counters that share a name and differ by labels.
type Counter struct {
	name   string
	help   string
	labels []string

	mu     sync.RWMutex
	values map[string]*atomic.Int64 // key: label values joined by \xff
}

// Gauge is a single value read when metrics are served.
type Gauge struct {
	name string
	help string
	read func() float64
}

// Registry holds every metric of a node.
type Registry struct {
	mu       sync.Mutex
	counters []*Counter
	gauges   []*Gauge
}

// NewCounter registers a counter family.
func (r *Registry) NewCounter(name, help string, labels ...string) *Counter {
	c := &Counter{name: name, help: help, labels: labels, values: map[string]*atomic.Int64{}}
	r.mu.Lock()
	r.counters = append(r.counters, c)
	r.mu.Unlock()
	return c
}

// NewGauge registers a gauge whose value is read by f.
func (r *Registry) NewGauge(name, help string, f func() float64) {
	r.mu.Lock()
	r.gauges = append(r.gauges, &Gauge{name: name, help: help, read: f})
	r.mu.Unlock()
}

// Add adds n to the counter with the given label values.
func (c *Counter) Add(n int64, values ...string) {
	if len(values) != len(c.labels) {
		panic(fmt.Sprintf("metrics: %s takes %d labels, got %d", c.name, len(c.labels), len(values)))
	}
	key := strings.Join(values, "\xff")
	c.mu.RLock()
	v, ok := c.values[key]
	c.mu.RUnlock()
	if !ok {
		c.mu.Lock()
		if v, ok = c.values[key]; !ok {
			v = new(atomic.Int64)
			c.values[key] = v
		}
		c.mu.Unlock()
	}
	v.Add(n)
}

// Inc adds one.
func (c *Counter) Inc(values ...string) { c.Add(1, values...) }

// Value returns the counter with the given label values.
func (c *Counter) Value(values ...string) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if v, ok := c.values[strings.Join(values, "\xff")]; ok {
		return v.Load()
	}
	return 0
}

// Write writes every metric in the text format.
func (r *Registry) Write(w io.Writer) {
	r.mu.Lock()
	counters := slices.Clone(r.counters)
	gauges := slices.Clone(r.gauges)
	r.mu.Unlock()

	for _, c := range counters {
		fmt.Fprintf(w, "# HELP %s %s\n# TYPE %s counter\n", c.name, c.help, c.name)
		c.mu.RLock()
		keys := make([]string, 0, len(c.values))
		for k := range c.values {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "%s%s %d\n", c.name, labelString(c.labels, strings.Split(k, "\xff")), c.values[k].Load())
		}
		c.mu.RUnlock()
	}
	for _, g := range gauges {
		fmt.Fprintf(w, "# HELP %s %s\n# TYPE %s gauge\n%s %g\n", g.name, g.help, g.name, g.name, g.read())
	}
}

// Handler serves the registry.
func (r *Registry) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		r.Write(w)
	})
}

func labelString(names, values []string) string {
	if len(names) == 0 {
		return ""
	}
	escape := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`)
	parts := make([]string, len(names))
	for i, n := range names {
		parts[i] = n + `="` + escape.Replace(values[i]) + `"`
	}
	return "{" + strings.Join(parts, ",") + "}"
}
