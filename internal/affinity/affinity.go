package affinity

import (
	"fmt"
	"sync"
)

// Rule describes a co-location or anti-affinity constraint between services.
type Rule struct {
	Service string
	Target  string
	Kind    string // "together" | "apart"
}

// Store holds affinity rules for services.
type Store struct {
	mu    sync.RWMutex
	rules map[string][]Rule // keyed by service
}

// New returns an initialised Store.
func New() *Store {
	return &Store{rules: make(map[string][]Rule)}
}

// Add registers an affinity rule for a service.
func (s *Store) Add(service, target, kind string) error {
	if service == "" || target == "" {
		return fmt.Errorf("affinity: service and target must not be empty")
	}
	if kind != "together" && kind != "apart" {
		return fmt.Errorf("affinity: unknown kind %q, must be 'together' or 'apart'", kind)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.rules[service] {
		if r.Target == target && r.Kind == kind {
			return nil // idempotent
		}
	}
	s.rules[service] = append(s.rules[service], Rule{Service: service, Target: target, Kind: kind})
	return nil
}

// Get returns all affinity rules registered for a service.
func (s *Store) Get(service string) []Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := append([]Rule{}, s.rules[service]...)
	return copy
}

// Conflicts returns services that violate "apart" rules given the running set.
func (s *Store) Conflicts(service string, running []string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	runningSet := make(map[string]bool, len(running))
	for _, r := range running {
		runningSet[r] = true
	}
	var out []string
	for _, rule := range s.rules[service] {
		if rule.Kind == "apart" && runningSet[rule.Target] {
			out = append(out, rule.Target)
		}
	}
	return out
}

// Remove deletes all affinity rules for a service.
func (s *Store) Remove(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.rules, service)
}

// All returns every registered rule across all services.
func (s *Store) All() []Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Rule
	for _, rules := range s.rules {
		out = append(out, rules...)
	}
	return out
}
