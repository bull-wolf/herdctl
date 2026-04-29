package limiter

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds concurrency limit config and current usage for a service.
type Entry struct {
	Service   string
	Max       int
	Active    int
	UpdatedAt time.Time
}

// Limiter enforces max-concurrency limits per service.
type Limiter struct {
	mu      sync.Mutex
	entries map[string]*Entry
}

// New creates a new Limiter.
func New() *Limiter {
	return &Limiter{entries: make(map[string]*Entry)}
}

// Set configures the max concurrency for a service.
func (l *Limiter) Set(service string, max int) error {
	if service == "" {
		return fmt.Errorf("service name must not be empty")
	}
	if max <= 0 {
		return fmt.Errorf("max must be greater than zero")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.entries[service]
	if !ok {
		e = &Entry{Service: service}
		l.entries[service] = e
	}
	e.Max = max
	e.UpdatedAt = time.Now()
	return nil
}

// Acquire attempts to acquire a concurrency slot for the service.
// Returns an error if the limit is exceeded or no limit is configured.
func (l *Limiter) Acquire(service string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.entries[service]
	if !ok {
		return fmt.Errorf("no limit configured for service %q", service)
	}
	if e.Active >= e.Max {
		return fmt.Errorf("concurrency limit reached for service %q (%d/%d)", service, e.Active, e.Max)
	}
	e.Active++
	e.UpdatedAt = time.Now()
	return nil
}

// Release decrements the active count for the service.
func (l *Limiter) Release(service string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.entries[service]
	if !ok {
		return fmt.Errorf("no limit configured for service %q", service)
	}
	if e.Active > 0 {
		e.Active--
		e.UpdatedAt = time.Now()
	}
	return nil
}

// Get returns the entry for a service.
func (l *Limiter) Get(service string) (Entry, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.entries[service]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a copy of all entries.
func (l *Limiter) All() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, 0, len(l.entries))
	for _, e := range l.entries {
		out = append(out, *e)
	}
	return out
}
