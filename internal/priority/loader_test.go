package priority

import (
	"testing"

	"github.com/seanmorris/herdctl/internal/config"
)

func makeConfig(services map[string]config.Service) *config.Config {
	return &config.Config{Services: services}
}

func TestLoadFromConfig_SetsPriorities(t *testing.T) {
	cfg := makeConfig(map[string]config.Service{
		"api": {Command: "go run .", Priority: 100},
		"db": {Command: "postgres", Priority: 10},
	})
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l := m.GetLevel("api"); l != 100 {
		t.Errorf("expected 100, got %d", l)
	}
	if l := m.GetLevel("db"); l != 10 {
		t.Errorf("expected 10, got %d", l)
	}
}

func TestLoadFromConfig_SkipsZeroValues(t *testing.T) {
	cfg := makeConfig(map[string]config.Service{
		"worker": {Command: "./worker", Priority: 0},
	})
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := m.Get("worker")
	if ok {
		t.Error("expected worker to be absent (zero priority skipped)")
	}
}

func TestLoadFromConfig_EmptyServices(t *testing.T) {
	cfg := makeConfig(map[string]config.Service{})
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.All()) != 0 {
		t.Error("expected empty manager")
	}
}

func TestLoadFromConfig_MultipleServices(t *testing.T) {
	cfg := makeConfig(map[string]config.Service{
		"a": {Command: "./a", Priority: 80},
		"b": {Command: "./b", Priority: 20},
		"c": {Command: "./c", Priority: 50},
	})
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.All()) != 3 {
		t.Errorf("expected 3 entries, got %d", len(m.All()))
	}
	sorted := m.Sorted([]string{"a", "b", "c"})
	if sorted[0] != "a" {
		t.Errorf("expected a first, got %s", sorted[0])
	}
}
