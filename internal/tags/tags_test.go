package tags

import (
	"testing"
)

func TestAdd_And_Get(t *testing.T) {
	s := New()
	if err := s.Add("api", "backend", "critical"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := s.Get("api")
	if len(got) != 2 || got[0] != "backend" || got[1] != "critical" {
		t.Errorf("expected [backend critical], got %v", got)
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	if got := s.Get("unknown"); got != nil {
		t.Errorf("expected nil for missing service, got %v", got)
	}
}

func TestAdd_DuplicateTag(t *testing.T) {
	s := New()
	s.Add("api", "backend")
	s.Add("api", "backend")
	if got := s.Get("api"); len(got) != 1 {
		t.Errorf("expected 1 tag, got %d", len(got))
	}
}

func TestAdd_EmptyService(t *testing.T) {
	s := New()
	if err := s.Add("", "backend"); err == nil {
		t.Error("expected error for empty service name")
	}
}

func TestAdd_EmptyTag(t *testing.T) {
	s := New()
	if err := s.Add("api", ""); err == nil {
		t.Error("expected error for empty tag value")
	}
}

func TestRemove_Tag(t *testing.T) {
	s := New()
	s.Add("api", "backend", "critical")
	s.Remove("api", "critical")
	if s.Has("api", "critical") {
		t.Error("expected tag 'critical' to be removed")
	}
	if !s.Has("api", "backend") {
		t.Error("expected tag 'backend' to remain")
	}
}

func TestRemove_NonExistentTag(t *testing.T) {
	s := New()
	// should not panic
	s.Remove("api", "ghost")
}

func TestHas_ReturnsFalseForMissing(t *testing.T) {
	s := New()
	if s.Has("api", "backend") {
		t.Error("expected false for service with no tags")
	}
}

func TestServicesWithTag(t *testing.T) {
	s := New()
	s.Add("api", "backend")
	s.Add("worker", "backend")
	s.Add("frontend", "ui")

	got := s.ServicesWithTag("backend")
	if len(got) != 2 || got[0] != "api" || got[1] != "worker" {
		t.Errorf("expected [api worker], got %v", got)
	}
}

func TestServicesWithTag_NoMatch(t *testing.T) {
	s := New()
	s.Add("api", "backend")
	if got := s.ServicesWithTag("ui"); len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}
