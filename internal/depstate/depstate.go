package depstate

import (
	"fmt"
	"sync"
)

// State represents the readiness state of a service's dependencies.
type State string

const (
	StatePending  State = "pending"
	StateReady    State = "ready"
	StateBlocked  State = "blocked"
)

// Entry holds dependency state info for a single service.
type Entry struct {
	Service  string
	State    State
	Blocking []string // services that are not yet ready
}

// Manager tracks whether a service's dependencies are satisfied.
type Manager struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New creates a new depstate Manager.
func New() *Manager {
	return &Manager{
		entries: make(map[string]*Entry),
	}
}

// Evaluate computes the dependency state for a service given a set of
// dependency names and a function that reports whether each dep is ready.
func (m *Manager) Evaluate(service string, deps []string, isReady func(string) bool) State {
	m.mu.Lock()
	defer m.mu.Unlock()

	blocking := []string{}
	for _, dep := range deps {
		if !isReady(dep) {
			blocking = append(blocking, dep)
		}
	}

	state := StateReady
	if len(blocking) > 0 {
		state = StateBlocked
	}

	m.entries[service] = &Entry{
		Service:  service,
		State:    state,
		Blocking: blocking,
	}
	return state
}

// Get returns the last evaluated Entry for a service.
func (m *Manager) Get(service string) (*Entry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	e, ok := m.entries[service]
	if !ok {
		return nil, fmt.Errorf("depstate: no entry for service %q", service)
	}
	return e, nil
}

// All returns a copy of all entries.
func (m *Manager) All() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]Entry, 0, len(m.entries))
	for _, e := range m.entries {
		out = append(out, *e)
	}
	return out
}

// Clear removes the entry for a service.
func (m *Manager) Clear(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, service)
}
