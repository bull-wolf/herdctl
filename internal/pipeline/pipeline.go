package pipeline

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Stage represents a named step in a pipeline.
type Stage struct {
	Name     string
	Status   string // pending, running, done, failed
	Started  time.Time
	Finished time.Time
	Err      error
}

// Pipeline tracks ordered stages for a service.
type Pipeline struct {
	mu       sync.RWMutex
	pipelines map[string][]*Stage
}

func New() *Pipeline {
	return &Pipeline{
		pipelines: make(map[string][]*Stage),
	}
}

// Register initialises a pipeline with the given stage names for a service.
func (p *Pipeline) Register(service string, stages []string) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	if len(stages) == 0 {
		return errors.New("at least one stage is required")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	list := make([]*Stage, len(stages))
	for i, name := range stages {
		list[i] = &Stage{Name: name, Status: "pending"}
	}
	p.pipelines[service] = list
	return nil
}

// Advance marks the next pending stage as running or done.
func (p *Pipeline) Advance(service string, failed bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	stages, ok := p.pipelines[service]
	if !ok {
		return fmt.Errorf("no pipeline registered for service %q", service)
	}
	for _, s := range stages {
		if s.Status == "running" {
			s.Finished = time.Now()
			if failed {
				s.Status = "failed"
			} else {
				s.Status = "done"
			}
			return nil
		}
		if s.Status == "pending" {
			s.Status = "running"
			s.Started = time.Now()
			return nil
		}
	}
	return fmt.Errorf("pipeline for %q has no pending or running stages", service)
}

// Get returns a copy of all stages for a service.
func (p *Pipeline) Get(service string) ([]*Stage, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	stages, ok := p.pipelines[service]
	if !ok {
		return nil, fmt.Errorf("no pipeline registered for service %q", service)
	}
	copy := make([]*Stage, len(stages))
	for i, s := range stages {
		tmp := *s
		copy[i] = &tmp
	}
	return copy, nil
}

// Reset clears all stages back to pending for a service.
func (p *Pipeline) Reset(service string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	stages, ok := p.pipelines[service]
	if !ok {
		return fmt.Errorf("no pipeline registered for service %q", service)
	}
	for _, s := range stages {
		s.Status = "pending"
		s.Started = time.Time{}
		s.Finished = time.Time{}
		s.Err = nil
	}
	return nil
}
