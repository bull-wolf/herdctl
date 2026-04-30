package priority

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	m := New()
	if err := m.Set("api", High); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := m.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Level != High {
		t.Errorf("expected High, got %d", e.Level)
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, ok := m.Get("unknown")
	if ok {
		t.Fatal("expected no entry for unknown service")
	}
}

func TestSet_InvalidArgs(t *testing.T) {
	m := New()
	if err := m.Set("", High); err == nil {
		t.Fatal("expected error for empty service")
	}
	if err := m.Set("svc", Level(-1)); err == nil {
		t.Fatal("expected error for negative level")
	}
}

func TestGetLevel_DefaultsToNormal(t *testing.T) {
	m := New()
	if l := m.GetLevel("unset"); l != Normal {
		t.Errorf("expected Normal (%d), got %d", Normal, l)
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	m := New()
	_ = m.Set("a", Low)
	_ = m.Set("b", High)
	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestSorted_OrdersByPriority(t *testing.T) {
	m := New()
	_ = m.Set("db", Low)
	_ = m.Set("cache", High)
	_ = m.Set("api", Normal)

	sorted := m.Sorted([]string{"db", "api", "cache"})
	if sorted[0] != "cache" {
		t.Errorf("expected cache first, got %s", sorted[0])
	}
	if sorted[1] != "api" {
		t.Errorf("expected api second, got %s", sorted[1])
	}
	if sorted[2] != "db" {
		t.Errorf("expected db third, got %s", sorted[2])
	}
}

func TestSorted_UnknownUsesNormal(t *testing.T) {
	m := New()
	_ = m.Set("db", Low)
	// "api" has no explicit priority — defaults to Normal
	sorted := m.Sorted([]string{"db", "api"})
	if sorted[0] != "api" {
		t.Errorf("expected api first (Normal > Low), got %s", sorted[0])
	}
}
