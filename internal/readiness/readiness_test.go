package readiness_test

import (
	"testing"

	"github.com/user/herdctl/internal/readiness"
)

func TestSet_And_Get(t *testing.T) {
	tr := readiness.New()
	err := tr.Set("api", readiness.StateReady, "all checks passed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := tr.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.State != readiness.StateReady {
		t.Errorf("expected StateReady, got %v", e.State)
	}
	if e.Reason != "all checks passed" {
		t.Errorf("unexpected reason: %s", e.Reason)
	}
}

func TestSet_InvalidArgs(t *testing.T) {
	tr := readiness.New()
	err := tr.Set("", readiness.StateReady, "")
	if err == nil {
		t.Fatal("expected error for empty service name")
	}
}

func TestGet_Missing(t *testing.T) {
	tr := readiness.New()
	_, ok := tr.Get("missing")
	if ok {
		t.Fatal("expected no entry for unknown service")
	}
}

func TestIsReady_True(t *testing.T) {
	tr := readiness.New()
	_ = tr.Set("worker", readiness.StateReady, "ok")
	if !tr.IsReady("worker") {
		t.Error("expected worker to be ready")
	}
}

func TestIsReady_False(t *testing.T) {
	tr := readiness.New()
	_ = tr.Set("worker", readiness.StateNotReady, "starting")
	if tr.IsReady("worker") {
		t.Error("expected worker to not be ready")
	}
}

func TestIsReady_Missing(t *testing.T) {
	tr := readiness.New()
	if tr.IsReady("ghost") {
		t.Error("expected missing service to not be ready")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	tr := readiness.New()
	_ = tr.Set("svc-a", readiness.StateReady, "")
	_ = tr.Set("svc-b", readiness.StateNotReady, "waiting")
	all := tr.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestClear_RemovesEntry(t *testing.T) {
	tr := readiness.New()
	_ = tr.Set("db", readiness.StateReady, "")
	tr.Clear("db")
	_, ok := tr.Get("db")
	if ok {
		t.Error("expected entry to be cleared")
	}
}

func TestState_String(t *testing.T) {
	if readiness.StateReady.String() != "ready" {
		t.Errorf("unexpected string for StateReady")
	}
	if readiness.StateNotReady.String() != "not_ready" {
		t.Errorf("unexpected string for StateNotReady")
	}
	if readiness.StateUnknown.String() != "unknown" {
		t.Errorf("unexpected string for StateUnknown")
	}
}
