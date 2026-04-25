package registry

import (
	"testing"
)

func TestRegister_And_Get(t *testing.T) {
	r := New()
	err := r.Register("api", "localhost", 8080, map[string]string{"env": "dev"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := r.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Host != "localhost" || e.Port != 8080 {
		t.Errorf("unexpected entry values: %+v", e)
	}
	if e.Meta["env"] != "dev" {
		t.Errorf("expected meta env=dev, got %v", e.Meta)
	}
}

func TestRegister_InvalidArgs(t *testing.T) {
	r := New()
	if err := r.Register("", "localhost", 8080, nil); err == nil {
		t.Error("expected error for empty service name")
	}
	if err := r.Register("api", "", 8080, nil); err == nil {
		t.Error("expected error for empty host")
	}
	if err := r.Register("api", "localhost", 0, nil); err == nil {
		t.Error("expected error for invalid port")
	}
	if err := r.Register("api", "localhost", 99999, nil); err == nil {
		t.Error("expected error for out-of-range port")
	}
}

func TestGet_Missing(t *testing.T) {
	r := New()
	_, ok := r.Get("nonexistent")
	if ok {
		t.Error("expected no entry for unregistered service")
	}
}

func TestRegister_Overwrites(t *testing.T) {
	r := New()
	_ = r.Register("svc", "localhost", 3000, nil)
	_ = r.Register("svc", "127.0.0.1", 4000, nil)
	e, ok := r.Get("svc")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Host != "127.0.0.1" || e.Port != 4000 {
		t.Errorf("expected overwritten values, got %+v", e)
	}
}

func TestDeregister(t *testing.T) {
	r := New()
	_ = r.Register("worker", "localhost", 5000, nil)
	if err := r.Deregister("worker"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := r.Get("worker")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestDeregister_NotFound(t *testing.T) {
	r := New()
	if err := r.Deregister("ghost"); err == nil {
		t.Error("expected error deregistering unknown service")
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	r := New()
	_ = r.Register("a", "localhost", 1001, nil)
	_ = r.Register("b", "localhost", 1002, nil)
	all := r.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
