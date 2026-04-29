package limiter

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	l := New()
	if err := l.Set("api", 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := l.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Max != 3 {
		t.Errorf("expected max=3, got %d", e.Max)
	}
}

func TestSet_InvalidArgs(t *testing.T) {
	l := New()
	if err := l.Set("", 3); err == nil {
		t.Error("expected error for empty service")
	}
	if err := l.Set("api", 0); err == nil {
		t.Error("expected error for zero max")
	}
	if err := l.Set("api", -1); err == nil {
		t.Error("expected error for negative max")
	}
}

func TestAcquire_And_Release(t *testing.T) {
	l := New()
	_ = l.Set("api", 2)

	if err := l.Acquire("api"); err != nil {
		t.Fatalf("unexpected error on first acquire: %v", err)
	}
	if err := l.Acquire("api"); err != nil {
		t.Fatalf("unexpected error on second acquire: %v", err)
	}
	if err := l.Acquire("api"); err == nil {
		t.Error("expected error when limit exceeded")
	}

	_ = l.Release("api")
	if err := l.Acquire("api"); err != nil {
		t.Errorf("expected acquire to succeed after release: %v", err)
	}
}

func TestAcquire_NoLimit(t *testing.T) {
	l := New()
	if err := l.Acquire("unknown"); err == nil {
		t.Error("expected error for unconfigured service")
	}
}

func TestRelease_NoLimit(t *testing.T) {
	l := New()
	if err := l.Release("unknown"); err == nil {
		t.Error("expected error for unconfigured service")
	}
}

func TestGet_Missing(t *testing.T) {
	l := New()
	_, ok := l.Get("missing")
	if ok {
		t.Error("expected missing entry")
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	l := New()
	_ = l.Set("api", 2)
	_ = l.Set("worker", 5)
	all := l.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
