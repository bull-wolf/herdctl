package snapshot_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/herdctl/internal/snapshot"
)

func TestCapture_And_Get(t *testing.T) {
	s := snapshot.New()
	s.Capture(snapshot.Entry{
		Service: "api",
		Status:  "running",
		Health:  "healthy",
	})
	e, ok := s.Get("api")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if e.Status != "running" {
		t.Errorf("expected status running, got %s", e.Status)
	}
	if e.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestGet_Missing(t *testing.T) {
	s := snapshot.New()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("expected no snapshot for unknown service")
	}
}

func TestCapture_Overwrites(t *testing.T) {
	s := snapshot.New()
	s.Capture(snapshot.Entry{Service: "db", Status: "starting"})
	s.Capture(snapshot.Entry{Service: "db", Status: "running"})
	e, _ := s.Get("db")
	if e.Status != "running" {
		t.Errorf("expected overwritten status running, got %s", e.Status)
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	s := snapshot.New()
	s.Capture(snapshot.Entry{Service: "api", Status: "running"})
	s.Capture(snapshot.Entry{Service: "db", Status: "running"})
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(all))
	}
}

func TestSaveToFile_WritesJSON(t *testing.T) {
	s := snapshot.New()
	s.Capture(snapshot.Entry{Service: "api", Status: "running", Health: "healthy"})

	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	if err := s.SaveToFile(path); err != nil {
		t.Fatalf("SaveToFile error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}
	var entries []snapshot.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Service != "api" {
		t.Errorf("unexpected service name: %s", entries[0].Service)
	}
}
