package drain

import (
	"sync"
	"time"
)

// Drainer tracks in-flight operations for a service and allows graceful
// shutdown by waiting for all active work to complete before a timeout.
type Drainer struct {
	mu       sync.Mutex
	counters map[string]int
	waiters  map[string][]chan struct{}
}

// New returns a new Drainer.
func New() *Drainer {
	return &Drainer{
		counters: make(map[string]int),
		waiters:  make(map[string][]chan struct{}),
	}
}

// Acquire increments the in-flight counter for the given service.
func (d *Drainer) Acquire(service string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.counters[service]++
}

// Release decrements the in-flight counter for the given service.
// If the counter reaches zero, any goroutines waiting on Drain are unblocked.
func (d *Drainer) Release(service string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.counters[service] > 0 {
		d.counters[service]--
	}
	if d.counters[service] == 0 {
		for _, ch := range d.waiters[service] {
			close(ch)
		}
		delete(d.waiters, service)
	}
}

// Active returns the current number of in-flight operations for a service.
func (d *Drainer) Active(service string) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.counters[service]
}

// Drain blocks until all in-flight operations for the given service complete
// or the timeout elapses. Returns true if drained cleanly, false on timeout.
func (d *Drainer) Drain(service string, timeout time.Duration) bool {
	d.mu.Lock()
	if d.counters[service] == 0 {
		d.mu.Unlock()
		return true
	}
	ch := make(chan struct{})
	d.waiters[service] = append(d.waiters[service], ch)
	d.mu.Unlock()

	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Reset clears all counters and waiters for a service.
func (d *Drainer) Reset(service string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.counters, service)
	for _, ch := range d.waiters[service] {
		close(ch)
	}
	delete(d.waiters, service)
}
