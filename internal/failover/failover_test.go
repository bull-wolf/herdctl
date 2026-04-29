package failover

import (
	"testing"
)

func TestRegister_And_Resolve(t *testing.T) {
	m := New()
	if err := m.Register("api", []string{"api-backup"}, StrategyPrimary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, active := m.Resolve("api")
	if active {
		t.Error("expected failover to be inactive before trigger")
	}
}

func TestTrigger_Primary(t *testing.T) {
	m := New()
	_ = m.Register("api", []string{"api-b", "api-c"}, StrategyPrimary)
	target, err := m.Trigger("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if target != "api-b" {
		t.Errorf("expected api-b, got %s", target)
	}
	// primary always returns first target
	target2, _ := m.Trigger("api")
	if target2 != "api-b" {
		t.Errorf("expected api-b on second trigger, got %s", target2)
	}
}

func TestTrigger_RoundRobin(t *testing.T) {
	m := New()
	_ = m.Register("worker", []string{"w1", "w2", "w3"}, StrategyRoundRobin)
	expected := []string{"w1", "w2", "w3", "w1"}
	for i, want := range expected {
		got, err := m.Trigger("worker")
		if err != nil {
			t.Fatalf("step %d: unexpected error: %v", i, err)
		}
		if got != want {
			t.Errorf("step %d: expected %s, got %s", i, want, got)
		}
	}
}

func TestTrigger_UnknownService(t *testing.T) {
	m := New()
	_, err := m.Trigger("ghost")
	if err == nil {
		t.Error("expected error for unknown service")
	}
}

func TestRecover_ClearsState(t *testing.T) {
	m := New()
	_ = m.Register("db", []string{"db-replica"}, StrategyPrimary)
	_, _ = m.Trigger("db")
	if err := m.Recover("db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, active := m.Resolve("db")
	if active {
		t.Error("expected failover to be inactive after recover")
	}
}

func TestRegister_InvalidArgs(t *testing.T) {
	m := New()
	if err := m.Register("", []string{"x"}, StrategyPrimary); err == nil {
		t.Error("expected error for empty service")
	}
	if err := m.Register("svc", []string{}, StrategyPrimary); err == nil {
		t.Error("expected error for empty targets")
	}
}

func TestAll_ReturnsEntries(t *testing.T) {
	m := New()
	_ = m.Register("a", []string{"a1"}, StrategyPrimary)
	_ = m.Register("b", []string{"b1"}, StrategyRoundRobin)
	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
