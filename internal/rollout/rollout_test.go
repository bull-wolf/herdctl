package rollout

import (
	"errors"
	"testing"
)

func TestBegin_And_Get(t *testing.T) {
	m := New()
	if err := m.Begin("api", StrategySequential); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := m.Get("api")
	if !ok {
		t.Fatal("expected state to exist")
	}
	if s.Service != "api" {
		t.Errorf("expected service api, got %s", s.Service)
	}
	if s.Strategy != StrategySequential {
		t.Errorf("expected sequential strategy, got %s", s.Strategy)
	}
	if s.Done {
		t.Error("expected Done to be false")
	}
}

func TestBegin_InvalidArgs(t *testing.T) {
	m := New()
	if err := m.Begin("", StrategyAll); err == nil {
		t.Error("expected error for empty service")
	}
	if err := m.Begin("svc", Strategy("unknown")); err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestComplete_MarksSuccess(t *testing.T) {
	m := New()
	_ = m.Begin("worker", StrategyAll)
	if err := m.Complete("worker", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, _ := m.Get("worker")
	if !s.Done {
		t.Error("expected Done to be true")
	}
	if s.Err != nil {
		t.Errorf("expected nil error, got %v", s.Err)
	}
}

func TestComplete_MarksFailure(t *testing.T) {
	m := New()
	_ = m.Begin("db", StrategyCanary)
	expected := errors.New("deploy failed")
	_ = m.Complete("db", expected)
	s, _ := m.Get("db")
	if s.Err == nil || s.Err.Error() != expected.Error() {
		t.Errorf("expected error %v, got %v", expected, s.Err)
	}
}

func TestComplete_MissingService(t *testing.T) {
	m := New()
	if err := m.Complete("ghost", nil); err == nil {
		t.Error("expected error for unknown service")
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, ok := m.Get("nonexistent")
	if ok {
		t.Error("expected ok=false for missing service")
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	m := New()
	_ = m.Begin("svc-a", StrategyAll)
	_ = m.Begin("svc-b", StrategyCanary)
	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
