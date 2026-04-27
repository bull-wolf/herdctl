package readiness

import (
	"fmt"
	"sync"
	"time"
)

// State represents the readiness state of a service.
type State int

const (
	StateUnknown State = iota
	StateReady
	StateNotReady
)

func (s State) String() string {
	switch s {
	case StateReady:
		return "ready"
	case StateNotReady:
		return "not_ready"
	default:
		return "unknown"
	}
}

// Entry holds readiness info for a single service.
type Entry struct {
	Service   string
	State     State
	Reason    string
	UpdatedAt time.Time
}

// Tracker manages readiness state for all services.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
	}
}

// Set records the readiness state for a service.
func (t *Tracker) Set(service string, state State, reason string) error {
	if service == "" {
		return fmt.Errorf("readiness: service name must not be empty")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[service] = Entry{
		Service:   service,
		State:     state,
		Reason:    reason,
		UpdatedAt: time.Now(),
	}
	return nil
}

// Get returns the readiness entry for a service.
func (t *Tracker) Get(service string) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[service]
	return e, ok
}

// IsReady returns true if the service is in a ready state.
func (t *Tracker) IsReady(service string) bool {
	e, ok := t.Get(service)
	return ok && e.State == StateReady
}

// All returns a snapshot of all readiness entries.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}

// Clear removes the readiness entry for a service.
func (t *Tracker) Clear(service string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, service)
}
