package metrics

import (
	"testing"
	"time"
)

func sampleEntry(service string, cpu, mem float64) Entry {
	return Entry{
		Service:   service,
		Timestamp: time.Now(),
		CPU:       cpu,
		MemoryMB:  mem,
		UptimeSec: 10,
	}
}

func TestRecord_And_Latest(t *testing.T) {
	c := New(10)
	c.Record(sampleEntry("api", 12.5, 64.0))

	e, ok := c.Latest("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.CPU != 12.5 {
		t.Errorf("expected CPU 12.5, got %v", e.CPU)
	}
}

func TestLatest_Missing(t *testing.T) {
	c := New(10)
	_, ok := c.Latest("nonexistent")
	if ok {
		t.Fatal("expected no entry for unknown service")
	}
}

func TestRecord_RollingWindow(t *testing.T) {
	c := New(3)
	for i := 0; i < 5; i++ {
		c.Record(sampleEntry("worker", float64(i), 0))
	}

	h := c.History("worker")
	if len(h) != 3 {
		t.Errorf("expected 3 samples (rolling), got %d", len(h))
	}
	if h[0].CPU != 2.0 {
		t.Errorf("expected oldest retained CPU=2, got %v", h[0].CPU)
	}
}

func TestHistory_ReturnsCopy(t *testing.T) {
	c := New(10)
	c.Record(sampleEntry("db", 5.0, 128.0))

	h := c.History("db")
	h[0].CPU = 999.0

	e, _ := c.Latest("db")
	if e.CPU == 999.0 {
		t.Error("History should return a copy, not a reference")
	}
}

func TestServices_ReturnsNames(t *testing.T) {
	c := New(10)
	c.Record(sampleEntry("api", 1, 1))
	c.Record(sampleEntry("worker", 2, 2))

	svcs := c.Services()
	if len(svcs) != 2 {
		t.Errorf("expected 2 services, got %d", len(svcs))
	}
}

func TestNew_DefaultsMaxLen(t *testing.T) {
	c := New(0)
	if c.maxLen != 60 {
		t.Errorf("expected default maxLen 60, got %d", c.maxLen)
	}
}
