package labels

import "testing"

func TestSet_And_Get(t *testing.T) {
	s := New()
	s.Set("api", "env", "production")
	v, ok := s.Get("api", "env")
	if !ok {
		t.Fatal("expected label to exist")
	}
	if v != "production" {
		t.Errorf("expected 'production', got %q", v)
	}
}

func TestGet_Missing(t *testing.T) {
	s := New()
	_, ok := s.Get("api", "env")
	if ok {
		t.Fatal("expected label to be missing")
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := New()
	s.Set("api", "env", "staging")
	s.Set("api", "env", "production")
	v, _ := s.Get("api", "env")
	if v != "production" {
		t.Errorf("expected overwritten value 'production', got %q", v)
	}
}

func TestAll_ReturnsAllLabels(t *testing.T) {
	s := New()
	s.Set("api", "env", "prod")
	s.Set("api", "team", "platform")
	m := s.All("api")
	if len(m) != 2 {
		t.Errorf("expected 2 labels, got %d", len(m))
	}
}

func TestAll_MissingService(t *testing.T) {
	s := New()
	if m := s.All("ghost"); m != nil {
		t.Errorf("expected nil for unknown service, got %v", m)
	}
}

func TestDelete_RemovesLabel(t *testing.T) {
	s := New()
	s.Set("api", "env", "prod")
	s.Delete("api", "env")
	_, ok := s.Get("api", "env")
	if ok {
		t.Fatal("expected label to be deleted")
	}
}

func TestClear_RemovesAllLabels(t *testing.T) {
	s := New()
	s.Set("api", "env", "prod")
	s.Set("api", "team", "platform")
	s.Clear("api")
	if m := s.All("api"); m != nil {
		t.Errorf("expected nil after clear, got %v", m)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := New()
	s.Set("api", "env", "prod")
	m := s.All("api")
	m["env"] = "mutated"
	v, _ := s.Get("api", "env")
	if v != "prod" {
		t.Errorf("store mutated via returned map: got %q", v)
	}
}
