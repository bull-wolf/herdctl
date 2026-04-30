package priority

import (
	"fmt"
	"sync"
)

// Level represents a service startup/execution priority.
type Level int

const (
	Low    Level = 10
	Normal Level = 50
	High   Level = 100
)

// Entry holds the priority configuration for a service.
type Entry struct {
	Service string
	Level   Level
}

// Manager stores and retrieves priority levels for services.
type Manager struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New creates a new priority Manager.
func New() *Manager {
	return &Manager{
		entries: make(map[string]Entry),
	}
}

// Set assigns a priority level to a service.
func (m *Manager) Set(service string, level Level) error {
	if service == "" {
		return fmt.Errorf("priority: service name must not be empty")
	}
	if level < 0 {
		return fmt.Errorf("priority: level must be non-negative")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[service] = Entry{Service: service, Level: level}
	return nil
}

// Get returns the priority entry for a service.
func (m *Manager) Get(service string) (Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[service]
	return e, ok
}

// Level returns the numeric priority for a service, defaulting to Normal.
func (m *Manager) GetLevel(service string) Level {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if e, ok := m.entries[service]; ok {
		return e.Level
	}
	return Normal
}

// All returns all priority entries.
func (m *Manager) All() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Entry, 0, len(m.entries))
	for _, e := range m.entries {
		out = append(out, e)
	}
	return out
}

// Sorted returns service names ordered from highest to lowest priority.
func (m *Manager) Sorted(services []string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	copy := append([]string(nil), services...)
	// simple insertion sort — service lists are typically small
	for i := 1; i < len(copy); i++ {
		for j := i; j > 0; j-- {
			a := m.levelOf(copy[j-1])
			b := m.levelOf(copy[j])
			if b > a {
				copy[j-1], copy[j] = copy[j], copy[j-1]
			} else {
				break
			}
		}
	}
	return copy
}

func (m *Manager) levelOf(service string) Level {
	if e, ok := m.entries[service]; ok {
		return e.Level
	}
	return Normal
}
