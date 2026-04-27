package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Entry tracks rate limit state for a single service.
type Entry struct {
	Service   string
	Requests  int
	WindowEnd time.Time
	Limit     int
}

// Remaining returns the number of requests still allowed in the current window.
func (e *Entry) Remaining() int {
	if e.Requests >= e.Limit {
		return 0
	}
	return e.Limit - e.Requests
}

// Limiter enforces per-service request rate limits over a sliding window.
type Limiter struct {
	mu      sync.Mutex
	entries map[string]*Entry
	window  time.Duration
}

// New creates a Limiter with the given window duration.
func New(window time.Duration) *Limiter {
	return &Limiter{
		entries: make(map[string]*Entry),
		window:  window,
	}
}

// Allow checks whether a request for the given service is within its limit.
// Returns true if allowed, false if the limit has been exceeded.
func (l *Limiter) Allow(service string, limit int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	e, ok := l.entries[service]
	if !ok || now.After(e.WindowEnd) {
		l.entries[service] = &Entry{
			Service:   service,
			Requests:  1,
			WindowEnd: now.Add(l.window),
			Limit:     limit,
		}
		return true
	}

	e.Limit = limit
	if e.Requests >= limit {
		return false
	}
	e.Requests++
	return true
}

// Get returns the current rate limit entry for a service.
func (l *Limiter) Get(service string) (*Entry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[service]
	if !ok {
		return nil, fmt.Errorf("ratelimit: no entry for service %q", service)
	}
	copy := *e
	return &copy, nil
}

// Reset clears the rate limit state for a service.
func (l *Limiter) Reset(service string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, service)
}

// All returns a snapshot of all current rate limit entries.
func (l *Limiter) All() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := make([]Entry, 0, len(l.entries))
	for _, e := range l.entries {
		out = append(out, *e)
	}
	return out
}
