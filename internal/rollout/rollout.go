package rollout

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Strategy defines how a rollout is performed.
type Strategy string

const (
	StrategyAll      Strategy = "all"
	StrategySequential Strategy = "sequential"
	StrategyCanary   Strategy = "canary"
)

// State represents the current rollout state for a service.
type State struct {
	Service   string
	Strategy  Strategy
	StartedAt time.Time
	Done      bool
	Err       error
}

// Manager tracks rollout state across services.
type Manager struct {
	mu     sync.RWMutex
	states map[string]*State
}

// New creates a new rollout Manager.
func New() *Manager {
	return &Manager{
		states: make(map[string]*State),
	}
}

// Begin records the start of a rollout for a service.
func (m *Manager) Begin(service string, strategy Strategy) error {
	if service == "" {
		return errors.New("rollout: service name must not be empty")
	}
	if strategy != StrategyAll && strategy != StrategySequential && strategy != StrategyCanary {
		return fmt.Errorf("rollout: unknown strategy %q", strategy)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[service] = &State{
		Service:   service,
		Strategy:  strategy,
		StartedAt: time.Now(),
		Done:      false,
	}
	return nil
}

// Complete marks a rollout as finished, optionally with an error.
func (m *Manager) Complete(service string, err error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.states[service]
	if !ok {
		return fmt.Errorf("rollout: no active rollout for service %q", service)
	}
	s.Done = true
	s.Err = err
	return nil
}

// Get returns the rollout state for a service.
func (m *Manager) Get(service string) (*State, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.states[service]
	if !ok {
		return nil, false
	}
	copy := *s
	return &copy, true
}

// All returns rollout states for all services.
func (m *Manager) All() []*State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*State, 0, len(m.states))
	for _, s := range m.states {
		copy := *s
		out = append(out, &copy)
	}
	return out
}
