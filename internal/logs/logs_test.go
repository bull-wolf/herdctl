package logs

import (
	"bytes"
	"strings"
	"testing"
)

func TestWrite_StoresEntry(t *testing.T) {
	c := New()
	c.Write("api", "started")
	c.Write("api", "listening on :8080")

	entries := c.Tail(0)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Service != "api" {
		t.Errorf("unexpected service: %s", entries[0].Service)
	}
}

func TestWrite_ForwardsToRegisteredWriter(t *testing.T) {
	c := New()
	var buf bytes.Buffer
	c.Register("db", &buf)
	c.Write("db", "postgres ready")

	if !strings.Contains(buf.String(), "postgres ready") {
		t.Errorf("expected log line in writer, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "db") {
		t.Errorf("expected service name in output, got: %s", buf.String())
	}
}

func TestTail_LimitsResults(t *testing.T) {
	c := New()
	for i := 0; i < 10; i++ {
		c.Write("worker", "tick")
	}

	entries := c.Tail(3)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestTailService_FiltersCorrectly(t *testing.T) {
	c := New()
	c.Write("api", "msg1")
	c.Write("db", "msg2")
	c.Write("api", "msg3")

	entries := c.TailService("api", 0)
	if len(entries) != 2 {
		t.Fatalf("expected 2 api entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Service != "api" {
			t.Errorf("unexpected service in filtered results: %s", e.Service)
		}
	}
}

func TestTailService_Limit(t *testing.T) {
	c := New()
	for i := 0; i < 5; i++ {
		c.Write("api", "line")
	}
	c.Write("db", "other")

	entries := c.TailService("api", 2)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}
