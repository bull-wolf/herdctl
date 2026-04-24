package restart

import (
	"fmt"
	"sync"
	"time"
)

// Policy defines the restart behaviour for a service.
type Policy string

const (
	PolicyNever  Policy = "never"
	PolicyAlways Policy = "always"
	PolicyOnFail Policy = "on-failure"
)

// Entry holds restart metadata for a single service.
type Entry struct {
	Policy   Policy
	Attempts int
	MaxRetry int
	LastAt   time.Time
	Cooldown time.Duration
}

// Manager tracks restart state for services.
type Manager struct {
	mu      sync.Mutex
	entries map[string]*Entry
}

// New creates a new restart Manager.
func New() *Manager {
	return &Manager{entries: make(map[string]*Entry)}
}

// Register configures a restart policy for a service.
func (m *Manager) Register(service string, policy Policy, maxRetry int, cooldown time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[service] = &Entry{
		Policy:   policy,
		MaxRetry: maxRetry,
		Cooldown: cooldown,
	}
}

// ShouldRestart returns true when the service is eligible for a restart
// given whether the last exit was a failure.
func (m *Manager) ShouldRestart(service string, failed bool) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[service]
	if !ok {
		return false, fmt.Errorf("restart: unknown service %q", service)
	}
	switch e.Policy {
	case PolicyNever:
		return false, nil
	case PolicyAlways:
		// fall through
	case PolicyOnFail:
		if !failed {
			return false, nil
		}
	}
	if e.MaxRetry > 0 && e.Attempts >= e.MaxRetry {
		return false, nil
	}
	if time.Since(e.LastAt) < e.Cooldown {
		return false, nil
	}
	return true, nil
}

// Record increments the attempt counter and updates the last-restart timestamp.
func (m *Manager) Record(service string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[service]
	if !ok {
		return fmt.Errorf("restart: unknown service %q", service)
	}
	e.Attempts++
	e.LastAt = time.Now()
	return nil
}

// Reset clears the attempt counter for a service (e.g. after a clean start).
func (m *Manager) Reset(service string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[service]
	if !ok {
		return fmt.Errorf("restart: unknown service %q", service)
	}
	e.Attempts = 0
	e.LastAt = time.Time{}
	return nil
}

// Get returns the Entry for a service.
func (m *Manager) Get(service string) (Entry, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[service]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}
