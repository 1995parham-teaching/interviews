// Package health checks backends actively.
//
// Each node checks every known backend on its own with GET /health. A backend
// is marked down after Fall failed checks in a row and up again after Rise
// successful checks in a row. Health is never gossiped: what one node can
// reach, another may not.
package health

import (
	"cmp"
	"context"
	"io"
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
	onFlip []func(backend string, up bool)
}

type state struct {
	up        bool
	failures  int // in a row
	successes int // in a row
	since     time.Time
	lastErr   string
}

// Status is a snapshot of one backend's health.
type Status struct {
	Backend string    `json:"backend"`
	Up      bool      `json:"up"`
	Since   time.Time `json:"since"`
	LastErr string    `json:"last_error,omitempty"`
}

// New returns a checker that probes GET /health with the given timeout.
func New(timeout time.Duration, fall, rise int) *Checker {
	return &Checker{
		Client: &http.Client{Timeout: timeout},
		Path:   "/health",
		Fall:   fall,
		Rise:   rise,
		states: make(map[string]*state),
	}
}

// OnFlip registers a function called, without the lock, when a backend goes
// down or comes back up.
func (c *Checker) OnFlip(f func(backend string, up bool)) {
	c.mu.Lock()
	c.onFlip = append(c.onFlip, f)
	c.mu.Unlock()
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
			err := c.probe(ctx, b)
			c.Record(b, err)
		})
	}
	wg.Wait()
}

// Record applies the result of one check; err is nil for a success.
func (c *Checker) Record(backend string, err error) {
	c.mu.Lock()
	s, found := c.states[backend]
	if !found {
		s = &state{up: true, since: time.Now()}
		c.states[backend] = s
	}
	flipped := false
	if err == nil {
		s.failures = 0
		s.successes++
		s.lastErr = ""
		if !s.up && s.successes >= c.Rise {
			s.up, s.since, flipped = true, time.Now(), true
		}
	} else {
		s.successes = 0
		s.failures++
		s.lastErr = err.Error()
		if s.up && s.failures >= c.Fall {
			s.up, s.since, flipped = false, time.Now(), true
		}
	}
	up := s.up
	fs := slices.Clone(c.onFlip)
	c.mu.Unlock()

	if flipped {
		for _, f := range fs {
			f(backend, up)
		}
	}
}

// Statuses returns the health of every known backend, sorted by name.
func (c *Checker) Statuses() []Status {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Status, 0, len(c.states))
	for b, s := range c.states {
		out = append(out, Status{Backend: b, Up: s.up, Since: s.since, LastErr: s.lastErr})
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

type statusError int

func (e statusError) Error() string { return "health answered " + http.StatusText(int(e)) }

func (c *Checker) probe(ctx context.Context, backend string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+backend+c.Path, nil)
	if err != nil {
		return err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	// Drain a little so the connection can be reused for the next check.
	io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return statusError(resp.StatusCode)
	}
	return nil
}
