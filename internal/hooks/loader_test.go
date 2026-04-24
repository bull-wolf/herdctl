package hooks_test

import (
	"testing"

	"github.com/user/herdctl/internal/config"
	"github.com/user/herdctl/internal/hooks"
)

func makeConfig(hookMap map[string]string) *config.Config {
	return &config.Config{
		Services: []config.Service{
			{
				Name:    "api",
				Command: "go run main.go",
				Hooks:   hookMap,
			},
		},
	}
}

func TestLoadFromConfig_ValidHooks(t *testing.T) {
	cfg := makeConfig(map[string]string{
		"before_start": "true",
		"after_stop":   "true",
	})
	r := hooks.New()
	if err := hooks.LoadFromConfig(cfg, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := r.List("api")
	if len(events) != 2 {
		t.Fatalf("expected 2 events registered, got %d", len(events))
	}
}

func TestLoadFromConfig_UnknownEvent(t *testing.T) {
	cfg := makeConfig(map[string]string{
		"on_crash": "true",
	})
	r := hooks.New()
	if err := hooks.LoadFromConfig(cfg, r); err == nil {
		t.Fatal("expected error for unknown hook event")
	}
}

func TestLoadFromConfig_FiresRegisteredCommand(t *testing.T) {
	cfg := makeConfig(map[string]string{
		"after_start": "true", // 'true' is a no-op shell command
	})
	r := hooks.New()
	if err := hooks.LoadFromConfig(cfg, r); err != nil {
		t.Fatalf("load error: %v", err)
	}
	if err := r.Fire("api", hooks.EventAfterStart); err != nil {
		t.Fatalf("fire error: %v", err)
	}
}

func TestLoadFromConfig_EmptyHooks(t *testing.T) {
	cfg := makeConfig(nil)
	r := hooks.New()
	if err := hooks.LoadFromConfig(cfg, r); err != nil {
		t.Fatalf("unexpected error with no hooks: %v", err)
	}
	if events := r.List("api"); len(events) != 0 {
		t.Fatalf("expected no events, got %d", len(events))
	}
}
