// Package notify provides a simple notification broker for broadcasting
// service lifecycle events to registered subscribers (e.g. CLI output,
// webhooks, or other internal components).
package notify

import (
	"fmt"
	"sync"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Event represents a single notification emitted by the system.
type Event struct {
	Service   string
	Level     Level
	Message   string
	Timestamp time.Time
}

// Handler is a function that receives a notification event.
type Handler func(Event)

// Broker manages subscriptions and dispatches events to all registered handlers.
type Broker struct {
	mu       sync.RWMutex
	handlers map[string]Handler
	history  []Event
	maxHist  int
}

// New creates a new Broker with the given history buffer size.
// If maxHistory is <= 0, a default of 100 is used.
func New(maxHistory int) *Broker {
	if maxHistory <= 0 {
		maxHistory = 100
	}
	return &Broker{
		handlers: make(map[string]Handler),
		maxHist:  maxHistory,
	}
}

// Subscribe registers a named handler to receive all future events.
// If a handler with the same name already exists it is replaced.
func (b *Broker) Subscribe(name string, h Handler) error {
	if name == "" {
		return fmt.Errorf("notify: subscriber name must not be empty")
	}
	if h == nil {
		return fmt.Errorf("notify: handler must not be nil")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = h
	return nil
}

// Unsubscribe removes the handler with the given name. It is a no-op if the
// name is not registered.
func (b *Broker) Unsubscribe(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, name)
}

// Emit dispatches an event to all registered handlers and appends it to the
// internal history buffer. Handlers are called synchronously in an undefined
// order.
func (b *Broker) Emit(service string, level Level, message string) {
	ev := Event{
		Service:   service,
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
	}

	b.mu.Lock()
	// Append to rolling history.
	b.history = append(b.history, ev)
	if len(b.history) > b.maxHist {
		b.history = b.history[len(b.history)-b.maxHist:]
	}
	// Copy handlers so we can call them outside the lock.
	copy := make([]Handler, 0, len(b.handlers))
	for _, h := range b.handlers {
		copy = append(copy, h)
	}
	b.mu.Unlock()

	for _, h := range copy {
		h(ev)
	}
}

// History returns a copy of the stored event history, oldest first.
func (b *Broker) History() []Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Event, len(b.history))
	copy(out, b.history)
	return out
}

// Subscribers returns the names of all currently registered handlers.
func (b *Broker) Subscribers() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	names := make([]string, 0, len(b.handlers))
	for name := range b.handlers {
		names = append(names, name)
	}
	return names
}
