package retention

import (
	"testing"
	"time"

	"github.com/user/herdctl/internal/config"
)

func makeConfig(services []config.Service) *config.Config {
	return &config.Config{Services: services}
}

func TestLoadFromConfig_SetsRetention(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "api", Command: "go run .", Retention: &config.RetentionConfig{MaxAgeSecs: 3600, MaxItems: 500}},
	})
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := m.Get("api")
	if !ok {
		t.Fatal("expected policy for api")
	}
	if p.MaxAge != time.Hour {
		t.Errorf("expected 1h, got %v", p.MaxAge)
	}
	if p.MaxItems != 500 {
		t.Errorf("expected 500, got %d", p.MaxItems)
	}
}

func TestLoadFromConfig_SkipsNilRetention(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "worker", Command: "./worker"},
	})
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := m.Get("worker"); ok {
		t.Fatal("expected no policy for worker")
	}
}

func TestLoadFromConfig_MultipleServices(t *testing.T) {
	cfg := makeConfig([]config.Service{
		{Name: "a", Command: "a", Retention: &config.RetentionConfig{MaxAgeSecs: 60, MaxItems: 100}},
		{Name: "b", Command: "b", Retention: &config.RetentionConfig{MaxAgeSecs: 120, MaxItems: 200}},
	})
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.All()) != 2 {
		t.Fatalf("expected 2 policies, got %d", len(m.All()))
	}
}

func TestLoadFromConfig_EmptyServices(t *testing.T) {
	cfg := makeConfig(nil)
	m := New()
	if err := LoadFromConfig(m, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.All()) != 0 {
		t.Fatal("expected no policies")
	}
}
