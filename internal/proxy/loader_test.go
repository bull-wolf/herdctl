package proxy

import (
	"testing"

	"github.com/example/herdctl/internal/config"
)

func makeConfig(services []config.Service) *config.Config {
	return &config.Config{Services: services}
}

func TestLoadFromConfig_RegistersProxies(t *testing.T) {
	p := New()
	cfg := makeConfig([]config.Service{
		{Name: "api", Command: "go run .", Proxy: &config.ProxyConfig{Port: 8080, Upstream: "localhost:9090"}},
		{Name: "web", Command: "npm start", Proxy: &config.ProxyConfig{Port: 3000, Upstream: "localhost:3001"}},
	})
	if err := LoadFromConfig(p, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(p.All()))
	}
	e, ok := p.Get("api")
	if !ok || e.Port != 8080 || e.Upstream != "localhost:9090" {
		t.Errorf("unexpected api entry: %+v", e)
	}
}

func TestLoadFromConfig_SkipsServicesWithoutProxy(t *testing.T) {
	p := New()
	cfg := makeConfig([]config.Service{
		{Name: "worker", Command: "./worker"},
	})
	if err := LoadFromConfig(p, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.All()) != 0 {
		t.Errorf("expected 0 entries, got %d", len(p.All()))
	}
}

func TestLoadFromConfig_EmptyServices(t *testing.T) {
	p := New()
	cfg := makeConfig([]config.Service{})
	if err := LoadFromConfig(p, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.All()) != 0 {
		t.Error("expected empty proxy table")
	}
}

func TestLoadFromConfig_InvalidProxy(t *testing.T) {
	p := New()
	cfg := makeConfig([]config.Service{
		{Name: "bad", Command: "./bad", Proxy: &config.ProxyConfig{Port: -1, Upstream: "localhost:9999"}},
	})
	if err := LoadFromConfig(p, cfg); err == nil {
		t.Error("expected error for invalid port")
	}
}
