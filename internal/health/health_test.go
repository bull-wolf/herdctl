package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheck_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := New()
	res := c.Check("web", ts.URL)

	if res.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %s", res.Status, res.Err)
	}
	if res.Service != "web" {
		t.Errorf("expected service 'web', got %s", res.Service)
	}
}

func TestCheck_Unhealthy_BadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := New()
	res := c.Check("api", ts.URL)

	if res.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestCheck_Unhealthy_Unreachable(t *testing.T) {
	c := New()
	res := c.Check("db", "http://127.0.0.1:19999")

	if res.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy for unreachable host, got %s", res.Status)
	}
	if res.Err == "" {
		t.Error("expected non-empty error message")
	}
}

func TestGet_ReturnsStoredResult(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := New()
	c.Check("svc", ts.URL)

	r, ok := c.Get("svc")
	if !ok {
		t.Fatal("expected result to be stored")
	}
	if r.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s", r.Status)
	}
}

func TestAll_ReturnsAllResults(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := New()
	c.Check("a", ts.URL)
	c.Check("b", ts.URL)

	all := c.All()
	if len(all) != 2 {
		t.Errorf("expected 2 results, got %d", len(all))
	}
}
