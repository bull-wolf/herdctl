package throttle

import (
	"testing"
	"time"
)

func TestAllow_UnderLimit(t *testing.T) {
	th := New(3, 10*time.Second)

	for i := 0; i < 3; i++ {
		ok, err := th.Allow("api")
		if !ok || err != nil {
			t.Fatalf("attempt %d: expected allowed, got ok=%v err=%v", i+1, ok, err)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	th := New(2, 10*time.Second)

	th.Allow("api") // 1
	th.Allow("api") // 2

	ok, err := th.Allow("api") // 3 — should be denied
	if ok {
		t.Fatal("expected deny after exceeding max attempts")
	}
	if err == nil {
		t.Fatal("expected error when throttled")
	}
}

func TestAllow_WindowReset(t *testing.T) {
	th := New(1, 50*time.Millisecond)

	ok, _ := th.Allow("worker")
	if !ok {
		t.Fatal("first attempt should be allowed")
	}

	ok, _ = th.Allow("worker")
	if ok {
		t.Fatal("second attempt within window should be denied")
	}

	time.Sleep(60 * time.Millisecond)

	ok, err := th.Allow("worker")
	if !ok || err != nil {
		t.Fatalf("after window expiry expected allow, got ok=%v err=%v", ok, err)
	}
}

func TestReset_ClearsCounter(t *testing.T) {
	th := New(1, 10*time.Second)

	th.Allow("db")
	th.Reset("db")

	ok, err := th.Allow("db")
	if !ok || err != nil {
		t.Fatalf("after reset expected allow, got ok=%v err=%v", ok, err)
	}
}

func TestGet_ReturnsEntry(t *testing.T) {
	th := New(5, 10*time.Second)
	th.Allow("svc")
	th.Allow("svc")

	e := th.Get("svc")
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if e.Attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", e.Attempts)
	}
}

func TestGet_Missing(t *testing.T) {
	th := New(5, 10*time.Second)
	if e := th.Get("unknown"); e != nil {
		t.Fatalf("expected nil for unknown service, got %+v", e)
	}
}

func TestAllow_IsolatedPerService(t *testing.T) {
	th := New(1, 10*time.Second)

	th.Allow("a")
	ok, err := th.Allow("b")
	if !ok || err != nil {
		t.Fatal("service b should be unaffected by service a's counter")
	}
}
