package cooldown

import (
	"testing"
	"time"
)

func TestSet_And_IsActive(t *testing.T) {
	m := New()
	if err := m.Set("svc", 500*time.Millisecond); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.IsActive("svc") {
		t.Fatal("expected cooldown to be active immediately after Set")
	}
}

func TestIsActive_Expired(t *testing.T) {
	m := New()
	_ = m.Set("svc", 10*time.Millisecond)
	time.Sleep(30 * time.Millisecond)
	if m.IsActive("svc") {
		t.Fatal("expected cooldown to have expired")
	}
}

func TestIsActive_Missing(t *testing.T) {
	m := New()
	if m.IsActive("ghost") {
		t.Fatal("expected false for unknown service")
	}
}

func TestSet_InvalidArgs(t *testing.T) {
	m := New()
	if err := m.Set("", time.Second); err == nil {
		t.Fatal("expected error for empty service")
	}
	if err := m.Set("svc", 0); err == nil {
		t.Fatal("expected error for zero duration")
	}
	if err := m.Set("svc", -time.Second); err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestGet_ReturnsEntry(t *testing.T) {
	m := New()
	_ = m.Set("api", time.Minute)
	e, ok := m.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Service != "api" {
		t.Fatalf("expected service 'api', got %q", e.Service)
	}
	if e.Duration != time.Minute {
		t.Fatalf("expected duration 1m, got %v", e.Duration)
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, ok := m.Get("nope")
	if ok {
		t.Fatal("expected ok=false for missing service")
	}
}

func TestClear_RemovesEntry(t *testing.T) {
	m := New()
	_ = m.Set("svc", time.Minute)
	m.Clear("svc")
	if m.IsActive("svc") {
		t.Fatal("expected cooldown to be inactive after Clear")
	}
	_, ok := m.Get("svc")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	m := New()
	_ = m.Set("a", time.Minute)
	_ = m.Set("b", time.Minute)
	all := m.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestRemainingTime_Positive(t *testing.T) {
	m := New()
	_ = m.Set("svc", time.Minute)
	e, _ := m.Get("svc")
	if e.RemainingTime() <= 0 {
		t.Fatal("expected positive remaining time")
	}
}

func TestRemainingTime_ZeroAfterExpiry(t *testing.T) {
	m := New()
	_ = m.Set("svc", 10*time.Millisecond)
	time.Sleep(30 * time.Millisecond)
	e, _ := m.Get("svc")
	if e.RemainingTime() != 0 {
		t.Fatalf("expected 0 remaining time, got %v", e.RemainingTime())
	}
}
