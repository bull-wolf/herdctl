package eventbus

import (
	"testing"

	"github.com/herdctl/herdctl/internal/config"
)

func makeConfig(names ...string) *config.Config {
	cfg := &config.Config{}
	for _, n := range names {
		cfg.Services = append(cfg.Services, config.Service{Name: n, Command: "echo " + n})
	}
	return cfg
}

func TestLoadFromConfig_RegistersSubscriptions(t *testing.T) {
	b := New()
	cfg := makeConfig("api", "db")
	if err := LoadFromConfig(b, cfg, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// each service gets 4 events: started, stopped, failed, restarted
	list := b.List()
	if len(list) != 8 {
		t.Fatalf("expected 8 subscriptions, got %d", len(list))
	}
}

func TestLoadFromConfig_EmptyServices(t *testing.T) {
	b := New()
	cfg := makeConfig()
	if err := LoadFromConfig(b, cfg, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.List()) != 0 {
		t.Fatal("expected no subscriptions for empty config")
	}
}

func TestLoadFromConfig_HandlerFires(t *testing.T) {
	b := New()
	cfg := makeConfig("worker")
	var msgs []string
	if err := LoadFromConfig(b, cfg, func(s string) { msgs = append(msgs, s) }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = b.Publish("worker", "started", nil)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
}

func TestLoadFromConfig_MultipleServices(t *testing.T) {
	b := New()
	cfg := makeConfig("a", "b", "c")
	if err := LoadFromConfig(b, cfg, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b.List()) != 12 {
		t.Fatalf("expected 12 subscriptions, got %d", len(b.List()))
	}
}
