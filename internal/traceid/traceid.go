package traceid

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Entry holds a trace ID and its metadata for a service.
type Entry struct {
	Service   string
	TraceID   string
	CreatedAt time.Time
}

// Store manages trace IDs per service.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
	rng     *rand.Rand
}

// New creates a new trace ID store.
func New() *Store {
	return &Store{
		entries: make(map[string]Entry),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate creates and stores a new trace ID for the given service.
func (s *Store) Generate(service string) (string, error) {
	if service == "" {
		return "", fmt.Errorf("traceid: service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	id := fmt.Sprintf("%016x", s.rng.Uint64())
	s.entries[service] = Entry{
		Service:   service,
		TraceID:   id,
		CreatedAt: time.Now(),
	}
	return id, nil
}

// Get returns the current trace entry for a service.
func (s *Store) Get(service string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[service]
	return e, ok
}

// Clear removes the trace ID for a service.
func (s *Store) Clear(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, service)
}

// All returns a copy of all trace entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}
