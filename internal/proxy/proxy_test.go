package proxy

import (
	"testing"
)

func TestRegister_And_Get(t *testing.T) {
	p := New()
	if err := p.Register("api", 8080, "localhost:9090"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := p.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Port != 8080 || e.Upstream != "localhost:9090" || !e.Enabled {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestRegister_InvalidArgs(t *testing.T) {
	p := New()
	if err := p.Register("", 8080, "localhost:9090"); err == nil {
		t.Error("expected error for empty service")
	}
	if err := p.Register("api", 0, "localhost:9090"); err == nil {
		t.Error("expected error for invalid port")
	}
	if err := p.Register("api", 8080, ""); err == nil {
		t.Error("expected error for empty upstream")
	}
}

func TestGet_Missing(t *testing.T) {
	p := New()
	_, ok := p.Get("unknown")
	if ok {
		t.Error("expected no entry for unknown service")
	}
}

func TestDisable_And_Enable(t *testing.T) {
	p := New()
	_ = p.Register("api", 8080, "localhost:9090")

	if err := p.Disable("api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := p.Get("api")
	if e.Enabled {
		t.Error("expected entry to be disabled")
	}

	if err := p.Enable("api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ = p.Get("api")
	if !e.Enabled {
		t.Error("expected entry to be enabled")
	}
}

func TestDisable_Missing(t *testing.T) {
	p := New()
	if err := p.Disable("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	p := New()
	_ = p.Register("api", 8080, "localhost:9090")
	_ = p.Register("web", 3000, "localhost:3001")
	all := p.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestDeregister(t *testing.T) {
	p := New()
	_ = p.Register("api", 8080, "localhost:9090")
	if err := p.Deregister("api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := p.Get("api")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestDeregister_Missing(t *testing.T) {
	p := New()
	if err := p.Deregister("ghost"); err == nil {
		t.Error("expected error for missing service")
	}
}
