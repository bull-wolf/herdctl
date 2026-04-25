package ratelimit

import (
	"testing"
	"time"
)

func TestAllow_UnderLimit(t *testing.T) {
	l := New(5 * time.Second)
	for i := 0; i < 3; i++ {
		if !l.Allow("svc-a", 5) {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := New(5 * time.Second)
	for i := 0; i < 3; i++ {
		l.Allow("svc-b", 3)
	}
	if l.Allow("svc-b", 3) {
		t.Fatal("expected request to be denied after limit reached")
	}
}

func TestAllow_WindowReset(t *testing.T) {
	l := New(50 * time.Millisecond)
	for i := 0; i < 2; i++ {
		l.Allow("svc-c", 2)
	}
	if l.Allow("svc-c", 2) {
		t.Fatal("expected denial before window reset")
	}
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("svc-c", 2) {
		t.Fatal("expected allow after window reset")
	}
}

func TestGet_ReturnsEntry(t *testing.T) {
	l := New(5 * time.Second)
	l.Allow("svc-d", 10)
	l.Allow("svc-d", 10)

	e, err := l.Get("svc-d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Requests != 2 {
		t.Errorf("expected 2 requests, got %d", e.Requests)
	}
	if e.Limit != 10 {
		t.Errorf("expected limit 10, got %d", e.Limit)
	}
}

func TestGet_Missing(t *testing.T) {
	l := New(5 * time.Second)
	_, err := l.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing service")
	}
}

func TestReset_ClearsEntry(t *testing.T) {
	l := New(5 * time.Second)
	l.Allow("svc-e", 2)
	l.Allow("svc-e", 2)
	l.Reset("svc-e")

	if !l.Allow("svc-e", 2) {
		t.Fatal("expected allow after reset")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	l := New(5 * time.Second)
	l.Allow("svc-f", 5)
	l.Allow("svc-g", 5)

	all := l.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
