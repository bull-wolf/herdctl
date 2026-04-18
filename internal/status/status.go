package status

import (
	"sync"
	"time"
)

// State represents the lifecycle state of a service.
type State string

const (
	StateStopped  State = "stopped"
	StateStarting State = "starting"
	StateRunning  State = "running"
	StateFailed   State = "failed"
)

// Entry holds the current status of a single service.
type Entry struct {
	Service   string
	State     State
	PID       int
	UpdatedAt time.Time
}

// Tracker stores and retrieves service status entries.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{entries: make(map[string]Entry)}
}

// Set updates the status for the given service.
func (t *Tracker) Set(service string, state State, pid int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[service] = Entry{
		Service:   service,
		State:     state,
		PID:       pid,
		UpdatedAt: time.Now(),
	}
}

// Get returns the status entry for the given service and whether it exists.
func (t *Tracker) Get(service string) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[service]
	return e, ok
}

// All returns a snapshot of all service status entries.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}
