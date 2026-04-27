package profile

import (
	"testing"
	"time"
)

func TestRecord_And_Latest(t *testing.T) {
	p := New(10)
	if err := p.Record("api", "start", 120*time.Millisecond); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := p.Latest("api")
	if !ok {
		t.Fatal("expected entry, got none")
	}
	if e.Duration != 120*time.Millisecond {
		t.Errorf("expected 120ms, got %v", e.Duration)
	}
	if e.Operation != "start" {
		t.Errorf("expected operation 'start', got %q", e.Operation)
	}
}

func TestLatest_Missing(t *testing.T) {
	p := New(10)
	_, ok := p.Latest("ghost")
	if ok {
		t.Error("expected no entry for unknown service")
	}
}

func TestRecord_InvalidArgs(t *testing.T) {
	p := New(10)
	if err := p.Record("", "start", time.Millisecond); err == nil {
		t.Error("expected error for empty service")
	}
	if err := p.Record("api", "", time.Millisecond); err == nil {
		t.Error("expected error for empty operation")
	}
}

func TestRecord_RollingWindow(t *testing.T) {
	p := New(3)
	for i := 0; i < 5; i++ {
		_ = p.Record("svc", "op", time.Duration(i)*time.Millisecond)
	}
	h := p.History("svc")
	if len(h) != 3 {
		t.Errorf("expected 3 entries, got %d", len(h))
	}
	if h[0].Duration != 2*time.Millisecond {
		t.Errorf("expected oldest retained entry to be 2ms, got %v", h[0].Duration)
	}
}

func TestAverage_ReturnsCorrectMean(t *testing.T) {
	p := New(10)
	_ = p.Record("db", "query", 100*time.Millisecond)
	_ = p.Record("db", "query", 200*time.Millisecond)
	_ = p.Record("db", "query", 300*time.Millisecond)
	avg, ok := p.Average("db")
	if !ok {
		t.Fatal("expected average result")
	}
	if avg != 200*time.Millisecond {
		t.Errorf("expected 200ms average, got %v", avg)
	}
}

func TestAverage_Missing(t *testing.T) {
	p := New(10)
	_, ok := p.Average("unknown")
	if ok {
		t.Error("expected no average for unknown service")
	}
}

func TestClear_RemovesEntries(t *testing.T) {
	p := New(10)
	_ = p.Record("worker", "run", 50*time.Millisecond)
	p.Clear("worker")
	h := p.History("worker")
	if len(h) != 0 {
		t.Errorf("expected empty history after clear, got %d entries", len(h))
	}
}
