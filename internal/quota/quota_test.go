package quota

import (
	"strings"
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	m := New()
	if err := m.Set("api", KindCPU, 2.0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := m.Get("api", KindCPU)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Limit != 2.0 {
		t.Errorf("expected limit 2.0, got %.2f", e.Limit)
	}
}

func TestSet_InvalidArgs(t *testing.T) {
	m := New()
	if err := m.Set("", KindCPU, 1.0); err == nil {
		t.Error("expected error for empty service")
	}
	if err := m.Set("api", KindCPU, -1.0); err == nil {
		t.Error("expected error for non-positive limit")
	}
	if err := m.Set("api", KindCPU, 0.0); err == nil {
		t.Error("expected error for zero limit")
	}
}

func TestRecord_UnderLimit(t *testing.T) {
	m := New()
	_ = m.Set("api", KindMemory, 512.0)
	if err := m.Record("api", KindMemory, 100.0); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	e, _ := m.Get("api", KindMemory)
	if e.Used != 100.0 {
		t.Errorf("expected used=100.0, got %.2f", e.Used)
	}
}

func TestRecord_ExceedsLimit(t *testing.T) {
	m := New()
	_ = m.Set("worker", KindProcs, 3.0)
	_ = m.Record("worker", KindProcs, 2.0)
	err := m.Record("worker", KindProcs, 2.0)
	if err == nil {
		t.Fatal("expected quota exceeded error")
	}
	if !strings.Contains(err.Error(), "exceeded") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRecord_NoQuota(t *testing.T) {
	m := New()
	if err := m.Record("ghost", KindCPU, 1.0); err == nil {
		t.Error("expected error for unknown service")
	}
}

func TestRecord_InvalidAmount(t *testing.T) {
	m := New()
	_ = m.Set("api", KindCPU, 4.0)
	if err := m.Record("api", KindCPU, 0.0); err == nil {
		t.Error("expected error for zero amount")
	}
	if err := m.Record("api", KindCPU, -1.0); err == nil {
		t.Error("expected error for negative amount")
	}
}

func TestReset_ClearsUsage(t *testing.T) {
	m := New()
	_ = m.Set("api", KindCPU, 4.0)
	_ = m.Record("api", KindCPU, 3.0)
	m.Reset("api", KindCPU)
	e, _ := m.Get("api", KindCPU)
	if e.Used != 0 {
		t.Errorf("expected used=0 after reset, got %.2f", e.Used)
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	m := New()
	_ = m.Set("api", KindCPU, 2.0)
	_ = m.Set("api", KindMemory, 256.0)
	_ = m.Set("worker", KindProcs, 5.0)
	all := m.All()
	if len(all) != 3 {
		t.Errorf("expected 3 entries, got %d", len(all))
	}
}

func TestExceeded_Flag(t *testing.T) {
	e := Entry{Limit: 10.0, Used: 10.0}
	if !e.Exceeded() {
		t.Error("expected Exceeded() to be true")
	}
	e.Used = 9.9
	if e.Exceeded() {
		t.Error("expected Exceeded() to be false")
	}
}
