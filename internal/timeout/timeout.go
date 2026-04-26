package timeout

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds timeout configuration and state for a service.
type Entry struct {
	Service  string
	Duration time.Duration
	Deadline time.Time
	Expired  bool
}

// Manager tracks per-service timeouts.
type Manager struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New creates a new timeout Manager.
func New() *Manager {
	return &Manager{
		entries: make(map[string]*Entry),
	}
}

// Set registers a timeout for the given service starting from now.
func (m *Manager) Set(service string, d time.Duration) error {
	if service == "" {
		return fmt.Errorf("service name must not be empty")
	}
	if d <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[service] = &Entry{
		Service:  service,
		Duration: d,
		Deadline: time.Now().Add(d),
		Expired:  false,
	}
	return nil
}

// Check evaluates whether the service timeout has elapsed.
// It marks the entry as expired if the deadline has passed.
func (m *Manager) Check(service string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[service]
	if !ok {
		return false, fmt.Errorf("no timeout set for service %q", service)
	}
	if !e.Expired && time.Now().After(e.Deadline) {
		e.Expired = true
	}
	return e.Expired, nil
}

// Get returns the Entry for a service without modifying state.
func (m *Manager) Get(service string) (*Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[service]
	if !ok {
		return nil, false
	}
	copy := *e
	return &copy, true
}

// Clear removes the timeout entry for a service.
func (m *Manager) Clear(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, service)
}

// All returns a snapshot of all current entries.
func (m *Manager) All() []*Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*Entry, 0, len(m.entries))
	for _, e := range m.entries {
		copy := *e
		result = append(result, &copy)
	}
	return result
}
