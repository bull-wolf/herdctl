package backoff_test

import (
	"testing"
	"time"

	"github.com/herdctl/herdctl/internal/backoff"
)

func TestNext_DefaultExponential(t *testing.T) {
	m := backoff.New()
	w1, err := m.Next("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w2, _ := m.Next("svc")
	if w2 <= w1 {
		t.Errorf("expected w2 (%v) > w1 (%v) for exponential backoff", w2, w1)
	}
}

func TestNext_FixedStrategy(t *testing.T) {
	m := backoff.New()
	_ = m.Configure("svc", backoff.StrategyFixed, 200*time.Millisecond, 10*time.Second)
	w1, _ := m.Next("svc")
	w2, _ := m.Next("svc")
	if w1 != w2 {
		t.Errorf("expected fixed wait, got w1=%v w2=%v", w1, w2)
	}
	if w1 != 200*time.Millisecond {
		t.Errorf("expected 200ms, got %v", w1)
	}
}

func TestNext_CapsAtMaxWait(t *testing.T) {
	m := backoff.New()
	_ = m.Configure("svc", backoff.StrategyExponential, 1*time.Second, 3*time.Second)
	var last time.Duration
	for i := 0; i < 10; i++ {
		last, _ = m.Next("svc")
	}
	if last != 3*time.Second {
		t.Errorf("expected max wait 3s, got %v", last)
	}
}

func TestReset_ClearsAttempts(t *testing.T) {
	m := backoff.New()
	_, _ = m.Next("svc")
	_, _ = m.Next("svc")
	m.Reset("svc")
	e, ok := m.Get("svc")
	if !ok {
		t.Fatal("expected entry to exist after reset")
	}
	if e.Attempts != 0 {
		t.Errorf("expected 0 attempts after reset, got %d", e.Attempts)
	}
}

func TestGet_Missing(t *testing.T) {
	m := backoff.New()
	_, ok := m.Get("unknown")
	if ok {
		t.Error("expected false for missing service")
	}
}

func TestGet_ReturnsEntry(t *testing.T) {
	m := backoff.New()
	_ = m.Configure("svc", backoff.StrategyFixed, 100*time.Millisecond, 5*time.Second)
	_, _ = m.Next("svc")
	e, ok := m.Get("svc")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", e.Attempts)
	}
	if e.LastWait != 100*time.Millisecond {
		t.Errorf("expected LastWait 100ms, got %v", e.LastWait)
	}
}

func TestConfigure_InvalidArgs(t *testing.T) {
	m := backoff.New()
	if err := m.Configure("", backoff.StrategyFixed, 100*time.Millisecond, 1*time.Second); err == nil {
		t.Error("expected error for empty service")
	}
	if err := m.Configure("svc", backoff.StrategyFixed, 0, 1*time.Second); err == nil {
		t.Error("expected error for zero base")
	}
	if err := m.Configure("svc", backoff.StrategyFixed, 100*time.Millisecond, 0); err == nil {
		t.Error("expected error for zero maxWait")
	}
}

func TestNext_EmptyService(t *testing.T) {
	m := backoff.New()
	_, err := m.Next("")
	if err == nil {
		t.Error("expected error for empty service")
	}
}
