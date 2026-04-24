package audit

import (
	"strings"
	"testing"
)

func TestRecord_And_All(t *testing.T) {
	a := New(100)
	a.Record(EventStart, "api", "service started")
	a.Record(EventStop, "api", "service stopped")

	entries := a.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Kind != EventStart {
		t.Errorf("expected EventStart, got %s", entries[0].Kind)
	}
	if entries[1].Service != "api" {
		t.Errorf("expected service 'api', got %s", entries[1].Service)
	}
}

func TestForService_FiltersCorrectly(t *testing.T) {
	a := New(100)
	a.Record(EventStart, "api", "started")
	a.Record(EventStart, "worker", "started")
	a.Record(EventStop, "api", "stopped")

	result := a.ForService("api")
	if len(result) != 2 {
		t.Fatalf("expected 2 entries for 'api', got %d", len(result))
	}
	for _, e := range result {
		if e.Service != "api" {
			t.Errorf("unexpected service %s in filtered results", e.Service)
		}
	}
}

func TestRecord_RollingWindow(t *testing.T) {
	a := New(3)
	a.Record(EventStart, "svc", "msg1")
	a.Record(EventStop, "svc", "msg2")
	a.Record(EventRestart, "svc", "msg3")
	a.Record(EventHook, "svc", "msg4")

	entries := a.All()
	if len(entries) != 3 {
		t.Fatalf("expected rolling window of 3, got %d", len(entries))
	}
	if entries[0].Message != "msg2" {
		t.Errorf("expected oldest retained to be msg2, got %s", entries[0].Message)
	}
}

func TestClear_RemovesAll(t *testing.T) {
	a := New(100)
	a.Record(EventStart, "api", "started")
	a.Clear()
	if len(a.All()) != 0 {
		t.Error("expected empty entries after Clear")
	}
}

func TestSummary_CountsKinds(t *testing.T) {
	a := New(100)
	a.Record(EventStart, "api", "")
	a.Record(EventStart, "worker", "")
	a.Record(EventStop, "api", "")
	a.Record(EventRestart, "api", "")

	summary := a.Summary()
	if !strings.Contains(summary, "start=2") {
		t.Errorf("expected start=2 in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "stop=1") {
		t.Errorf("expected stop=1 in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "restart=1") {
		t.Errorf("expected restart=1 in summary, got: %s", summary)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	a := New(100)
	a.Record(EventStart, "api", "started")
	entries := a.All()
	entries[0].Service = "mutated"

	original := a.All()
	if original[0].Service == "mutated" {
		t.Error("All() should return a copy, not a reference")
	}
}
