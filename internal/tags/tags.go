package tags

import (
	"fmt"
	"sort"
	"sync"
)

// Store manages arbitrary string tags associated with services.
// Tags allow grouping and filtering services by user-defined categories.
type Store struct {
	mu   sync.RWMutex
	data map[string]map[string]struct{} // service -> set of tags
}

// New returns an initialised tag Store.
func New() *Store {
	return &Store{
		data: make(map[string]map[string]struct{}),
	}
}

// Add attaches one or more tags to a service. Duplicate tags are silently ignored.
func (s *Store) Add(service string, tags ...string) error {
	if service == "" {
		return fmt.Errorf("tags: service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data[service] == nil {
		s.data[service] = make(map[string]struct{})
	}
	for _, t := range tags {
		if t == "" {
			return fmt.Errorf("tags: tag value must not be empty")
		}
		s.data[service][t] = struct{}{}
	}
	return nil
}

// Remove detaches a tag from a service. No-op if the tag is not present.
func (s *Store) Remove(service, tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if set, ok := s.data[service]; ok {
		delete(set, tag)
	}
}

// Get returns the sorted list of tags for a service.
func (s *Store) Get(service string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	set, ok := s.data[service]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// Has reports whether a service carries the given tag.
func (s *Store) Has(service, tag string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if set, ok := s.data[service]; ok {
		_, found := set[tag]
		return found
	}
	return false
}

// ServicesWithTag returns a sorted list of services that carry the given tag.
func (s *Store) ServicesWithTag(tag string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []string
	for svc, set := range s.data {
		if _, ok := set[tag]; ok {
			out = append(out, svc)
		}
	}
	sort.Strings(out)
	return out
}
