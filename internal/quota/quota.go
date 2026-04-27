package quota

import (
	"errors"
	"fmt"
	"sync"
)

// Kind represents a type of resource quota.
type Kind string

const (
	KindCPU    Kind = "cpu"
	KindMemory Kind = "memory"
	KindProcs  Kind = "procs"
)

// Entry holds the limit and current usage for a service resource.
type Entry struct {
	Service string
	Kind    Kind
	Limit   float64
	Used    float64
}

// Exceeded returns true if usage has reached or surpassed the limit.
func (e Entry) Exceeded() bool {
	return e.Used >= e.Limit
}

// Manager tracks per-service resource quotas.
type Manager struct {
	mu      sync.RWMutex
	entries map[string]map[Kind]*Entry
}

// New creates a new quota Manager.
func New() *Manager {
	return &Manager{
		entries: make(map[string]map[Kind]*Entry),
	}
}

// Set registers or updates a quota limit for a service and kind.
func (m *Manager) Set(service string, kind Kind, limit float64) error {
	if service == "" {
		return errors.New("quota: service name must not be empty")
	}
	if limit <= 0 {
		return fmt.Errorf("quota: limit must be positive, got %.2f", limit)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.entries[service] == nil {
		m.entries[service] = make(map[Kind]*Entry)
	}
	m.entries[service][kind] = &Entry{Service: service, Kind: kind, Limit: limit}
	return nil
}

// Record adds usage against a service quota. Returns an error if exceeded.
func (m *Manager) Record(service string, kind Kind, amount float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	kinds, ok := m.entries[service]
	if !ok {
		return fmt.Errorf("quota: no quota set for service %q", service)
	}
	e, ok := kinds[kind]
	if !ok {
		return fmt.Errorf("quota: no %s quota set for service %q", kind, service)
	}
	e.Used += amount
	if e.Exceeded() {
		return fmt.Errorf("quota: %s limit exceeded for service %q (used=%.2f limit=%.2f)", kind, service, e.Used, e.Limit)
	}
	return nil
}

// Get returns the quota entry for a service and kind.
func (m *Manager) Get(service string, kind Kind) (Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if kinds, ok := m.entries[service]; ok {
		if e, ok := kinds[kind]; ok {
			return *e, true
		}
	}
	return Entry{}, false
}

// Reset clears usage (but not the limit) for a service and kind.
func (m *Manager) Reset(service string, kind Kind) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if kinds, ok := m.entries[service]; ok {
		if e, ok := kinds[kind]; ok {
			e.Used = 0
		}
	}
}

// All returns a flat list of all quota entries.
func (m *Manager) All() []Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []Entry
	for _, kinds := range m.entries {
		for _, e := range kinds {
			out = append(out, *e)
		}
	}
	return out
}
