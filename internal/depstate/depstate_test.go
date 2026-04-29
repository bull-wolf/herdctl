package depstate

import (
	"testing"
)

func TestEvaluate_NoBlocking(t *testing.T) {
	m := New()
	state := m.Evaluate("api", []string{"db", "cache"}, func(dep string) bool {
		return true
	})
	if state != StateReady {
		t.Fatalf("expected ready, got %s", state)
	}
}

func TestEvaluate_Blocked(t *testing.T) {
	m := New()
	state := m.Evaluate("api", []string{"db", "cache"}, func(dep string) bool {
		return dep == "db" // cache not ready
	})
	if state != StateBlocked {
		t.Fatalf("expected blocked, got %s", state)
	}
	e, err := m.Get("api")
	if err != nil {
		t.Fatal(err)
	}
	if len(e.Blocking) != 1 || e.Blocking[0] != "cache" {
		t.Fatalf("unexpected blocking: %v", e.Blocking)
	}
}

func TestEvaluate_NoDeps(t *testing.T) {
	m := New()
	state := m.Evaluate("worker", []string{}, func(dep string) bool {
		return false
	})
	if state != StateReady {
		t.Fatalf("expected ready with no deps, got %s", state)
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, err := m.Get("unknown")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	m := New()
	m.Evaluate("svc-a", []string{}, func(string) bool { return true })
	m.Evaluate("svc-b", []string{"svc-a"}, func(string) bool { return false })

	all := m.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestClear_RemovesEntry(t *testing.T) {
	m := New()
	m.Evaluate("api", []string{}, func(string) bool { return true })
	m.Clear("api")
	_, err := m.Get("api")
	if err == nil {
		t.Fatal("expected error after clear")
	}
}

func TestEvaluate_OverwritesPrevious(t *testing.T) {
	m := New()
	m.Evaluate("api", []string{"db"}, func(string) bool { return false })
	state := m.Evaluate("api", []string{"db"}, func(string) bool { return true })
	if state != StateReady {
		t.Fatalf("expected ready after re-evaluation, got %s", state)
	}
	e, _ := m.Get("api")
	if len(e.Blocking) != 0 {
		t.Fatalf("expected no blocking after re-evaluation")
	}
}
