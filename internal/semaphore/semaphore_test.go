package semaphore

import (
	"testing"
)

func TestSet_And_State(t *testing.T) {
	s := New()
	if err := s.Set("api", 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	active, max, ok := s.State("api")
	if !ok {
		t.Fatal("expected state to exist")
	}
	if active != 0 || max != 3 {
		t.Errorf("expected active=0 max=3, got active=%d max=%d", active, max)
	}
}

func TestSet_InvalidArgs(t *testing.T) {
	s := New()
	if err := s.Set("", 1); err == nil {
		t.Error("expected error for empty service")
	}
	if err := s.Set("api", 0); err == nil {
		t.Error("expected error for zero max")
	}
}

func TestAcquire_And_Release(t *testing.T) {
	s := New()
	_ = s.Set("worker", 2)

	if err := s.Acquire("worker"); err != nil {
		t.Fatalf("unexpected error on first acquire: %v", err)
	}
	if err := s.Acquire("worker"); err != nil {
		t.Fatalf("unexpected error on second acquire: %v", err)
	}

	active, _, _ := s.State("worker")
	if active != 2 {
		t.Errorf("expected active=2, got %d", active)
	}

	_ = s.Release("worker")
	active, _, _ = s.State("worker")
	if active != 1 {
		t.Errorf("expected active=1 after release, got %d", active)
	}
}

func TestAcquire_LimitReached(t *testing.T) {
	s := New()
	_ = s.Set("db", 1)
	_ = s.Acquire("db")

	if err := s.Acquire("db"); err == nil {
		t.Error("expected error when limit is reached")
	}
}

func TestAcquire_UnknownService(t *testing.T) {
	s := New()
	if err := s.Acquire("ghost"); err == nil {
		t.Error("expected error for unknown service")
	}
}

func TestRelease_UnknownService(t *testing.T) {
	s := New()
	if err := s.Release("ghost"); err == nil {
		t.Error("expected error for unknown service")
	}
}

func TestState_Missing(t *testing.T) {
	s := New()
	_, _, ok := s.State("missing")
	if ok {
		t.Error("expected ok=false for missing service")
	}
}
