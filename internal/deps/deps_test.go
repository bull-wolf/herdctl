package deps

import (
	"testing"
)

func TestOrder_NoDeps(t *testing.T) {
	g := New(map[string][]string{
		"api": {},
		"db":  {},
	})
	order, err := g.Order([]string{"api", "db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 2 {
		t.Fatalf("expected 2 services, got %d", len(order))
	}
}

func TestOrder_WithDeps(t *testing.T) {
	g := New(map[string][]string{
		"api":    {"db"},
		"worker": {"api", "db"},
		"db":     {},
	})
	order, err := g.Order([]string{"worker"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pos := func(name string) int {
		for i, s := range order {
			if s == name {
				return i
			}
		}
		return -1
	}
	if pos("db") >= pos("api") {
		t.Errorf("db should come before api")
	}
	if pos("api") >= pos("worker") {
		t.Errorf("api should come before worker")
	}
}

func TestOrder_CycleDetected(t *testing.T) {
	g := New(map[string][]string{
		"a": {"b"},
		"b": {"a"},
	})
	_, err := g.Order([]string{"a"})
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestDeps_ReturnsDirect(t *testing.T) {
	g := New(map[string][]string{
		"api": {"db", "cache"},
	})
	deps := g.Deps("api")
	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %d", len(deps))
	}
}
