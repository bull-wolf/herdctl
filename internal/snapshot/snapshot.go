package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry represents a point-in-time snapshot of a service's state.
type Entry struct {
	Service   string            `json:"service"`
	Status    string            `json:"status"`
	Health    string            `json:"health"`
	Env       map[string]string `json:"env,omitempty"`
	CapturedAt time.Time        `json:"captured_at"`
}

// Store holds snapshots keyed by service name.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New creates an empty snapshot Store.
func New() *Store {
	return &Store{
		entries: make(map[string]Entry),
	}
}

// Capture records a snapshot for the given service.
func (s *Store) Capture(e Entry) {
	e.CapturedAt = time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[e.Service] = e
}

// Get returns the latest snapshot for a service, or false if none exists.
func (s *Store) Get(service string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[service]
	return e, ok
}

// All returns a copy of all stored snapshots.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

// SaveToFile serialises all snapshots to a JSON file at the given path.
func (s *Store) SaveToFile(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	entries := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		entries = append(entries, e)
	}
	if err := enc.Encode(entries); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}
