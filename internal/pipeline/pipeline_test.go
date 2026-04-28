package pipeline

import (
	"testing"
)

func TestRegister_And_Get(t *testing.T) {
	p := New()
	err := p.Register("api", []string{"build", "test", "deploy"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stages, err := p.Get("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stages) != 3 {
		t.Fatalf("expected 3 stages, got %d", len(stages))
	}
	if stages[0].Name != "build" || stages[0].Status != "pending" {
		t.Errorf("unexpected first stage: %+v", stages[0])
	}
}

func TestRegister_InvalidArgs(t *testing.T) {
	p := New()
	if err := p.Register("", []string{"build"}); err == nil {
		t.Error("expected error for empty service")
	}
	if err := p.Register("api", nil); err == nil {
		t.Error("expected error for empty stages")
	}
}

func TestAdvance_ProgressesStages(t *testing.T) {
	p := New()
	_ = p.Register("svc", []string{"build", "deploy"})

	// first advance: pending -> running
	if err := p.Advance("svc", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stages, _ := p.Get("svc")
	if stages[0].Status != "running" {
		t.Errorf("expected running, got %s", stages[0].Status)
	}

	// second advance: running -> done, next pending -> running
	_ = p.Advance("svc", false)
	stages, _ = p.Get("svc")
	if stages[0].Status != "done" {
		t.Errorf("expected done, got %s", stages[0].Status)
	}
	if stages[1].Status != "running" {
		t.Errorf("expected running, got %s", stages[1].Status)
	}
}

func TestAdvance_FailedStage(t *testing.T) {
	p := New()
	_ = p.Register("svc", []string{"build"})
	_ = p.Advance("svc", false) // pending -> running
	_ = p.Advance("svc", true)  // running -> failed
	stages, _ := p.Get("svc")
	if stages[0].Status != "failed" {
		t.Errorf("expected failed, got %s", stages[0].Status)
	}
}

func TestAdvance_UnknownService(t *testing.T) {
	p := New()
	if err := p.Advance("ghost", false); err == nil {
		t.Error("expected error for unknown service")
	}
}

func TestReset_ClearsStages(t *testing.T) {
	p := New()
	_ = p.Register("svc", []string{"build", "deploy"})
	_ = p.Advance("svc", false)
	_ = p.Advance("svc", false)
	_ = p.Reset("svc")
	stages, _ := p.Get("svc")
	for _, s := range stages {
		if s.Status != "pending" {
			t.Errorf("expected pending after reset, got %s", s.Status)
		}
	}
}

func TestGet_Missing(t *testing.T) {
	p := New()
	_, err := p.Get("nope")
	if err == nil {
		t.Error("expected error for missing service")
	}
}
