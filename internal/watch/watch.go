package watch

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Event represents a file change event.
type Event struct {
	Path string
	Op   string
}

// Watcher monitors a set of file paths for changes and emits events.
type Watcher struct {
	mu       sync.Mutex
	paths    map[string]time.Time
	Events   chan Event
	Errors   chan error
	stopCh   chan struct{}
	interval time.Duration
}

// New creates a new Watcher with the given poll interval.
func New(interval time.Duration) *Watcher {
	return &Watcher{
		paths:    make(map[string]time.Time),
		Events:   make(chan Event, 16),
		Errors:   make(chan error, 4),
		stopCh:   make(chan struct{}),
		interval: interval,
	}
}

// Add registers a file or glob pattern to be watched.
func (w *Watcher) Add(pattern string) error {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		matches = []string{pattern}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, p := range matches {
		info, err := os.Stat(p)
		if err != nil {
			w.paths[p] = time.Time{}
			continue
		}
		w.paths[p] = info.ModTime()
	}
	return nil
}

// Start begins polling for file changes in a background goroutine.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-w.stopCh:
				return
			case <-ticker.C:
				w.poll()
			}
		}
	}()
}

// Stop halts the watcher.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

func (w *Watcher) poll() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for p, lastMod := range w.paths {
		info, err := os.Stat(p)
		if err != nil {
			if !lastMod.IsZero() {
				w.paths[p] = time.Time{}
				w.Events <- Event{Path: p, Op: "remove"}
			}
			continue
		}
		if info.ModTime().After(lastMod) {
			w.paths[p] = info.ModTime()
			op := "write"
			if lastMod.IsZero() {
				op = "create"
			}
			w.Events <- Event{Path: p, Op: op}
		}
	}
}
