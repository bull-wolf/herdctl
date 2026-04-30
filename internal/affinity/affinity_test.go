package affinity

import (
	"testing"
)

func TestAdd_And_Get(t *testing.T) {
	s := New()
	if err := s.Add("web", "db", "together"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rules := s.Get("web")
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Target != "db" || rules[0].Kind != "together" {
		t.Errorf("unexpected rule: %+v", rules[0])
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	if rules := s.Get("unknown"); len(rules) != 0 {
		t.Errorf("expected empty slice, got %v", rules)
	}
}

func TestAdd_InvalidArgs(t *testing.T) {
	s := New()
	if err := s.Add("", "db", "together"); err == nil {
		t.Error("expected error for empty service")
	}
	if err := s.Add("web", "", "together"); err == nil {
		t.Error("expected error for empty target")
	}
	if err := s.Add("web", "db", "invalid"); err == nil {
		t.Error("expected error for unknown kind")
	}
}

func TestAdd_Idempotent(t *testing.T) {
	s := New()
	_ = s.Add("web", "cache", "apart")
	_ = s.Add("web", "cache", "apart")
	if len(s.Get("web")) != 1 {
		t.Error("duplicate rule should not be added")
	}
}

func TestConflicts_DetectsApart(t *testing.T) {
	s := New()
	_ = s.Add("web", "legacy", "apart")
	_ = s.Add("web", "db", "together")
	conflicts := s.Conflicts("web", []string{"legacy", "db", "cache"})
	if len(conflicts) != 1 || conflicts[0] != "legacy" {
		t.Errorf("expected [legacy], got %v", conflicts)
	}
}

func TestConflicts_NoConflict(t *testing.T) {
	s := New()
	_ = s.Add("web", "legacy", "apart")
	if c := s.Conflicts("web", []string{"db"}); len(c) != 0 {
		t.Errorf("expected no conflicts, got %v", c)
	}
}

func TestRemove_ClearsRules(t *testing.T) {
	s := New()
	_ = s.Add("web", "db", "together")
	s.Remove("web")
	if rules := s.Get("web"); len(rules) != 0 {
		t.Errorf("expected empty after remove, got %v", rules)
	}
}

func TestAll_ReturnsAllRules(t *testing.T) {
	s := New()
	_ = s.Add("web", "db", "together")
	_ = s.Add("api", "cache", "apart")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 rules, got %d", len(all))
	}
}
