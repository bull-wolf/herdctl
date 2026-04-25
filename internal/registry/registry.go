package registry

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds metadata about a registered service instance.
type Entry struct {
	Service     string
	Host        string
	Port        int
	RegisteredAt time.Time
	Meta        map[string]string
}

// Registry tracks service registration metadata for discovery.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New returns an initialized Registry.
func New() *Registry {
	return &Registry{
		entries: make(map[string]*Entry),
	}
}

// Register adds or replaces a service entry in the registry.
func (r *Registry) Register(service, host string, port int, meta map[string]string) error {
	if service == "" {
		return fmt.Errorf("service name must not be empty")
	}
	if host == "" {
		return fmt.Errorf("host must not be empty")
	}
	if port <= 0 || port > 65535 {
		return fmt.Errorf("port %d is out of valid range", port)
	}

	copied := make(map[string]string, len(meta))
	for k, v := range meta {
		copied[k] = v
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[service] = &Entry{
		Service:      service,
		Host:         host,
		Port:         port,
		RegisteredAt: time.Now(),
		Meta:         copied,
	}
	return nil
}

// Deregister removes a service entry from the registry.
func (r *Registry) Deregister(service string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[service]; !ok {
		return fmt.Errorf("service %q not registered", service)
	}
	delete(r.entries, service)
	return nil
}

// Get returns the entry for a named service.
func (r *Registry) Get(service string) (*Entry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[service]
	if !ok {
		return nil, false
	}
	copy := *e
	return &copy, true
}

// All returns a snapshot of all registered entries.
func (r *Registry) All() []*Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Entry, 0, len(r.entries))
	for _, e := range r.entries {
		copy := *e
		out = append(out, &copy)
	}
	return out
}
