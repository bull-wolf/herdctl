package watch

import (
	"log"
	"time"

	"github.com/user/herdctl/internal/config"
)

// ServiceWatcher maps service names to their file watchers and triggers
// a restart callback when a watched file changes.
type ServiceWatcher struct {
	watchers map[string]*Watcher
	OnChange func(serviceName string, event Event)
}

// NewServiceWatcher creates a ServiceWatcher from the loaded config.
// Services with a non-empty "watch" list will have watchers registered.
func NewServiceWatcher(cfg *config.Config, onChange func(string, Event)) *ServiceWatcher {
	sw := &ServiceWatcher{
		watchers: make(map[string]*Watcher),
		OnChange: onChange,
	}
	for _, svc := range cfg.Services {
		if len(svc.Watch) == 0 {
			continue
		}
		w := New(500 * time.Millisecond)
		for _, pattern := range svc.Watch {
			if err := w.Add(pattern); err != nil {
				log.Printf("watch: service %s: failed to add pattern %q: %v", svc.Name, pattern, err)
			}
		}
		sw.watchers[svc.Name] = w
	}
	return sw
}

// Start begins watching all registered services.
func (sw *ServiceWatcher) Start() {
	for name, w := range sw.watchers {
		name, w := name, w
		w.Start()
		go func() {
			for ev := range w.Events {
				if sw.OnChange != nil {
					sw.OnChange(name, ev)
				}
			}
		}()
	}
}

// Stop halts all active watchers.
func (sw *ServiceWatcher) Stop() {
	for _, w := range sw.watchers {
		w.Stop()
	}
}

// Watching returns the list of service names that have active watchers.
func (sw *ServiceWatcher) Watching() []string {
	names := make([]string, 0, len(sw.watchers))
	for name := range sw.watchers {
		names = append(names, name)
	}
	return names
}
