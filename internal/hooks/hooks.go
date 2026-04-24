package hooks

import (
	"fmt"
	"sync"
)

// Event represents a lifecycle event that can trigger hooks.
type Event string

const (
	EventBeforeStart Event = "before_start"
	EventAfterStart  Event = "after_start"
	EventBeforeStop  Event = "before_stop"
	EventAfterStop   Event = "after_stop"
)

// HookFunc is a function invoked when a hook fires.
type HookFunc func(service string, event Event) error

// Registry stores and dispatches hooks for services.
type Registry struct {
	mu    sync.RWMutex
	hooks map[string]map[Event][]HookFunc
}

// New creates a new hook Registry.
func New() *Registry {
	return &Registry{
		hooks: make(map[string]map[Event][]HookFunc),
	}
}

// Register adds a HookFunc for the given service and event.
func (r *Registry) Register(service string, event Event, fn HookFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.hooks[service] == nil {
		r.hooks[service] = make(map[Event][]HookFunc)
	}
	r.hooks[service][event] = append(r.hooks[service][event], fn)
}

// Fire invokes all hooks registered for the given service and event.
// All hooks are called; errors are collected and returned as a combined error.
func (r *Registry) Fire(service string, event Event) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var errs []error
	for _, fn := range r.hooks[service][event] {
		if err := fn(service, event); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("hook errors for %s/%s: %v", service, event, errs)
}

// List returns the events that have hooks registered for a service.
func (r *Registry) List(service string) []Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]Event, 0, len(r.hooks[service]))
	for e := range r.hooks[service] {
		events = append(events, e)
	}
	return events
}

// Clear removes all hooks for a given service.
func (r *Registry) Clear(service string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.hooks, service)
}
