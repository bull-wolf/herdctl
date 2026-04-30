package retention

import (
	"fmt"
	"sync"
	"time"
)

// Policy defines how long log/metric data is retained for a service.
type Policy struct {
	Service  string
	MaxAge   time.Duration
	MaxItems int
}

// Entry records when a policy was applied.
type Entry struct {
	Service   string
	AppliedAt time.Time
	Purged    int
}

// Manager stores and applies retention policies.
type Manager struct {
	mu       sync.RWMutex
	policies map[string]Policy
	history  []Entry
}

// New returns an initialised Manager.
func New() *Manager {
	return &Manager{
		policies: make(map[string]Policy),
	}
}

// Set registers a retention policy for a service.
func (m *Manager) Set(service string, maxAge time.Duration, maxItems int) error {
	if service == "" {
		return fmt.Errorf("service name must not be empty")
	}
	if maxAge < 0 || maxItems < 0 {
		return fmt.Errorf("maxAge and maxItems must be non-negative")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.policies[service] = Policy{Service: service, MaxAge: maxAge, MaxItems: maxItems}
	return nil
}

// Get returns the policy for a service, and whether it exists.
func (m *Manager) Get(service string) (Policy, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.policies[service]
	return p, ok
}

// Apply records a purge event for a service.
func (m *Manager) Apply(service string, purged int) error {
	if service == "" {
		return fmt.Errorf("service name must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.history = append(m.history, Entry{
		Service:   service,
		AppliedAt: time.Now(),
		Purged:    purged,
	})
	return nil
}

// History returns all recorded purge events.
func (m *Manager) History() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Entry, len(m.history))
	copy(out, m.history)
	return out
}

// All returns all registered policies.
func (m *Manager) All() []Policy {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Policy, 0, len(m.policies))
	for _, p := range m.policies {
		out = append(out, p)
	}
	return out
}
