// Package health checks backends actively.
//
// Each node checks every backend on its own, by calling GET /health on it. A
// backend is marked down after Fall failed checks in a row and up again after
// Rise successful checks in a row. Health is never gossiped: what one node can
// reach, another may not.
package health

import (
	"cmp"
	"context"
	"net/http"
	"slices"
	"sync"
	"time"
)

// Checker keeps the health of a set of backends.
type Checker struct {
	Client *http.Client
	Path   string
	Fall   int
	Rise   int

	mu     sync.Mutex
	states map[string]*state
}

type state struct {
	up        bool
	failures  int
	successes int
	changed   time.Time
}

// Status is a snapshot of one backend's health.
type Status struct {
	Backend string    `json:"backend"`
	Up      bool      `json:"up"`
	Since   time.Time `json:"since"`
}

// New returns a checker with the usual defaults: GET /health, a one second
// timeout, down after 3 failures, up after 2 successes.
func New() *Checker {
	return &Checker{
		Client: &http.Client{Timeout: time.Second},
		Path:   "/health",
		Fall:   3,
		Rise:   2,
		states: make(map[string]*state),
	}
}

// Healthy reports whether backend may take traffic. A backend that has not
// been checked yet counts as healthy, so a new backend takes traffic at once.
func (c *Checker) Healthy(backend string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	s, ok := c.states[backend]
	return !ok || s.up
}

// CheckAll checks each backend once, in parallel, and forgets backends that
// are not in the list any more.
func (c *Checker) CheckAll(ctx context.Context, backends []string) {
	c.mu.Lock()
	for b := range c.states {
		if !slices.Contains(backends, b) {
			delete(c.states, b)
		}
	}
	c.mu.Unlock()

	var wg sync.WaitGroup
	for _, b := range backends {
		wg.Go(func() {
			c.Record(b, c.probe(ctx, b))
		})
	}
	wg.Wait()
}

// Record applies the result of one check.
func (c *Checker) Record(backend string, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	s, found := c.states[backend]
	if !found {
		s = &state{up: true, changed: time.Now()}
		c.states[backend] = s
	}
	if ok {
		s.successes++
		if !s.up && s.successes >= c.Rise {
			s.up = true
			s.successes = 0
			s.changed = time.Now()
		}
		return
	}
	s.failures++
	s.successes = 0
	if s.up && s.failures >= c.Fall {
		s.up = false
		s.failures = 0
		s.changed = time.Now()
	}
}

// Statuses returns the health of every known backend, sorted by name.
func (c *Checker) Statuses() []Status {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Status, 0, len(c.states))
	for b, s := range c.states {
		out = append(out, Status{Backend: b, Up: s.up, Since: s.changed})
	}
	slices.SortFunc(out, func(a, b Status) int { return cmp.Compare(a.Backend, b.Backend) })
	return out
}

// Run checks the backends returned by list every interval until ctx is done.
func (c *Checker) Run(ctx context.Context, interval time.Duration, list func() []string) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		c.CheckAll(ctx, list())
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

func (c *Checker) probe(ctx context.Context, backend string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+backend+c.Path, nil)
	if err != nil {
		return false
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
