package traceid

import (
	"testing"
)

func TestGenerate_And_Get(t *testing.T) {
	s := New()
	id, err := s.Generate("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty trace ID")
	}
	e, ok := s.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.TraceID != id {
		t.Errorf("expected %q, got %q", id, e.TraceID)
	}
	if e.Service != "api" {
		t.Errorf("expected service 'api', got %q", e.Service)
	}
}

func TestGenerate_EmptyService(t *testing.T) {
	s := New()
	_, err := s.Generate("")
	if err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestGenerate_Unique(t *testing.T) {
	s := New()
	id1, _ := s.Generate("svc")
	id2, _ := s.Generate("svc")
	if id1 == id2 {
		t.Error("expected unique trace IDs on regeneration")
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get("unknown")
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestClear_RemovesEntry(t *testing.T) {
	s := New()
	s.Generate("worker")
	s.Clear("worker")
	_, ok := s.Get("worker")
	if ok {
		t.Fatal("expected entry to be cleared")
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	s := New()
	s.Generate("a")
	s.Generate("b")
	s.Generate("c")
	all := s.All()
	if len(all) != 3 {
		t.Errorf("expected 3 entries, got %d", len(all))
	}
}

func TestAll_Empty(t *testing.T) {
	s := New()
	all := s.All()
	if len(all) != 0 {
		t.Errorf("expected 0 entries on new store, got %d", len(all))
	}
}

func TestClear_NonExistent(t *testing.T) {
	// Clearing a service that was never added should not panic or error.
	s := New()
	s.Clear("nonexistent")
	_, ok := s.Get("nonexistent")
	if ok {
		t.Fatal("expected entry to remain absent after clearing nonexistent service")
	}
}
