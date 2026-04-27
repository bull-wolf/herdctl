package quota

import (
	"testing"

	"github.com/user/herdctl/internal/config"
)

func makeConfig(services []config.Service) *config.Config {
	return &config.Config{Services: services}
}

func TestLoadFromConfig_SetsQuotas(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{
			Name:    "api",
			Command: "go run main.go",
			Quota:   config.Quota{CPU: 1.5, Memory: 512, Procs: 4},
		},
	})
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e, ok := m.Get("api", KindCPU); !ok || e.Limit != 1.5 {
		t.Errorf("expected cpu limit 1.5, got %+v", e)
	}
	if e, ok := m.Get("api", KindMemory); !ok || e.Limit != 512 {
		t.Errorf("expected memory limit 512, got %+v", e)
	}
	if e, ok := m.Get("api", KindProcs); !ok || e.Limit != 4 {
		t.Errorf("expected procs limit 4, got %+v", e)
	}
}

func TestLoadFromConfig_SkipsZeroValues(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "worker", Command: "./worker", Quota: config.Quota{}},
	})
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if all := m.All(); len(all) != 0 {
		t.Errorf("expected no entries for zero quotas, got %d", len(all))
	}
}

func TestLoadFromConfig_MultipleServices(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "api", Command: "go run .", Quota: config.Quota{CPU: 2.0}},
		{Name: "db", Command: "postgres", Quota: config.Quota{Memory: 1024}},
	})
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m.Get("api", KindCPU); !ok {
		t.Error("expected cpu quota for api")
	}
	if _, ok := m.Get("db", KindMemory); !ok {
		t.Error("expected memory quota for db")
	}
}

func TestLoadFromConfig_EmptyServices(t *testing.T) {
	cfg := makeConfig(nil)
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error on empty config: %v", err)
	}
}
