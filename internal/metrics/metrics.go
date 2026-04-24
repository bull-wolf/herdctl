package metrics

import (
	"sync"
	"time"
)

// Entry holds a single recorded metric sample for a service.
type Entry struct {
	Service   string
	Timestamp time.Time
	CPU       float64 // percentage 0-100
	MemoryMB  float64
	UptimeSec int64
}

// Collector stores rolling metric samples per service.
type Collector struct {
	mu      sync.RWMutex
	samples map[string][]Entry
	maxLen  int
}

// New creates a Collector that retains up to maxLen samples per service.
func New(maxLen int) *Collector {
	if maxLen <= 0 {
		maxLen = 60
	}
	return &Collector{
		samples: make(map[string][]Entry),
		maxLen:  maxLen,
	}
}

// Record appends a new metric entry for the given service.
func (c *Collector) Record(e Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	list := c.samples[e.Service]
	list = append(list, e)
	if len(list) > c.maxLen {
		list = list[len(list)-c.maxLen:]
	}
	c.samples[e.Service] = list
}

// Latest returns the most recent Entry for a service, and whether it exists.
func (c *Collector) Latest(service string) (Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	list, ok := c.samples[service]
	if !ok || len(list) == 0 {
		return Entry{}, false
	}
	return list[len(list)-1], true
}

// History returns a copy of all recorded samples for a service.
func (c *Collector) History(service string) []Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	list := c.samples[service]
	out := make([]Entry, len(list))
	copy(out, list)
	return out
}

// Services returns the names of all services that have recorded metrics.
func (c *Collector) Services() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.samples))
	for k := range c.samples {
		names = append(names, k)
	}
	return names
}
