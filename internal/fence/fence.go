package fence

import (
	"fmt"
	"sync"
	"time"
)

// State represents the current state of a fence gate.
type State string

const (
	StateOpen   State = "open"
	StateClosed State = "closed"
)

// Entry holds the fence state for a single service.
type Entry struct {
	Service   string
	State     State
	Reason    string
	UpdatedAt time.Time
}

// Fence controls whether a service is allowed to proceed (open) or blocked (closed).
type Fence struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New returns a new Fence.
func New() *Fence {
	return &Fence{
		entries: make(map[string]*Entry),
	}
}

// Open marks the fence gate as open for the given service, optionally with a reason.
func (f *Fence) Open(service, reason string) error {
	if service == "" {
		return fmt.Errorf("fence: service name required")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.entries[service] = &Entry{
		Service:   service,
		State:     StateOpen,
		Reason:    reason,
		UpdatedAt: time.Now(),
	}
	return nil
}

// Close marks the fence gate as closed for the given service, optionally with a reason.
func (f *Fence) Close(service, reason string) error {
	if service == "" {
		return fmt.Errorf("fence: service name required")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.entries[service] = &Entry{
		Service:   service,
		State:     StateClosed,
		Reason:    reason,
		UpdatedAt: time.Now(),
	}
	return nil
}

// IsOpen returns true if the fence gate is open for the given service.
// Services with no entry are considered open by default.
func (f *Fence) IsOpen(service string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	e, ok := f.entries[service]
	if !ok {
		return true
	}
	return e.State == StateOpen
}

// Get returns the Entry for the given service and whether it exists.
func (f *Fence) Get(service string) (Entry, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	e, ok := f.entries[service]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a copy of all fence entries.
func (f *Fence) All() []Entry {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make([]Entry, 0, len(f.entries))
	for _, e := range f.entries {
		out = append(out, *e)
	}
	return out
}
