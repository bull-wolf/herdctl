package limiter

import (
	"testing"

	"github.com/user/herdctl/internal/config"
)

func makeConfig(services map[string]config.Service) *config.Config {
	return &config.Config{Services: services}
}

func TestLoadFromConfig_SetsLimits(t *testing.T) {
	cfg := makeConfig(map[string]config.Service{
		"api":    {Command: "go run .", Concurrency: 4},
		"worker": {Command: "./worker", Concurrency: 2},
	})
	l := New()
	if err := LoadFromConfig(l, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := l.Get("api")
	if !ok || e.Max != 4 {
		t.Errorf("expected api max=4, got %+v", e)
	}
	e, ok = l.Get("worker")
	if !ok || e.Max != 2 {
		t.Errorf("expected worker max=2, got %+v", e)
	}
}

func TestLoadFromConfig_SkipsZeroValues(t *testing.T) {
	cfg := makeConfig(map[string]config.Service{
		"api":    {Command: "go run .", Concurrency: 0},
		"worker": {Command: "./worker", Concurrency: 3},
	})
	l := New()
	if err := LoadFromConfig(l, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := l.Get("api"); ok {
		t.Error("expected api to be skipped")
	}
	if _, ok := l.Get("worker"); !ok {
		t.Error("expected worker to be registered")
	}
}

func TestLoadFromConfig_EmptyServices(t *testing.T) {
	cfg := makeConfig(map[string]config.Service{})
	l := New()
	if err := LoadFromConfig(l, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.All()) != 0 {
		t.Error("expected no entries")
	}
}

func TestLoadFromConfig_MultipleServices(t *testing.T) {
	cfg := makeConfig(map[string]config.Service{
		"a": {Command: "a", Concurrency: 1},
		"b": {Command: "b", Concurrency: 5},
		"c": {Command: "c", Concurrency: 10},
	})
	l := New()
	if err := LoadFromConfig(l, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.All()) != 3 {
		t.Errorf("expected 3 entries, got %d", len(l.All()))
	}
}
