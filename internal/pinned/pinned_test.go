package pinned

import (
	"testing"
)

func TestPin_And_Get(t *testing.T) {
	s := New()
	if err := s.Pin("api", "v1.2.3", "alice", "stability"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Version != "v1.2.3" {
		t.Errorf("expected version v1.2.3, got %s", e.Version)
	}
	if e.PinnedBy != "alice" {
		t.Errorf("expected pinnedBy alice, got %s", e.PinnedBy)
	}
}

func TestPin_InvalidArgs(t *testing.T) {
	s := New()
	if err := s.Pin("", "v1.0.0", "", ""); err == nil {
		t.Error("expected error for empty service")
	}
	if err := s.Pin("api", "", "", ""); err == nil {
		t.Error("expected error for empty version")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get("unknown")
	if ok {
		t.Error("expected no entry for unknown service")
	}
}

func TestIsPinned(t *testing.T) {
	s := New()
	if s.IsPinned("api") {
		t.Error("expected service to not be pinned initially")
	}
	_ = s.Pin("api", "v2.0.0", "bob", "")
	if !s.IsPinned("api") {
		t.Error("expected service to be pinned")
	}
}

func TestUnpin_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Pin("api", "v1.0.0", "alice", "")
	if err := s.Unpin("api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsPinned("api") {
		t.Error("expected service to be unpinned")
	}
}

func TestUnpin_NotPinned(t *testing.T) {
	s := New()
	if err := s.Unpin("api"); err == nil {
		t.Error("expected error when unpinning non-pinned service")
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	s := New()
	_ = s.Pin("api", "v1.0.0", "alice", "")
	_ = s.Pin("worker", "v2.1.0", "bob", "perf")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestAll_Empty(t *testing.T) {
	s := New()
	if len(s.All()) != 0 {
		t.Error("expected empty result")
	}
}
