package backoff

import (
	"errors"
	"sync"
	"time"
)

// Strategy defines the backoff strategy type.
type Strategy string

const (
	StrategyFixed       Strategy = "fixed"
	StrategyExponential Strategy = "exponential"

	defaultBase    = 500 * time.Millisecond
	defaultMaxWait = 30 * time.Second
)

// Entry holds backoff state for a single service.
type Entry struct {
	Service  string
	Attempts int
	LastWait time.Duration
	Strategy Strategy
	Base     time.Duration
	MaxWait  time.Duration
}

// Manager tracks backoff state per service.
type Manager struct {
	mu      sync.Mutex
	entries map[string]*Entry
}

// New creates a new backoff Manager.
func New() *Manager {
	return &Manager{
		entries: make(map[string]*Entry),
	}
}

// Configure sets the backoff strategy and parameters for a service.
func (m *Manager) Configure(service string, strategy Strategy, base, maxWait time.Duration) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	if base <= 0 {
		return errors.New("base duration must be positive")
	}
	if maxWait <= 0 {
		return errors.New("max wait duration must be positive")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[service] = &Entry{
		Service:  service,
		Strategy: strategy,
		Base:     base,
		MaxWait:  maxWait,
	}
	return nil
}

// Next increments the attempt counter and returns the next wait duration.
func (m *Manager) Next(service string) (time.Duration, error) {
	if service == "" {
		return 0, errors.New("service name must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[service]
	if !ok {
		e = &Entry{
			Service:  service,
			Strategy: StrategyExponential,
			Base:     defaultBase,
			MaxWait:  defaultMaxWait,
		}
		m.entries[service] = e
	}
	e.Attempts++
	var wait time.Duration
	switch e.Strategy {
	case StrategyFixed:
		wait = e.Base
	default:
		shift := e.Attempts - 1
		if shift > 30 {
			shift = 30
		}
		wait = e.Base * (1 << uint(shift))
	}
	if wait > e.MaxWait {
		wait = e.MaxWait
	}
	e.LastWait = wait
	return wait, nil
}

// Reset clears the attempt counter for a service.
func (m *Manager) Reset(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if e, ok := m.entries[service]; ok {
		e.Attempts = 0
		e.LastWait = 0
	}
}

// Get returns the current backoff entry for a service.
func (m *Manager) Get(service string) (Entry, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[service]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}
