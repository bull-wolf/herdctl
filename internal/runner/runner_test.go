package runner_test

import (
	"testing"

	"github.com/user/herdctl/internal/config"
	"github.com/user/herdctl/internal/runner"
)

func testConfig() *config.Config {
	return &config.Config{
		Services: map[string]config.Service{
			"echo": {
				Command: "echo hello",
			},
		},
	}
}

func TestStart_UnknownService(t *testing.T) {
	r := runner.New(testConfig())
	if err := r.Start("nope"); err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestStart_And_Stop(t *testing.T) {
	r := runner.New(testConfig())

	if err := r.Start("echo"); err != nil {
		t.Fatalf("unexpected start error: %v", err)
	}

	statuses := r.Status()
	if !statuses["echo"].Running {
		t.Fatal("expected echo to be running")
	}

	if err := r.Stop("echo"); err != nil {
		t.Fatalf("unexpected stop error: %v", err)
	}

	statuses = r.Status()
	if statuses["echo"].Running {
		t.Fatal("expected echo to be stopped")
	}
}

func TestStart_AlreadyRunning(t *testing.T) {
	r := runner.New(testConfig())

	if err := r.Start("echo"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer r.Stop("echo") //nolint

	if err := r.Start("echo"); err == nil {
		t.Fatal("expected error when starting already-running service")
	}
}

func TestStop_NotRunning(t *testing.T) {
	r := runner.New(testConfig())
	if err := r.Stop("echo"); err == nil {
		t.Fatal("expected error when stopping non-running service")
	}
}
