package env

import (
	"os"
	"strings"
	"testing"
)

func TestResolve_GlobalOverridesOS(t *testing.T) {
	os.Setenv("HERD_TEST_KEY", "os_value")
	t.Cleanup(func() { os.Unsetenv("HERD_TEST_KEY") })

	r := New(map[string]string{"HERD_TEST_KEY": "global_value"})
	env := r.Resolve("svc")

	for _, e := range env {
		if strings.HasPrefix(e, "HERD_TEST_KEY=") {
			if e != "HERD_TEST_KEY=global_value" {
				t.Fatalf("expected global_value, got %s", e)
			}
			return
		}
	}
	t.Fatal("HERD_TEST_KEY not found in resolved env")
}

func TestResolve_ServiceOverridesGlobal(t *testing.T) {
	r := New(map[string]string{"PORT": "8080"})
	r.SetService("api", map[string]string{"PORT": "9090"})

	env := r.Resolve("api")
	for _, e := range env {
		if strings.HasPrefix(e, "PORT=") {
			if e != "PORT=9090" {
				t.Fatalf("expected PORT=9090, got %s", e)
			}
			return
		}
	}
	t.Fatal("PORT not found in resolved env")
}

func TestResolve_OtherServiceUnaffected(t *testing.T) {
	r := New(map[string]string{"PORT": "8080"})
	r.SetService("api", map[string]string{"PORT": "9090"})

	env := r.Resolve("worker")
	for _, e := range env {
		if e == "PORT=9090" {
			t.Fatal("worker should not have api override")
		}
	}
}

func TestGet_ServiceKey(t *testing.T) {
	r := New(map[string]string{"X": "global"})
	r.SetService("svc", map[string]string{"X": "svc_val"})

	v, ok := r.Get("svc", "X")
	if !ok || v != "svc_val" {
		t.Fatalf("expected svc_val, got %q ok=%v", v, ok)
	}
}

func TestGet_FallsBackToGlobal(t *testing.T) {
	r := New(map[string]string{"Y": "global_y"})
	v, ok := r.Get("svc", "Y")
	if !ok || v != "global_y" {
		t.Fatalf("expected global_y, got %q ok=%v", v, ok)
	}
}

func TestGet_Missing(t *testing.T) {
	r := New(nil)
	_, ok := r.Get("svc", "DEFINITELY_NOT_SET_XYZ")
	if ok {
		t.Fatal("expected key to be missing")
	}
}
