package runner

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/user/herdctl/internal/config"
	"github.com/user/herdctl/internal/logs"
)

// Runner manages the lifecycle of configured services.
type Runner struct {
	mu        sync.Mutex
	cfg       *config.Config
	processes map[string]*exec.Cmd
	Logs      *logs.Collector
}

// New creates a Runner for the given config.
func New(cfg *config.Config) *Runner {
	return &Runner{
		cfg:       cfg,
		processes: make(map[string]*exec.Cmd),
		Logs:      logs.New(),
	}
}

// Start launches a service by name, respecting dependency order.
func (r *Runner) Start(name string) error {
	svc, ok := r.cfg.Services[name]
	if !ok {
		return fmt.Errorf("unknown service: %s", name)
	}

	for _, dep := range svc.Deps {
		r.mu.Lock()
		_, running := r.processes[dep]
		r.mu.Unlock()
		if !running {
			if err := r.Start(dep); err != nil {
				return fmt.Errorf("dependency %s failed: %w", dep, err)
			}
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, running := r.processes[name]; running {
		return fmt.Errorf("service %s is already running", name)
	}

	cmd := exec.Command("sh", "-c", svc.Command)

	pr, pw, err := createPipe()
	if err == nil {
		cmd.Stdout = pw
		cmd.Stderr = pw
		r.Logs.Register(name, pr)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", name, err)
	}

	r.processes[name] = cmd
	r.Logs.Write(name, fmt.Sprintf("started (pid %d)", cmd.Process.Pid))
	return nil
}

// Stop terminates a running service.
func (r *Runner) Stop(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cmd, ok := r.processes[name]
	if !ok {
		return fmt.Errorf("service %s is not running", name)
	}

	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop %s: %w", name, err)
	}

	delete(r.processes, name)
	r.Logs.Write(name, "stopped")
	return nil
}

// Running returns the names of all currently running services.
func (r *Runner) Running() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	names := make([]string, 0, len(r.processes))
	for n := range r.processes {
		names = append(names, n)
	}
	return names
}

// createPipe is a thin wrapper so tests can override it.
var createPipe = func() (*readWriter, *readWriter, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

type readWriter struct{}

func (rw *readWriter) Write(p []byte) (int, error) { return len(p), nil }
func (rw *readWriter) Read(p []byte) (int, error)  { return 0, nil }
