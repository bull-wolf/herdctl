package circuit_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/herdctl/internal/circuit"
)

func TestAllow_ClosedByDefault(t *testing.T) {
	b := circuit.New(3, time.Second)
	if err := b.Allow("svc"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensCircuit(t *testing.T) {
	b := circuit.New(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure("svc", errors.New("boom"))
	}
	e, ok := b.Get("svc")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.State != circuit.StateOpen {
		t.Errorf("expected open, got %s", e.State)
	}
}

func TestAllow_BlocksWhenOpen(t *testing.T) {
	b := circuit.New(2, 10*time.Second)
	b.RecordFailure("svc", errors.New("e"))
	b.RecordFailure("svc", errors.New("e"))
	if err := b.Allow("svc"); err == nil {
		t.Fatal("expected error when circuit is open")
	}
}

func TestAllow_HalfOpenAfterCooldown(t *testing.T) {
	b := circuit.New(1, 10*time.Millisecond)
	b.RecordFailure("svc", errors.New("e"))
	time.Sleep(20 * time.Millisecond)
	if err := b.Allow("svc"); err != nil {
		t.Fatalf("expected nil after cooldown, got %v", err)
	}
	e, _ := b.Get("svc")
	if e.State != circuit.StateHalfOpen {
		t.Errorf("expected half-open, got %s", e.State)
	}
}

func TestRecordSuccess_ClosesCircuit(t *testing.T) {
	b := circuit.New(2, time.Second)
	b.RecordFailure("svc", errors.New("e"))
	b.RecordFailure("svc", errors.New("e"))
	b.RecordSuccess("svc")
	e, _ := b.Get("svc")
	if e.State != circuit.StateClosed {
		t.Errorf("expected closed, got %s", e.State)
	}
	if e.Failures != 0 {
		t.Errorf("expected 0 failures, got %d", e.Failures)
	}
}

func TestGet_Missing(t *testing.T) {
	b := circuit.New(3, time.Second)
	_, ok := b.Get("missing")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestReset_ClearsEntry(t *testing.T) {
	b := circuit.New(1, time.Second)
	b.RecordFailure("svc", errors.New("e"))
	b.Reset("svc")
	_, ok := b.Get("svc")
	if ok {
		t.Fatal("expected entry to be cleared after reset")
	}
}

func TestAllow_BelowThreshold(t *testing.T) {
	b := circuit.New(5, time.Second)
	b.RecordFailure("svc", errors.New("e"))
	b.RecordFailure("svc", errors.New("e"))
	if err := b.Allow("svc"); err != nil {
		t.Fatalf("expected nil below threshold, got %v", err)
	}
}
