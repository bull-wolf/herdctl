package timeout

import (
	"testing"
	"time"
)

func TestSet_And_Get(t *testing.T) {
	m := New()
	err := m.Set("api", 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := m.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Service != "api" {
		t.Errorf("expected service 'api', got %q", e.Service)
	}
	if e.Duration != 5*time.Second {
		t.Errorf("expected duration 5s, got %v", e.Duration)
	}
}

func TestSet_InvalidArgs(t *testing.T) {
	m := New()
	if err := m.Set("", time.Second); err == nil {
		t.Error("expected error for empty service name")
	}
	if err := m.Set("api", 0); err == nil {
		t.Error("expected error for zero duration")
	}
	if err := m.Set("api", -time.Second); err == nil {
		t.Error("expected error for negative duration")
	}
}

func TestCheck_NotExpired(t *testing.T) {
	m := New()
	_ = m.Set("api", 10*time.Second)
	expired, err := m.Check("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expired {
		t.Error("expected not expired")
	}
}

func TestCheck_Expired(t *testing.T) {
	m := New()
	_ = m.Set("api", 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	expired, err := m.Check("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !expired {
		t.Error("expected expired")
	}
}

func TestCheck_Missing(t *testing.T) {
	m := New()
	_, err := m.Check("unknown")
	if err == nil {
		t.Error("expected error for missing service")
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, ok := m.Get("ghost")
	if ok {
		t.Error("expected missing entry")
	}
}

func TestClear_RemovesEntry(t *testing.T) {
	m := New()
	_ = m.Set("api", time.Second)
	m.Clear("api")
	_, ok := m.Get("api")
	if ok {
		t.Error("expected entry to be removed after clear")
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	m := New()
	_ = m.Set("api", time.Second)
	_ = m.Set("worker", 2*time.Second)
	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
