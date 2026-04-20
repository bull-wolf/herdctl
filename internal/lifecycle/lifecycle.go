// Package lifecycle coordinates the startup and shutdown of services,
// respecting dependency order, health checks, and environment resolution.
package lifecycle

import (
	"fmt"
	"log"
	"time"

	"github.com/user/herdctl/internal/config"
	"github.com/user/herdctl/internal/deps"
	"github.com/user/herdctl/internal/env"
	"github.com/user/herdctl/internal/health"
	"github.com/user/herdctl/internal/runner"
	"github.com/user/herdctl/internal/status"
)

// Manager orchestrates service start/stop across the full dependency graph.
type Manager struct {
	cfg    *config.Config
	deps   *deps.Graph
	runner *runner.Runner
	health *health.Checker
	status *status.Store
	env    *env.Resolver
}

// New creates a new lifecycle Manager wired up with the provided subsystems.
func New(
	cfg *config.Config,
	d *deps.Graph,
	r *runner.Runner,
	h *health.Checker,
	s *status.Store,
	e *env.Resolver,
) *Manager {
	return &Manager{
		cfg:    cfg,
		deps:   d,
		runner: r,
		health: h,
		status: s,
		env:    e,
	}
}

// StartAll starts all services in dependency order.
// Services with unsatisfied dependencies are skipped with an error logged.
func (m *Manager) StartAll() error {
	order, err := m.deps.Order()
	if err != nil {
		return fmt.Errorf("resolving dependency order: %w", err)
	}

	for _, name := range order {
		if err := m.startOne(name); err != nil {
			log.Printf("[lifecycle] failed to start %q: %v", name, err)
			m.status.Set(name, "error")
			return fmt.Errorf("starting service %q: %w", name, err)
		}
	}
	return nil
}

// StopAll stops all services in reverse dependency order.
func (m *Manager) StopAll() error {
	order, err := m.deps.Order()
	if err != nil {
		return fmt.Errorf("resolving dependency order: %w", err)
	}

	// Reverse the order for shutdown.
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	var firstErr error
	for _, name := range order {
		if err := m.runner.Stop(name); err != nil {
			log.Printf("[lifecycle] failed to stop %q: %v", name, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		m.status.Set(name, "stopped")
	}
	return firstErr
}

// StartService starts a single named service after ensuring its dependencies are running.
func (m *Manager) StartService(name string) error {
	direct := m.deps.Deps(name)
	for _, dep := range direct {
		st, _ := m.status.Get(dep)
		if st != "running" {
			return fmt.Errorf("dependency %q of %q is not running (status: %q)", dep, name, st)
		}
	}
	return m.startOne(name)
}

// startOne starts a single service and waits for its health check to pass if configured.
func (m *Manager) startOne(name string) error {
	if err := m.runner.Start(name); err != nil {
		return err
	}
	m.status.Set(name, "starting")

	svc, ok := m.cfg.Services[name]
	if !ok {
		return fmt.Errorf("service %q not found in config", name)
	}

	// If a health check URL is configured, poll until healthy or timeout.
	if svc.HealthCheck != "" {
		if err := m.waitHealthy(name, svc.HealthCheck, 15*time.Second); err != nil {
			return fmt.Errorf("health check failed for %q: %w", name, err)
		}
	}

	m.status.Set(name, "running")
	log.Printf("[lifecycle] service %q is running", name)
	return nil
}

// waitHealthy polls the health checker until the service is healthy or the deadline is exceeded.
func (m *Manager) waitHealthy(name, url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		result := m.health.Check(name, url)
		if result.Healthy {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timed out after %s waiting for %q to become healthy", timeout, name)
}
