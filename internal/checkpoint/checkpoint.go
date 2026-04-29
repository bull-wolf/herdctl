package checkpoint

import (
	"fmt"
	"sync"
	"time"
)

// Status represents the state of a checkpoint.
type Status string

const (
	StatusPending Status = "pending"
	StatusPassed  Status = "passed"
	StatusFailed  Status = "failed"
)

// Entry holds checkpoint data for a service.
type Entry struct {
	Service   string
	Name      string
	Status    Status
	Message   string
	UpdatedAt time.Time
}

// Manager tracks named checkpoints per service.
type Manager struct {
	mu      sync.RWMutex
	entries map[string]map[string]*Entry // service -> name -> entry
}

// New creates a new checkpoint Manager.
func New() *Manager {
	return &Manager{
		entries: make(map[string]map[string]*Entry),
	}
}

// Set records or updates a checkpoint for a service.
func (m *Manager) Set(service, name string, status Status, message string) error {
	if service == "" || name == "" {
		return fmt.Errorf("checkpoint: service and name must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.entries[service]; !ok {
		m.entries[service] = make(map[string]*Entry)
	}
	m.entries[service][name] = &Entry{
		Service:   service,
		Name:      name,
		Status:    status,
		Message:   message,
		UpdatedAt: time.Now(),
	}
	return nil
}

// Get returns the checkpoint entry for a service and name.
func (m *Manager) Get(service, name string) (*Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if svc, ok := m.entries[service]; ok {
		if e, ok := svc[name]; ok {
			copy := *e
			return &copy, true
		}
	}
	return nil, false
}

// AllForService returns all checkpoints for a given service.
func (m *Manager) AllForService(service string) []*Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []*Entry
	for _, e := range m.entries[service] {
		copy := *e
		out = append(out, &copy)
	}
	return out
}

// Clear removes all checkpoints for a service.
func (m *Manager) Clear(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, service)
}

// AllPassed returns true if every checkpoint for a service has passed.
func (m *Manager) AllPassed(service string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	svc, ok := m.entries[service]
	if !ok || len(svc) == 0 {
		return false
	}
	for _, e := range svc {
		if e.Status != StatusPassed {
			return false
		}
	}
	return true
}
