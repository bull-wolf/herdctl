package throttle

import (
	"fmt"
	"sync"
	"time"
)

// Entry tracks restart attempts for a single service.
type Entry struct {
	Attempts  int
	LastReset time.Time
}

// Throttle enforces a maximum number of restarts per service within a rolling window.
type Throttle struct {
	mu      sync.Mutex
	entries map[string]*Entry
	max     int
	window  time.Duration
}

// New creates a Throttle that allows at most maxAttempts restarts within window.
func New(maxAttempts int, window time.Duration) *Throttle {
	return &Throttle{
		entries: make(map[string]*Entry),
		max:     maxAttempts,
		window:  window,
	}
}

// Allow reports whether the service is permitted to restart.
// It increments the attempt counter; callers should check before restarting.
func (t *Throttle) Allow(service string) (bool, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	e, ok := t.entries[service]
	if !ok || now.Sub(e.LastReset) >= t.window {
		t.entries[service] = &Entry{Attempts: 1, LastReset: now}
		return true, nil
	}

	if e.Attempts >= t.max {
		return false, fmt.Errorf(
			"throttle: service %q exceeded %d restarts within %s",
			service, t.max, t.window,
		)
	}

	e.Attempts++
	return true, nil
}

// Reset clears the attempt counter for a service (e.g. after a clean run).
func (t *Throttle) Reset(service string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, service)
}

// Get returns the current entry for a service, or nil if not tracked.
func (t *Throttle) Get(service string) *Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[service]
	if !ok {
		return nil
	}
	copy := *e
	return &copy
}
