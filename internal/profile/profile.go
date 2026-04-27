package profile

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds a single profiling sample for a service operation.
type Entry struct {
	Service   string
	Operation string
	Duration  time.Duration
	RecordedAt time.Time
}

// Profiler records timing data for service operations.
type Profiler struct {
	mu      sync.RWMutex
	entries map[string][]Entry
	maxPerService int
}

// New creates a new Profiler with the given rolling window size per service.
func New(maxPerService int) *Profiler {
	if maxPerService <= 0 {
		maxPerService = 50
	}
	return &Profiler{
		entries:       make(map[string][]Entry),
		maxPerService: maxPerService,
	}
}

// Record stores a duration sample for the given service and operation.
func (p *Profiler) Record(service, operation string, d time.Duration) error {
	if service == "" {
		return fmt.Errorf("profile: service name must not be empty")
	}
	if operation == "" {
		return fmt.Errorf("profile: operation must not be empty")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	e := Entry{
		Service:    service,
		Operation:  operation,
		Duration:   d,
		RecordedAt: time.Now(),
	}
	p.entries[service] = append(p.entries[service], e)
	if len(p.entries[service]) > p.maxPerService {
		p.entries[service] = p.entries[service][len(p.entries[service])-p.maxPerService:]
	}
	return nil
}

// Latest returns the most recent entry for a service, or false if none exists.
func (p *Profiler) Latest(service string) (Entry, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	entries := p.entries[service]
	if len(entries) == 0 {
		return Entry{}, false
	}
	return entries[len(entries)-1], true
}

// History returns a copy of all recorded entries for a service.
func (p *Profiler) History(service string) []Entry {
	p.mu.RLock()
	defer p.mu.RUnlock()
	src := p.entries[service]
	out := make([]Entry, len(src))
	copy(out, src)
	return out
}

// Average returns the mean duration of all recorded entries for a service.
// Returns 0 and false if no entries exist.
func (p *Profiler) Average(service string) (time.Duration, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	entries := p.entries[service]
	if len(entries) == 0 {
		return 0, false
	}
	var total time.Duration
	for _, e := range entries {
		total += e.Duration
	}
	return total / time.Duration(len(entries)), true
}

// Clear removes all profiling data for a service.
func (p *Profiler) Clear(service string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.entries, service)
}
