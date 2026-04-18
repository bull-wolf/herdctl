package env

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// Resolver holds per-service environment variable overrides.
type Resolver struct {
	mu      sync.RWMutex
	global  map[string]string
	service map[string]map[string]string
}

// New creates a new Resolver.
func New(global map[string]string) *Resolver {
	g := make(map[string]string, len(global))
	for k, v := range global {
		g[k] = v
	}
	return &Resolver{
		global:  g,
		service: make(map[string]map[string]string),
	}
}

// SetService sets environment overrides for a specific service.
func (r *Resolver) SetService(service string, vars map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	m := make(map[string]string, len(vars))
	for k, v := range vars {
		m[k] = v
	}
	r.service[service] = m
}

// Resolve returns the merged environment for a service as a slice of "KEY=VALUE"
// strings suitable for exec.Cmd.Env. Precedence: service > global > os env.
func (r *Resolver) Resolve(service string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	merged := make(map[string]string)

	// Start with OS environment.
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			merged[parts[0]] = parts[1]
		}
	}

	// Apply globals.
	for k, v := range r.global {
		merged[k] = v
	}

	// Apply service-specific overrides.
	if svc, ok := r.service[service]; ok {
		for k, v := range svc {
			merged[k] = v
		}
	}

	result := make([]string, 0, len(merged))
	for k, v := range merged {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return result
}

// Get returns the resolved value of a single key for a service.
func (r *Resolver) Get(service, key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if svc, ok := r.service[service]; ok {
		if v, ok := svc[key]; ok {
			return v, true
		}
	}
	if v, ok := r.global[key]; ok {
		return v, true
	}
	return os.LookupEnv(key)
}
