package heartbeat

import (
	"testing"
	"time"
)

func TestRegister_And_Get(t *testing.T) {
	m := New()
	if err := m.Register("api", time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, err := m.Get("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Service != "api" {
		t.Errorf("expected service 'api', got %q", e.Service)
	}
	if e.Interval != time.Second {
		t.Errorf("expected interval 1s, got %v", e.Interval)
	}
}

func TestRegister_InvalidArgs(t *testing.T) {
	m := New()
	if err := m.Register("", time.Second); err == nil {
		t.Error("expected error for empty service")
	}
	if err := m.Register("svc", 0); err == nil {
		t.Error("expected error for zero interval")
	}
}

func TestBeat_ResetsMissed(t *testing.T) {
	m := New()
	_ = m.Register("worker", 50*time.Millisecond)

	// Force a stale check to increment missed
	m.Check(time.Now().Add(time.Second))

	e, _ := m.Get("worker")
	if e.Missed == 0 {
		t.Fatal("expected missed > 0 after stale check")
	}

	if err := m.Beat("worker"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ = m.Get("worker")
	if e.Missed != 0 {
		t.Errorf("expected missed to be reset to 0, got %d", e.Missed)
	}
}

func TestBeat_UnknownService(t *testing.T) {
	m := New()
	if err := m.Beat("ghost"); err == nil {
		t.Error("expected error for unknown service")
	}
}

func TestCheck_DetectsStale(t *testing.T) {
	m := New()
	_ = m.Register("db", 100*time.Millisecond)

	future := time.Now().Add(time.Second)
	stale := m.Check(future)
	if len(stale) != 1 || stale[0] != "db" {
		t.Errorf("expected [db] to be stale, got %v", stale)
	}

	e, _ := m.Get("db")
	if e.Missed != 1 {
		t.Errorf("expected missed=1, got %d", e.Missed)
	}
}

func TestCheck_NoStale(t *testing.T) {
	m := New()
	_ = m.Register("cache", 10*time.Second)

	stale := m.Check(time.Now())
	if len(stale) != 0 {
		t.Errorf("expected no stale services, got %v", stale)
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, err := m.Get("missing")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	m := New()
	_ = m.Register("a", time.Second)
	_ = m.Register("b", 2*time.Second)

	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
