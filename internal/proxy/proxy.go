package proxy

import (
	"errors"
	"fmt"
	"sync"
)

// Entry represents a proxy rule mapping a service to a local port via an upstream target.
type Entry struct {
	Service  string
	Port     int
	Upstream string
	Enabled  bool
}

// Proxy manages port-forwarding proxy rules for services.
type Proxy struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns a new Proxy instance.
func New() *Proxy {
	return &Proxy{
		entries: make(map[string]Entry),
	}
}

// Register adds or replaces a proxy rule for the given service.
func (p *Proxy) Register(service string, port int, upstream string) error {
	if service == "" {
		return errors.New("proxy: service name must not be empty")
	}
	if port <= 0 || port > 65535 {
		return fmt.Errorf("proxy: invalid port %d for service %q", port, service)
	}
	if upstream == "" {
		return fmt.Errorf("proxy: upstream must not be empty for service %q", service)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.entries[service] = Entry{
		Service:  service,
		Port:     port,
		Upstream: upstream,
		Enabled:  true,
	}
	return nil
}

// Get returns the proxy entry for the given service.
func (p *Proxy) Get(service string) (Entry, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	e, ok := p.entries[service]
	return e, ok
}

// Disable marks a proxy rule as disabled without removing it.
func (p *Proxy) Disable(service string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	e, ok := p.entries[service]
	if !ok {
		return fmt.Errorf("proxy: no entry for service %q", service)
	}
	e.Enabled = false
	p.entries[service] = e
	return nil
}

// Enable marks a previously disabled proxy rule as active.
func (p *Proxy) Enable(service string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	e, ok := p.entries[service]
	if !ok {
		return fmt.Errorf("proxy: no entry for service %q", service)
	}
	e.Enabled = true
	p.entries[service] = e
	return nil
}

// All returns a copy of all registered proxy entries.
func (p *Proxy) All() []Entry {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Entry, 0, len(p.entries))
	for _, e := range p.entries {
		out = append(out, e)
	}
	return out
}

// Deregister removes the proxy rule for the given service.
func (p *Proxy) Deregister(service string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.entries[service]; !ok {
		return fmt.Errorf("proxy: no entry for service %q", service)
	}
	delete(p.entries, service)
	return nil
}
