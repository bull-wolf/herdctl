package runner

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/user/herdctl/internal/config"
)

// ServiceStatus tracks the state of a running service.
type ServiceStatus struct {
	Name    string
	Cmd     *exec.Cmd
	Running bool
	Err     error
}

// Runner manages starting and stopping services.
type Runner struct {
	cfg      *config.Config
	mu       sync.Mutex
	services map[string]*ServiceStatus
}

// New creates a Runner for the given config.
func New(cfg *config.Config) *Runner {
	return &Runner{
		cfg:      cfg,
		services: make(map[string]*ServiceStatus),
	}
}

// Start launches a single service by name.
func (r *Runner) Start(name string) error {
	svc, ok := r.cfg.Services[name]
	if !ok {
		return fmt.Errorf("service %q not found", name)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if s, exists := r.services[name]; exists && s.Running {
		return fmt.Errorf("service %q is already running", name)
	}

	cmd := exec.Command("sh", "-c", svc.Command)
	cmd.Dir = svc.Dir

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %q: %w", name, err)
	}

	r.services[name] = &ServiceStatus{Name: name, Cmd: cmd, Running: true}
	return nil
}

// Stop terminates a running service by name.
func (r *Runner) Stop(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.services[name]
	if !ok || !s.Running {
		return fmt.Errorf("service %q is not running", name)
	}

	if err := s.Cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop %q: %w", name, err)
	}

	s.Running = false
	return nil
}

// Status returns the current status snapshot for all services.
func (r *Runner) Status() map[string]*ServiceStatus {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make(map[string]*ServiceStatus, len(r.services))
	for k, v := range r.services {
		copy := *v
		out[k] = &copy
	}
	return out
}
