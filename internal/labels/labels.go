package labels

import "sync"

// Store manages arbitrary key-value labels attached to services.
type Store struct {
	mu     sync.RWMutex
	entries map[string]map[string]string
}

// New returns an initialised label Store.
func New() *Store {
	return &Store{
		entries: make(map[string]map[string]string),
	}
}

// Set attaches a label key=value to the given service.
func (s *Store) Set(service, key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.entries[service] == nil {
		s.entries[service] = make(map[string]string)
	}
	s.entries[service][key] = value
}

// Get returns the value for a label key on the given service.
// The second return value is false when the key does not exist.
func (s *Store) Get(service, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if m, ok := s.entries[service]; ok {
		v, found := m[key]
		return v, found
	}
	return "", false
}

// All returns a copy of all labels for a service.
// Returns nil when the service has no labels.
func (s *Store) All(service string) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.entries[service]
	if !ok {
		return nil
	}
	copy := make(map[string]string, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy
}

// Delete removes a single label from a service.
func (s *Store) Delete(service, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok := s.entries[service]; ok {
		delete(m, key)
	}
}

// Clear removes all labels for a service.
func (s *Store) Clear(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, service)
}
