package ports

import (
	"fmt"
	"net"
	"sync"
)

// Allocation tracks which ports have been assigned to services.
type Allocation struct {
	mu      sync.RWMutex
	assigned map[string]int // service name -> port
	reverse  map[int]string // port -> service name
}

// New returns an initialised Allocation registry.
func New() *Allocation {
	return &Allocation{
		assigned: make(map[string]int),
		reverse:  make(map[int]string),
	}
}

// Assign records a static port for a named service.
// Returns an error if the port is already taken by another service.
func (a *Allocation) Assign(service string, port int) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if owner, conflict := a.reverse[port]; conflict && owner != service {
		return fmt.Errorf("port %d already assigned to service %q", port, owner)
	}
	a.assigned[service] = port
	a.reverse[port] = service
	return nil
}

// Get returns the port assigned to a service and whether one exists.
func (a *Allocation) Get(service string) (int, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	port, ok := a.assigned[service]
	return port, ok
}

// Release removes the port assignment for a service.
func (a *Allocation) Release(service string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if port, ok := a.assigned[service]; ok {
		delete(a.reverse, port)
		delete(a.assigned, service)
	}
}

// All returns a snapshot of every service-to-port mapping.
func (a *Allocation) All() map[string]int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make(map[string]int, len(a.assigned))
	for k, v := range a.assigned {
		out[k] = v
	}
	return out
}

// IsFree returns true when nothing is currently listening on the given TCP port.
func IsFree(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}
