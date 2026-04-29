package cascade

import (
	"testing"
)

func TestRegister_And_Get(t *testing.T) {
	m := New()
	if err := m.Register("api", PolicyStop, []string{"worker", "cache"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r, ok := m.Get("api")
	if !ok {
		t.Fatal("expected rule to exist")
	}
	if r.Policy != PolicyStop {
		t.Errorf("expected stop, got %s", r.Policy)
	}
	if len(r.Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(r.Targets))
	}
}

func TestRegister_InvalidArgs(t *testing.T) {
	m := New()
	if err := m.Register("", PolicyStop, nil); err == nil {
		t.Error("expected error for empty service")
	}
	if err := m.Register("api", Policy("unknown"), nil); err == nil {
		t.Error("expected error for unknown policy")
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, ok := m.Get("ghost")
	if ok {
		t.Error("expected no rule for unknown service")
	}
}

func TestEvaluate_ReturnsTargets(t *testing.T) {
	m := New()
	_ = m.Register("db", PolicyRestart, []string{"api", "worker"})
	policy, targets := m.Evaluate("db")
	if policy != PolicyRestart {
		t.Errorf("expected restart, got %s", policy)
	}
	if len(targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(targets))
	}
}

func TestEvaluate_IgnorePolicy(t *testing.T) {
	m := New()
	_ = m.Register("sidecar", PolicyIgnore, []string{"api"})
	policy, targets := m.Evaluate("sidecar")
	if policy != PolicyIgnore {
		t.Errorf("expected ignore, got %s", policy)
	}
	if len(targets) != 0 {
		t.Errorf("expected no targets for ignore policy, got %d", len(targets))
	}
}

func TestEvaluate_NoRule(t *testing.T) {
	m := New()
	policy, targets := m.Evaluate("missing")
	if policy != PolicyIgnore {
		t.Errorf("expected ignore for missing rule, got %s", policy)
	}
	if len(targets) != 0 {
		t.Error("expected empty targets for missing rule")
	}
}

func TestDeregister(t *testing.T) {
	m := New()
	_ = m.Register("api", PolicyStop, []string{"worker"})
	m.Deregister("api")
	_, ok := m.Get("api")
	if ok {
		t.Error("expected rule to be removed")
	}
}

func TestAll_ReturnsAllRules(t *testing.T) {
	m := New()
	_ = m.Register("a", PolicyStop, nil)
	_ = m.Register("b", PolicyRestart, nil)
	if len(m.All()) != 2 {
		t.Errorf("expected 2 rules, got %d", len(m.All()))
	}
}
