package signal

import (
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestRegister_And_Notify(t *testing.T) {
	h := New()
	h.Start()
	defer h.Stop()

	var mu sync.Mutex
	received := []os.Signal{}

	h.Register(func(sig os.Signal) {
		mu.Lock()
		received = append(received, sig)
		mu.Unlock()
	})

	h.Notify(syscall.SIGINT)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 1 {
		t.Fatalf("expected 1 signal, got %d", len(received))
	}
	if received[0] != syscall.SIGINT {
		t.Errorf("expected SIGINT, got %v", received[0])
	}
}

func TestRegister_MultipleListeners(t *testing.T) {
	h := New()
	h.Start()
	defer h.Stop()

	var wg sync.WaitGroup
	wg.Add(2)

	h.Register(func(sig os.Signal) { wg.Done() })
	h.Register(func(sig os.Signal) { wg.Done() })

	h.Notify(syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for both listeners to fire")
	}
}

func TestStop_Idempotent(t *testing.T) {
	h := New()
	h.Start()
	h.Stop()
	h.Stop() // should not panic
}

func TestNotify_NoListeners(t *testing.T) {
	h := New()
	h.Start()
	defer h.Stop()

	// Should not block or panic with no listeners registered
	h.Notify(syscall.SIGINT)
	time.Sleep(30 * time.Millisecond)
}
