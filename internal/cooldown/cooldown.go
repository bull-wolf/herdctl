package cooldown

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds cooldown state for a single service.
type Entry struct {
	Service   string
	Duration  time.Duration
	ActiveAt  time.Time
}

// IsActive returns true if the cooldown period has not yet elapsed.
func (e Entry) IsActive() bool {
	return time.Now().Before(e.ActiveAt.Add(e.Duration))
}

// RemainingTime returns how long is left in the cooldown window.
func (e Entry) RemainingTime() time.Duration {
	remaining := time.Until(e.ActiveAt.Add(e.Duration))
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Manager tracks per-service cooldown windows.
type Manager struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an initialised Manager.
func New() *Manager {
	return &Manager{entries: make(map[string]Entry)}
}

// Set activates a cooldown for the given service starting now.
func (m *Manager) Set(service string, d time.Duration) error {
	if service == "" {
		return fmt.Errorf("cooldown: service name must not be empty")
	}
	if d <= 0 {
		return fmt.Errorf("cooldown: duration must be positive")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[service] = Entry{
		Service:  service,
		Duration: d,
		ActiveAt: time.Now(),
	}
	return nil
}

// IsActive returns true if the named service is currently in a cooldown window.
func (m *Manager) IsActive(service string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[service]
	if !ok {
		return false
	}
	return e.IsActive()
}

// Get returns the Entry for a service and whether it exists.
func (m *Manager) Get(service string) (Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[service]
	return e, ok
}

// Clear removes the cooldown record for a service.
func (m *Manager) Clear(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, service)
}

// All returns a snapshot of every tracked entry.
func (m *Manager) All() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Entry, 0, len(m.entries))
	for _, e := range m.entries {
		out = append(out, e)
	}
	return out
}
