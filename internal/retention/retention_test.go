package retention

import (
	"testing"
	"time"
)

func TestSet_And_Get(t *testing.T) {
	m := New()
	if err := m.Set("api", 24*time.Hour, 1000); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := m.Get("api")
	if !ok {
		t.Fatal("expected policy to exist")
	}
	if p.MaxAge != 24*time.Hour || p.MaxItems != 1000 {
		t.Errorf("unexpected policy: %+v", p)
	}
}

func TestGet_Missing(t *testing.T) {
	m := New()
	_, ok := m.Get("unknown")
	if ok {
		t.Fatal("expected no policy for unknown service")
	}
}

func TestSet_InvalidArgs(t *testing.T) {
	m := New()
	if err := m.Set("", time.Hour, 10); err == nil {
		t.Fatal("expected error for empty service")
	}
	if err := m.Set("svc", -1, 10); err == nil {
		t.Fatal("expected error for negative maxAge")
	}
	if err := m.Set("svc", time.Hour, -1); err == nil {
		t.Fatal("expected error for negative maxItems")
	}
}

func TestApply_RecordsHistory(t *testing.T) {
	m := New()
	_ = m.Set("worker", time.Hour, 500)
	if err := m.Apply("worker", 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h := m.History()
	if len(h) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(h))
	}
	if h[0].Service != "worker" || h[0].Purged != 42 {
		t.Errorf("unexpected history entry: %+v", h[0])
	}
}

func TestApply_EmptyService(t *testing.T) {
	m := New()
	if err := m.Apply("", 10); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestAll_ReturnsPolicies(t *testing.T) {
	m := New()
	_ = m.Set("a", time.Minute, 100)
	_ = m.Set("b", time.Hour, 200)
	all := m.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 policies, got %d", len(all))
	}
}

func TestHistory_ReturnsCopy(t *testing.T) {
	m := New()
	_ = m.Apply("svc", 5)
	h := m.History()
	h[0].Purged = 999
	if m.History()[0].Purged == 999 {
		t.Fatal("history should return a copy")
	}
}
