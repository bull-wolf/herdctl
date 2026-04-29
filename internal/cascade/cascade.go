package cascade

import (
	"fmt"
	"sync"
)

// Policy defines how a cascade event propagates to dependents.
type Policy string

const (
	PolicyStop    Policy = "stop"    // stop dependents when source stops
	PolicyRestart Policy = "restart" // restart dependents when source restarts
	PolicyIgnore  Policy = "ignore"  // do not propagate
)

// Rule describes cascade behaviour for a service.
type Rule struct {
	Service  string
	Policy   Policy
	Targets  []string
}

// Manager stores and evaluates cascade rules.
type Manager struct {
	mu    sync.RWMutex
	rules map[string]Rule
}

// New returns an initialised Manager.
func New() *Manager {
	return &Manager{rules: make(map[string]Rule)}
}

// Register adds or replaces the cascade rule for a service.
func (m *Manager) Register(service string, policy Policy, targets []string) error {
	if service == "" {
		return fmt.Errorf("cascade: service name must not be empty")
	}
	if policy != PolicyStop && policy != PolicyRestart && policy != PolicyIgnore {
		return fmt.Errorf("cascade: unknown policy %q", policy)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rules[service] = Rule{Service: service, Policy: policy, Targets: append([]string(nil), targets...)}
	return nil
}

// Get returns the rule registered for service.
func (m *Manager) Get(service string) (Rule, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rules[service]
	return r, ok
}

// Evaluate returns the targets affected by a cascade event from service.
// If no rule exists or the policy is ignore, an empty slice is returned.
func (m *Manager) Evaluate(service string) (Policy, []string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rules[service]
	if !ok || r.Policy == PolicyIgnore {
		return PolicyIgnore, nil
	}
	return r.Policy, append([]string(nil), r.Targets...)
}

// All returns every registered rule.
func (m *Manager) All() []Rule {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Rule, 0, len(m.rules))
	for _, r := range m.rules {
		out = append(out, r)
	}
	return out
}

// Deregister removes the rule for service.
func (m *Manager) Deregister(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.rules, service)
}
