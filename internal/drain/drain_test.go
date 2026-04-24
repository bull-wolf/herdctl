package drain_test

import (
	"testing"
	"time"

	"github.com/user/herdctl/internal/drain"
)

func TestActive_StartsAtZero(t *testing.T) {
	d := drain.New()
	if got := d.Active("svc"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAcquire_And_Release(t *testing.T) {
	d := drain.New()
	d.Acquire("svc")
	d.Acquire("svc")
	if got := d.Active("svc"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
	d.Release("svc")
	if got := d.Active("svc"); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestDrain_ImmediateWhenIdle(t *testing.T) {
	d := drain.New()
	ok := d.Drain("svc", 100*time.Millisecond)
	if !ok {
		t.Fatal("expected drain to succeed immediately when idle")
	}
}

func TestDrain_WaitsForRelease(t *testing.T) {
	d := drain.New()
	d.Acquire("svc")

	go func() {
		time.Sleep(30 * time.Millisecond)
		d.Release("svc")
	}()

	ok := d.Drain("svc", 200*time.Millisecond)
	if !ok {
		t.Fatal("expected drain to succeed after release")
	}
}

func TestDrain_Timeout(t *testing.T) {
	d := drain.New()
	d.Acquire("svc")

	ok := d.Drain("svc", 30*time.Millisecond)
	if ok {
		t.Fatal("expected drain to time out")
	}
}

func TestReset_ClearsCounter(t *testing.T) {
	d := drain.New()
	d.Acquire("svc")
	d.Acquire("svc")
	d.Reset("svc")
	if got := d.Active("svc"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestReset_UnblocksWaiters(t *testing.T) {
	d := drain.New()
	d.Acquire("svc")

	done := make(chan bool, 1)
	go func() {
		done <- d.Drain("svc", 500*time.Millisecond)
	}()

	time.Sleep(20 * time.Millisecond)
	d.Reset("svc")

	select {
	case result := <-done:
		_ = result // reset closes channel so drain returns true
	case <-time.After(200 * time.Millisecond):
		t.Fatal("drain did not unblock after reset")
	}
}
