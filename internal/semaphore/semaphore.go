package semaphore

import (
	"fmt"
	"sync"
)

// Semaphore manages named counting semaphores to limit concurrent access per service.
type Semaphore struct {
	mu      sync.Mutex
	entries map[string]*entry
}

type entry struct {
	max     int
	active  int
	waiting int
}

// New returns an initialised Semaphore.
func New() *Semaphore {
	return &Semaphore{entries: make(map[string]*entry)}
}

// Set configures the maximum concurrency for a service.
func (s *Semaphore) Set(service string, max int) error {
	if service == "" {
		return fmt.Errorf("semaphore: service name must not be empty")
	}
	if max <= 0 {
		return fmt.Errorf("semaphore: max must be greater than zero")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[service] = &entry{max: max}
	return nil
}

// Acquire attempts to acquire a slot for the given service.
// Returns an error if the limit is reached or the service is not configured.
func (s *Semaphore) Acquire(service string) error {
	if service == "" {
		return fmt.Errorf("semaphore: service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[service]
	if !ok {
		return fmt.Errorf("semaphore: no semaphore configured for service %q", service)
	}
	if e.active >= e.max {
		e.waiting++
		return fmt.Errorf("semaphore: limit reached for service %q (%d/%d)", service, e.active, e.max)
	}
	e.active++
	return nil
}

// Release frees a slot for the given service.
func (s *Semaphore) Release(service string) error {
	if service == "" {
		return fmt.Errorf("semaphore: service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[service]
	if !ok {
		return fmt.Errorf("semaphore: no semaphore configured for service %q", service)
	}
	if e.active > 0 {
		e.active--
	}
	if e.waiting > 0 {
		e.waiting--
	}
	return nil
}

// State returns the current active and max counts for a service.
func (s *Semaphore) State(service string) (active, max int, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, found := s.entries[service]
	if !found {
		return 0, 0, false
	}
	return e.active, e.max, true
}
