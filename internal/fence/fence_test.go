package fence

import (
	"testing"
)

func TestOpen_And_IsOpen(t *testing.T) {
	f := New()
	if err := f.Open("api", "deployment complete"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.IsOpen("api") {
		t.Error("expected gate to be open")
	}
}

func TestClose_And_IsOpen(t *testing.T) {
	f := New()
	_ = f.Open("api", "")
	if err := f.Close("api", "maintenance"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.IsOpen("api") {
		t.Error("expected gate to be closed")
	}
}

func TestIsOpen_DefaultsToOpen(t *testing.T) {
	f := New()
	if !f.IsOpen("unknown") {
		t.Error("expected unknown service to be open by default")
	}
}

func TestOpen_EmptyService(t *testing.T) {
	f := New()
	if err := f.Open("", "reason"); err == nil {
		t.Error("expected error for empty service")
	}
}

func TestClose_EmptyService(t *testing.T) {
	f := New()
	if err := f.Close("", "reason"); err == nil {
		t.Error("expected error for empty service")
	}
}

func TestGet_ReturnsEntry(t *testing.T) {
	f := New()
	_ = f.Close("worker", "overloaded")
	e, ok := f.Get("worker")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.State != StateClosed {
		t.Errorf("expected closed, got %s", e.State)
	}
	if e.Reason != "overloaded" {
		t.Errorf("expected reason 'overloaded', got %s", e.Reason)
	}
}

func TestGet_Missing(t *testing.T) {
	f := New()
	_, ok := f.Get("ghost")
	if ok {
		t.Error("expected no entry for unknown service")
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	f := New()
	_ = f.Open("api", "")
	_ = f.Close("worker", "paused")
	entries := f.All()
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestOpen_Overwrites(t *testing.T) {
	f := New()
	_ = f.Close("api", "maintenance")
	_ = f.Open("api", "back online")
	if !f.IsOpen("api") {
		t.Error("expected gate to be open after re-open")
	}
	e, _ := f.Get("api")
	if e.Reason != "back online" {
		t.Errorf("expected updated reason, got %s", e.Reason)
	}
}
