package eventbus

import (
	"sync"
	"testing"
)

func TestSubscribe_And_Publish(t *testing.T) {
	b := New()
	var got []string
	_ = b.Subscribe("api", "started", func(svc, ev string, p map[string]string) {
		got = append(got, svc+":"+ev)
	})
	_ = b.Publish("api", "started", nil)
	if len(got) != 1 || got[0] != "api:started" {
		t.Fatalf("expected [api:started], got %v", got)
	}
}

func TestPublish_NoSubscribers(t *testing.T) {
	b := New()
	if err := b.Publish("api", "stopped", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSubscribe_InvalidArgs(t *testing.T) {
	b := New()
	if err := b.Subscribe("", "started", func(_, _ string, _ map[string]string) {}); err == nil {
		t.Fatal("expected error for empty service")
	}
	if err := b.Subscribe("api", "", func(_, _ string, _ map[string]string) {}); err == nil {
		t.Fatal("expected error for empty event")
	}
}

func TestPublish_InvalidArgs(t *testing.T) {
	b := New()
	if err := b.Publish("", "started", nil); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestUnsubscribe_RemovesHandlers(t *testing.T) {
	b := New()
	called := false
	_ = b.Subscribe("api", "started", func(_, _ string, _ map[string]string) { called = true })
	b.Unsubscribe("api", "started")
	_ = b.Publish("api", "started", nil)
	if called {
		t.Fatal("handler should not have been called after unsubscribe")
	}
}

func TestPublish_MultipleHandlers(t *testing.T) {
	b := New()
	var mu sync.Mutex
	count := 0
	inc := func(_, _ string, _ map[string]string) {
		mu.Lock()
		count++
		mu.Unlock()
	}
	_ = b.Subscribe("db", "ready", inc)
	_ = b.Subscribe("db", "ready", inc)
	_ = b.Publish("db", "ready", nil)
	if count != 2 {
		t.Fatalf("expected 2 handler calls, got %d", count)
	}
}

func TestPublish_PayloadDelivered(t *testing.T) {
	b := New()
	var received map[string]string
	_ = b.Subscribe("svc", "info", func(_, _ string, p map[string]string) { received = p })
	_ = b.Publish("svc", "info", map[string]string{"key": "val"})
	if received["key"] != "val" {
		t.Fatalf("expected val, got %q", received["key"])
	}
}

func TestList_ReturnsSubscriptions(t *testing.T) {
	b := New()
	_ = b.Subscribe("api", "started", func(_, _ string, _ map[string]string) {})
	_ = b.Subscribe("db", "stopped", func(_, _ string, _ map[string]string) {})
	list := b.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}
}
