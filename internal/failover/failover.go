package failover

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Strategy defines how failover selects a replacement service.
type Strategy string

const (
	StrategyRoundRobin Strategy = "round-robin"
	StrategyPrimary    Strategy = "primary"
)

// Entry holds failover configuration for a service.
type Entry struct {
	Service   string
	Targets   []string
	Strategy  Strategy
	Active    bool
	FailedAt  time.Time
	Current   string
}

// Manager tracks failover state for services.
type Manager struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	cursors map[string]int
}

// New returns a new failover Manager.
func New() *Manager {
	return &Manager{
		entries: make(map[string]*Entry),
		cursors: make(map[string]int),
	}
}

// Register configures failover targets for a service.
func (m *Manager) Register(service string, targets []string, strategy Strategy) error {
	if service == "" {
		return errors.New("failover: service name required")
	}
	if len(targets) == 0 {
		return errors.New("failover: at least one target required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[service] = &Entry{
		Service:  service,
		Targets:  targets,
		Strategy: strategy,
		Active:   false,
	}
	m.cursors[service] = 0
	return nil
}

// Trigger activates failover for a service and returns the selected target.
func (m *Manager) Trigger(service string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[service]
	if !ok {
		return "", fmt.Errorf("failover: no entry for service %q", service)
	}
	var target string
	switch e.Strategy {
	case StrategyPrimary:
		target = e.Targets[0]
	case StrategyRoundRobin:
		idx := m.cursors[service] % len(e.Targets)
		target = e.Targets[idx]
		m.cursors[service] = idx + 1
	default:
		target = e.Targets[0]
	}
	e.Active = true
	e.FailedAt = time.Now()
	e.Current = target
	return target, nil
}

// Resolve returns the current failover target, or empty string if not active.
func (m *Manager) Resolve(service string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[service]
	if !ok || !e.Active {
		return "", false
	}
	return e.Current, true
}

// Recover clears the active failover state for a service.
func (m *Manager) Recover(service string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[service]
	if !ok {
		return fmt.Errorf("failover: no entry for service %q", service)
	}
	e.Active = false
	e.Current = ""
	m.cursors[service] = 0
	return nil
}

// All returns a snapshot of all registered entries.
func (m *Manager) All() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Entry, 0, len(m.entries))
	for _, e := range m.entries {
		out = append(out, *e)
	}
	return out
}
