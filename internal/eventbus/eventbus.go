package eventbus

import (
	"fmt"
	"sync"
)

// Handler is a function that receives an event payload.
type Handler func(service, event string, payload map[string]string)

// Entry represents a subscription.
type Entry struct {
	Service string
	Event   string
	Handler Handler
}

// Bus is a simple in-process pub/sub event bus keyed by service and event.
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler // key: "service:event"
}

// New creates a new Bus.
func New() *Bus {
	return &Bus{
		subscribers: make(map[string][]Handler),
	}
}

func key(service, event string) string {
	return service + ":" + event
}

// Subscribe registers a handler for the given service and event.
func (b *Bus) Subscribe(service, event string, h Handler) error {
	if service == "" || event == "" {
		return fmt.Errorf("eventbus: service and event must not be empty")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	k := key(service, event)
	b.subscribers[k] = append(b.subscribers[k], h)
	return nil
}

// Publish delivers payload to all handlers subscribed to service+event.
func (b *Bus) Publish(service, event string, payload map[string]string) error {
	if service == "" || event == "" {
		return fmt.Errorf("eventbus: service and event must not be empty")
	}
	b.mu.RLock()
	handlers := make([]Handler, len(b.subscribers[key(service, event)]))
	copy(handlers, b.subscribers[key(service, event)])
	b.mu.RUnlock()
	for _, h := range handlers {
		h(service, event, payload)
	}
	return nil
}

// Unsubscribe removes all handlers for the given service and event.
func (b *Bus) Unsubscribe(service, event string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subscribers, key(service, event))
}

// List returns all active subscriptions as a slice of Entry (handler omitted).
func (b *Bus) List() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	var out []Entry
	for k, hs := range b.subscribers {
		var svc, ev string
		fmt.Sscanf(k, "%s", &svc) // fallback; we parse manually below
		_ = hs
		// parse key manually
		for i, c := range k {
			if c == ':' {
				svc = k[:i]
				ev = k[i+1:]
				break
			}
		}
		for range hs {
			out = append(out, Entry{Service: svc, Event: ev})
		}
	}
	return out
}
