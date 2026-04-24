package audit

import (
	"fmt"
	"sync"
	"time"
)

// EventKind represents the type of audit event.
type EventKind string

const (
	EventStart   EventKind = "start"
	EventStop    EventKind = "stop"
	EventRestart EventKind = "restart"
	EventHook    EventKind = "hook"
	EventConfig  EventKind = "config_reload"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time
	Kind      EventKind
	Service   string
	Message   string
}

// Auditor records lifecycle and operational events.
type Auditor struct {
	mu      sync.RWMutex
	entries []Entry
	maxSize int
}

// New creates an Auditor that retains up to maxSize entries.
func New(maxSize int) *Auditor {
	if maxSize <= 0 {
		maxSize = 500
	}
	return &Auditor{maxSize: maxSize}
}

// Record appends a new audit entry.
func (a *Auditor) Record(kind EventKind, service, message string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	entry := Entry{
		Timestamp: time.Now(),
		Kind:      kind,
		Service:   service,
		Message:   message,
	}
	a.entries = append(a.entries, entry)
	if len(a.entries) > a.maxSize {
		a.entries = a.entries[len(a.entries)-a.maxSize:]
	}
}

// All returns a copy of all stored entries.
func (a *Auditor) All() []Entry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	copy := make([]Entry, len(a.entries))
	for i, e := range a.entries {
		copy[i] = e
	}
	return copy
}

// ForService returns entries matching the given service name.
func (a *Auditor) ForService(service string) []Entry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var result []Entry
	for _, e := range a.entries {
		if e.Service == service {
			result = append(result, e)
		}
	}
	return result
}

// Clear removes all stored entries.
func (a *Auditor) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = nil
}

// Summary returns a human-readable count per event kind.
func (a *Auditor) Summary() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	counts := make(map[EventKind]int)
	for _, e := range a.entries {
		counts[e.Kind]++
	}
	return fmt.Sprintf("start=%d stop=%d restart=%d hook=%d config_reload=%d",
		counts[EventStart], counts[EventStop], counts[EventRestart],
		counts[EventHook], counts[EventConfig])
}
