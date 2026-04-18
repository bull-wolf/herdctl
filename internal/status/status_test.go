package status_test

import (
	"testing"

	"herdctl/internal/status"
)

func TestSet_And_Get(t *testing.T) {
	tr := status.New()
	tr.Set("api", status.StateRunning, 1234)

	e, ok := tr.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.State != status.StateRunning {
		t.Errorf("expected running, got %s", e.State)
	}
	if e.PID != 1234 {
		t.Errorf("expected pid 1234, got %d", e.PID)
	}
}

func TestGet_Missing(t *testing.T) {
	tr := status.New()
	_, ok := tr.Get("nonexistent")
	if ok {
		t.Fatal("expected no entry for unknown service")
	}
}

func TestSet_Overwrites(t *testing.T) {
	tr := status.New()
	tr.Set("worker", status.StateStarting, 0)
	tr.Set("worker", status.StateRunning, 5678)

	e, _ := tr.Get("worker")
	if e.State != status.StateRunning {
		t.Errorf("expected running after overwrite, got %s", e.State)
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	tr := status.New()
	tr.Set("api", status.StateRunning, 1)
	tr.Set("db", status.StateStopped, 0)
	tr.Set("cache", status.StateFailed, 0)

	all := tr.All()
	if len(all) != 3 {
		t.Errorf("expected 3 entries, got %d", len(all))
	}
}

func TestAll_Empty(t *testing.T) {
	tr := status.New()
	if len(tr.All()) != 0 {
		t.Error("expected empty result for new tracker")
	}
}
