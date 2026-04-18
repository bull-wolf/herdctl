package health

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Status represents the health state of a service.
type Status string

const (
	StatusUnknown  Status = "unknown"
	StatusHealthy  Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

// Result holds the latest health check result for a service.
type Result struct {
	Service   string
	Status    Status
	CheckedAt time.Time
	Err       string
}

// Checker runs HTTP health checks against service endpoints.
type Checker struct {
	mu      sync.RWMutex
	results map[string]Result
	client  *http.Client
}

// New creates a new Checker with a default HTTP client.
func New() *Checker {
	return &Checker{
		results: make(map[string]Result),
		client:  &http.Client{Timeout: 3 * time.Second},
	}
}

// Check performs an HTTP GET against url and records the result for service.
func (c *Checker) Check(service, url string) Result {
	r := Result{
		Service:   service,
		CheckedAt: time.Now(),
	}

	resp, err := c.client.Get(url)
	if err != nil {
		r.Status = StatusUnhealthy
		r.Err = err.Error()
	} else {
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			r.Status = StatusHealthy
		} else {
			r.Status = StatusUnhealthy
			r.Err = fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
		}
	}

	c.mu.Lock()
	c.results[service] = r
	c.mu.Unlock()
	return r
}

// Get returns the last recorded result for a service.
func (c *Checker) Get(service string) (Result, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	r, ok := c.results[service]
	return r, ok
}

// All returns a copy of all recorded results.
func (c *Checker) All() []Result {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Result, 0, len(c.results))
	for _, r := range c.results {
		out = append(out, r)
	}
	return out
}
