package heartbeat

import (
	"fmt"
	"sync"
	"time"
)

// Entry records the last heartbeat time and interval for a service.
type Entry struct {
	Service  string
	Interval time.Duration
	LastBeat time.Time
	Missed   int
}

// Manager tracks heartbeat state for services.
type Manager struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New returns a new heartbeat Manager.
func New() *Manager {
	return &Manager{
		entries: make(map[string]*Entry),
	}
}

// Register registers a service with the given heartbeat interval.
func (m *Manager) Register(service string, interval time.Duration) error {
	if service == "" {
		return fmt.Errorf("heartbeat: service name must not be empty")
	}
	if interval <= 0 {
		return fmt.Errorf("heartbeat: interval must be positive")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[service] = &Entry{
		Service:  service,
		Interval: interval,
		LastBeat: time.Now(),
	}
	return nil
}

// Beat records a heartbeat for the given service, resetting the missed counter.
func (m *Manager) Beat(service string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[service]
	if !ok {
		return fmt.Errorf("heartbeat: service %q not registered", service)
	}
	e.LastBeat = time.Now()
	e.Missed = 0
	return nil
}

// Check evaluates all services and increments the missed counter for any
// service whose last beat exceeds its interval. Returns a list of service
// names that are considered stale.
func (m *Manager) Check(now time.Time) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	var stale []string
	for _, e := range m.entries {
		if now.Sub(e.LastBeat) > e.Interval {
			e.Missed++
			stale = append(stale, e.Service)
		}
	}
	return stale
}

// Get returns the Entry for the given service, or an error if not found.
func (m *Manager) Get(service string) (Entry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[service]
	if !ok {
		return Entry{}, fmt.Errorf("heartbeat: service %q not found", service)
	}
	return *e, nil
}

// All returns a copy of all registered entries.
func (m *Manager) All() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Entry, 0, len(m.entries))
	for _, e := range m.entries {
		out = append(out, *e)
	}
	return out
}
