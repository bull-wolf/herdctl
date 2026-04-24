package hooks_test

import (
	"errors"
	"testing"

	"github.com/user/herdctl/internal/hooks"
)

func TestRegister_And_Fire(t *testing.T) {
	r := hooks.New()
	called := false
	r.Register("api", hooks.EventAfterStart, func(svc string, ev hooks.Event) error {
		called = true
		return nil
	})

	if err := r.Fire("api", hooks.EventAfterStart); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected hook to be called")
	}
}

func TestFire_NoHooks(t *testing.T) {
	r := hooks.New()
	if err := r.Fire("api", hooks.EventBeforeStart); err != nil {
		t.Fatalf("expected no error for missing hook, got: %v", err)
	}
}

func TestFire_MultipleHooks(t *testing.T) {
	r := hooks.New()
	count := 0
	for i := 0; i < 3; i++ {
		r.Register("worker", hooks.EventBeforeStop, func(svc string, ev hooks.Event) error {
			count++
			return nil
		})
	}
	r.Fire("worker", hooks.EventBeforeStop)
	if count != 3 {
		t.Fatalf("expected 3 hooks called, got %d", count)
	}
}

func TestFire_CollectsErrors(t *testing.T) {
	r := hooks.New()
	r.Register("db", hooks.EventAfterStop, func(svc string, ev hooks.Event) error {
		return errors.New("hook failed")
	})
	r.Register("db", hooks.EventAfterStop, func(svc string, ev hooks.Event) error {
		return errors.New("another failure")
	})

	err := r.Fire("db", hooks.EventAfterStop)
	if err == nil {
		t.Fatal("expected error from failing hooks")
	}
}

func TestList_ReturnsRegisteredEvents(t *testing.T) {
	r := hooks.New()
	r.Register("api", hooks.EventBeforeStart, func(svc string, ev hooks.Event) error { return nil })
	r.Register("api", hooks.EventAfterStop, func(svc string, ev hooks.Event) error { return nil })

	events := r.List("api")
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestClear_RemovesHooks(t *testing.T) {
	r := hooks.New()
	called := false
	r.Register("api", hooks.EventAfterStart, func(svc string, ev hooks.Event) error {
		called = true
		return nil
	})
	r.Clear("api")
	r.Fire("api", hooks.EventAfterStart)
	if called {
		t.Fatal("hook should not be called after Clear")
	}
}
