package ports

import (
	"net"
	"testing"
)

func TestAssign_And_Get(t *testing.T) {
	a := New()
	if err := a.Assign("api", 8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	port, ok := a.Get("api")
	if !ok || port != 8080 {
		t.Fatalf("expected port 8080, got %d (ok=%v)", port, ok)
	}
}

func TestAssign_Conflict(t *testing.T) {
	a := New()
	_ = a.Assign("api", 8080)
	err := a.Assign("worker", 8080)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}

func TestAssign_SameServiceIdempotent(t *testing.T) {
	a := New()
	_ = a.Assign("api", 8080)
	if err := a.Assign("api", 8080); err != nil {
		t.Fatalf("re-assigning same port to same service should be fine: %v", err)
	}
}

func TestGet_Missing(t *testing.T) {
	a := New()
	_, ok := a.Get("ghost")
	if ok {
		t.Fatal("expected ok=false for unknown service")
	}
}

func TestRelease(t *testing.T) {
	a := New()
	_ = a.Assign("api", 8080)
	a.Release("api")
	_, ok := a.Get("api")
	if ok {
		t.Fatal("expected port to be released")
	}
	// port should now be available for another service
	if err := a.Assign("other", 8080); err != nil {
		t.Fatalf("expected port to be free after release: %v", err)
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	a := New()
	_ = a.Assign("api", 8080)
	_ = a.Assign("db", 5432)
	all := a.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all["api"] != 8080 || all["db"] != 5432 {
		t.Fatalf("unexpected snapshot: %v", all)
	}
}

func TestIsFree_OccupiedPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skip("cannot bind a test listener")
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port
	if IsFree(port) {
		t.Fatalf("expected port %d to be occupied", port)
	}
}
