package watch

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAdd_TracksFile(t *testing.T) {
	w := New(100 * time.Millisecond)
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.txt")
	os.WriteFile(f, []byte("hello"), 0644)

	if err := w.Add(f); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	w.mu.Lock()
	_, ok := w.paths[f]
	w.mu.Unlock()
	if !ok {
		t.Error("expected file to be tracked")
	}
}

func TestPoll_DetectsWrite(t *testing.T) {
	w := New(50 * time.Millisecond)
	tmp := t.TempDir()
	f := filepath.Join(tmp, "change.txt")
	os.WriteFile(f, []byte("v1"), 0644)

	w.Add(f)

	time.Sleep(10 * time.Millisecond)
	os.WriteFile(f, []byte("v2"), 0644)
	w.poll()

	select {
	case ev := <-w.Events:
		if ev.Path != f {
			t.Errorf("expected path %s, got %s", f, ev.Path)
		}
		if ev.Op != "write" {
			t.Errorf("expected op 'write', got %s", ev.Op)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for write event")
	}
}

func TestPoll_DetectsRemove(t *testing.T) {
	w := New(50 * time.Millisecond)
	tmp := t.TempDir()
	f := filepath.Join(tmp, "gone.txt")
	os.WriteFile(f, []byte("data"), 0644)

	w.Add(f)
	os.Remove(f)
	w.poll()

	select {
	case ev := <-w.Events:
		if ev.Op != "remove" {
			t.Errorf("expected op 'remove', got %s", ev.Op)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for remove event")
	}
}

func TestStart_And_Stop(t *testing.T) {
	w := New(20 * time.Millisecond)
	tmp := t.TempDir()
	f := filepath.Join(tmp, "live.txt")
	os.WriteFile(f, []byte("init"), 0644)
	w.Add(f)
	w.Start()

	time.Sleep(10 * time.Millisecond)
	os.WriteFile(f, []byte("updated"), 0644)

	select {
	case ev := <-w.Events:
		if ev.Op != "write" {
			t.Errorf("unexpected op: %s", ev.Op)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for event from running watcher")
	}
	w.Stop()
}
