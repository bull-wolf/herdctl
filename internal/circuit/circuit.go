package circuit

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State string

const (
	StateClosed   State = "closed"
	StateOpen     State = "open"
	StateHalfOpen State = "half-open"
)

// Entry holds circuit breaker state for a single service.
type Entry struct {
	Service   string
	State     State
	Failures  int
	LastError error
	OpenedAt  time.Time
}

// Breaker manages circuit breaker state per service.
type Breaker struct {
	mu        sync.RWMutex
	entries   map[string]*Entry
	threshold int
	cooldown  time.Duration
}

// New creates a Breaker with the given failure threshold and cooldown duration.
func New(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{
		entries:   make(map[string]*Entry),
		threshold: threshold,
		cooldown:  cooldown,
	}
}

// RecordFailure increments the failure count for a service and may open the circuit.
func (b *Breaker) RecordFailure(service string, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	e := b.getOrCreate(service)
	e.Failures++
	e.LastError = err
	if e.State == StateClosed && e.Failures >= b.threshold {
		e.State = StateOpen
		e.OpenedAt = time.Now()
	}
}

// RecordSuccess resets failures and closes the circuit for a service.
func (b *Breaker) RecordSuccess(service string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	e := b.getOrCreate(service)
	e.Failures = 0
	e.LastError = nil
	e.State = StateClosed
}

// Allow returns nil if the service is allowed to proceed, or an error if the circuit is open.
func (b *Breaker) Allow(service string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	e := b.getOrCreate(service)
	switch e.State {
	case StateOpen:
		if time.Since(e.OpenedAt) >= b.cooldown {
			e.State = StateHalfOpen
			return nil
		}
		return fmt.Errorf("circuit open for service %q: %w", service, errors.New("too many failures"))
	case StateHalfOpen:
		return nil
	}
	return nil
}

// Get returns the current Entry for a service, or false if not found.
func (b *Breaker) Get(service string) (Entry, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	e, ok := b.entries[service]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Reset clears the circuit state for a service.
func (b *Breaker) Reset(service string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.entries, service)
}

func (b *Breaker) getOrCreate(service string) *Entry {
	if e, ok := b.entries[service]; ok {
		return e
	}
	e := &Entry{Service: service, State: StateClosed}
	b.entries[service] = e
	return e
}
