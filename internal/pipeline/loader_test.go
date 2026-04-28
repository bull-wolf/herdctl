package pipeline

import (
	"testing"

	"github.com/example/herdctl/internal/config"
)

func makeConfig(services []config.Service) *config.Config {
	return &config.Config{Services: services}
}

func TestLoadFromConfig_RegistersStages(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "api", Command: "go run .", Pipeline: []string{"build", "test", "deploy"}},
		{Name: "worker", Command: "./worker", Pipeline: []string{"build", "deploy"}},
	})
	p := New()
	if err := LoadFromConfig(p, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stages, err := p.Get("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stages) != 3 {
		t.Errorf("expected 3 stages for api, got %d", len(stages))
	}
	stages, err = p.Get("worker")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stages) != 2 {
		t.Errorf("expected 2 stages for worker, got %d", len(stages))
	}
}

func TestLoadFromConfig_SkipsServicesWithoutPipeline(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "api", Command: "go run ."},
	})
	p := New()
	if err := LoadFromConfig(p, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := p.Get("api")
	if err == nil {
		t.Error("expected error for service with no pipeline registered")
	}
}

func TestLoadFromConfig_EmptyServices(t *testing.T) {
	cfg := makeConfig(nil)
	p := New()
	if err := LoadFromConfig(p, cfg); err != nil {
		t.Fatalf("unexpected error on empty services: %v", err)
	}
}
