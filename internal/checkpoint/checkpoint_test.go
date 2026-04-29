package checkpoint

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	m := New()
	err := m.Set("api", "db-ready", StatusPassed, "connected")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := m.Get("api", "db-ready")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Status != StatusPassed || e.Message != "connected" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, ok := m.Get("api", "nonexistent")
	if ok {
		t.Fatal("expected no entry")
	}
}

func TestSet_InvalidArgs(t *testing.T) {
	m := New()
	if err := m.Set("", "name", StatusPending, ""); err == nil {
		t.Error("expected error for empty service")
	}
	if err := m.Set("svc", "", StatusPending, ""); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestSet_Overwrites(t *testing.T) {
	m := New()
	_ = m.Set("api", "db-ready", StatusPending, "waiting")
	_ = m.Set("api", "db-ready", StatusFailed, "timeout")
	e, _ := m.Get("api", "db-ready")
	if e.Status != StatusFailed {
		t.Errorf("expected failed, got %s", e.Status)
	}
}

func TestAllForService(t *testing.T) {
	m := New()
	_ = m.Set("api", "db-ready", StatusPassed, "ok")
	_ = m.Set("api", "cache-ready", StatusPending, "")
	_ = m.Set("worker", "queue-ready", StatusPassed, "ok")
	entries := m.AllForService("api")
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestClear_RemovesService(t *testing.T) {
	m := New()
	_ = m.Set("api", "db-ready", StatusPassed, "ok")
	m.Clear("api")
	entries := m.AllForService("api")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after clear, got %d", len(entries))
	}
}

func TestAllPassed_True(t *testing.T) {
	m := New()
	_ = m.Set("api", "db-ready", StatusPassed, "ok")
	_ = m.Set("api", "cache-ready", StatusPassed, "ok")
	if !m.AllPassed("api") {
		t.Error("expected all passed")
	}
}

func TestAllPassed_False(t *testing.T) {
	m := New()
	_ = m.Set("api", "db-ready", StatusPassed, "ok")
	_ = m.Set("api", "cache-ready", StatusPending, "")
	if m.AllPassed("api") {
		t.Error("expected not all passed")
	}
}

func TestAllPassed_NoEntries(t *testing.T) {
	m := New()
	if m.AllPassed("unknown") {
		t.Error("expected false for unknown service")
	}
}
