package pinned

import (
	"fmt"
	"sync"
	"time"
)

// Entry represents a pinned version for a service.
type Entry struct {
	Service   string
	Version   string
	PinnedAt  time.Time
	PinnedBy  string
	Reason    string
}

// Store tracks pinned versions per service.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New creates a new pinned version Store.
func New() *Store {
	return &Store{
		entries: make(map[string]Entry),
	}
}

// Pin sets a pinned version for the given service.
func (s *Store) Pin(service, version, pinnedBy, reason string) error {
	if service == "" {
		return fmt.Errorf("service name must not be empty")
	}
	if version == "" {
		return fmt.Errorf("version must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[service] = Entry{
		Service:  service,
		Version:  version,
		PinnedAt: time.Now(),
		PinnedBy: pinnedBy,
		Reason:   reason,
	}
	return nil
}

// Unpin removes the pinned version for the given service.
func (s *Store) Unpin(service string) error {
	if service == "" {
		return fmt.Errorf("service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[service]; !ok {
		return fmt.Errorf("service %q is not pinned", service)
	}
	delete(s.entries, service)
	return nil
}

// Get returns the pinned entry for a service, if any.
func (s *Store) Get(service string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[service]
	return e, ok
}

// IsPinned reports whether a service has a pinned version.
func (s *Store) IsPinned(service string) bool {
	_, ok := s.Get(service)
	return ok
}

// All returns a copy of all pinned entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}
