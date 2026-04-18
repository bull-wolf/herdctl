package logs

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Entry represents a single log line from a service.
type Entry struct {
	Service   string
	Timestamp time.Time
	Line      string
}

// Collector aggregates log streams from multiple services.
type Collector struct {
	mu      sync.Mutex
	entries []Entry
	writers map[string]io.Writer
}

// New creates a new Collector.
func New() *Collector {
	return &Collector{
		writers: make(map[string]io.Writer),
	}
}

// Register adds a writer for a named service (e.g. os.Pipe or bytes.Buffer).
func (c *Collector) Register(service string, w io.Writer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.writers[service] = w
}

// Write records a log line for the given service and forwards it to the
// registered writer (falling back to stdout).
func (c *Collector) Write(service, line string) {
	entry := Entry{
		Service:   service,
		Timestamp: time.Now(),
		Line:      line,
	}

	c.mu.Lock()
	c.entries = append(c.entries, entry)
	w, ok := c.writers[service]
	c.mu.Unlock()

	formatted := fmt.Sprintf("[%s] %s | %s\n", entry.Timestamp.Format("15:04:05"), service, line)
	if ok {
		fmt.Fprint(w, formatted)
	} else {
		fmt.Fprint(os.Stdout, formatted)
	}
}

// Tail returns the last n entries across all services, or all if n <= 0.
func (c *Collector) Tail(n int) []Entry {
	c.mu.Lock()
	defer c.mu.Unlock()
	if n <= 0 || n >= len(c.entries) {
		copy := make([]Entry, len(c.entries))
		copy = append([]Entry{}, c.entries...)
		return copy
	}
	return append([]Entry{}, c.entries[len(c.entries)-n:]...)
}

// TailService returns the last n entries for a specific service.
func (c *Collector) TailService(service string, n int) []Entry {
	c.mu.Lock()
	defer c.mu.Unlock()
	var filtered []Entry
	for _, e := range c.entries {
		if e.Service == service {
			filtered = append(filtered, e)
		}
	}
	if n <= 0 || n >= len(filtered) {
		return filtered
	}
	return filtered[len(filtered)-n:]
}
