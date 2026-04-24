package restart

import (
	"testing"
	"time"
)

func TestShouldRestart_Never(t *testing.T) {
	m := New()
	m.Register("svc", PolicyNever, 3, 0)
	ok, err := m.ShouldRestart("svc", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected no restart for policy=never")
	}
}

func TestShouldRestart_Always(t *testing.T) {
	m := New()
	m.Register("svc", PolicyAlways, 5, 0)
	ok, err := m.ShouldRestart("svc", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected restart for policy=always")
	}
}

func TestShouldRestart_OnFail_NoFailure(t *testing.T) {
	m := New()
	m.Register("svc", PolicyOnFail, 5, 0)
	ok, _ := m.ShouldRestart("svc", false)
	if ok {
		t.Fatal("expected no restart when exit was clean")
	}
}

func TestShouldRestart_OnFail_WithFailure(t *testing.T) {
	m := New()
	m.Register("svc", PolicyOnFail, 5, 0)
	ok, err := m.ShouldRestart("svc", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected restart on failure")
	}
}

func TestShouldRestart_MaxRetryExceeded(t *testing.T) {
	m := New()
	m.Register("svc", PolicyAlways, 2, 0)
	_ = m.Record("svc")
	_ = m.Record("svc")
	ok, _ := m.ShouldRestart("svc", true)
	if ok {
		t.Fatal("expected no restart after max retries exceeded")
	}
}

func TestShouldRestart_CooldownActive(t *testing.T) {
	m := New()
	m.Register("svc", PolicyAlways, 10, 5*time.Second)
	_ = m.Record("svc") // sets LastAt to now
	ok, _ := m.ShouldRestart("svc", true)
	if ok {
		t.Fatal("expected no restart while cooldown is active")
	}
}

func TestRecord_And_Reset(t *testing.T) {
	m := New()
	m.Register("svc", PolicyAlways, 5, 0)
	_ = m.Record("svc")
	_ = m.Record("svc")
	e, ok := m.Get("svc")
	if !ok || e.Attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", e.Attempts)
	}
	_ = m.Reset("svc")
	e, _ = m.Get("svc")
	if e.Attempts != 0 {
		t.Fatalf("expected 0 attempts after reset, got %d", e.Attempts)
	}
}

func TestShouldRestart_UnknownService(t *testing.T) {
	m := New()
	_, err := m.ShouldRestart("ghost", true)
	if err == nil {
		t.Fatal("expected error for unknown service")
	}
}
