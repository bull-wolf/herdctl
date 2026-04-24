package signal

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Handler listens for OS signals and notifies registered listeners.
type Handler struct {
	mu        sync.RWMutex
	listeners []func(os.Signal)
	ch        chan os.Signal
	stopCh    chan struct{}
	stopped   bool
}

// New creates a new Handler that listens for SIGINT and SIGTERM.
func New() *Handler {
	return &Handler{
		ch:     make(chan os.Signal, 1),
		stopCh: make(chan struct{}),
	}
}

// Register adds a listener function that will be called when a signal is received.
func (h *Handler) Register(fn func(os.Signal)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.listeners = append(h.listeners, fn)
}

// Start begins listening for OS signals in a background goroutine.
func (h *Handler) Start() {
	signal.Notify(h.ch, syscall.SIGINT, syscall.SIGTERM)
	go h.loop()
}

// Stop cancels signal listening and unregisters the channel.
func (h *Handler) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.stopped {
		return
	}
	h.stopped = true
	signal.Stop(h.ch)
	close(h.stopCh)
}

// Notify sends a signal directly to the handler (useful for testing).
func (h *Handler) Notify(sig os.Signal) {
	h.ch <- sig
}

func (h *Handler) loop() {
	for {
		select {
		case sig := <-h.ch:
			h.mu.RLock()
			listeners := make([]func(os.Signal), len(h.listeners))
			copy(listeners, h.listeners)
			h.mu.RUnlock()
			for _, fn := range listeners {
				fn(sig)
			}
		case <-h.stopCh:
			return
		}
	}
}
