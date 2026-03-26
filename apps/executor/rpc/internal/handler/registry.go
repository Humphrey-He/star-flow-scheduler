package handler

import (
	"fmt"
	"sync"
)

type Registry struct {
	mu       sync.RWMutex
	handlers map[string]JobHandler
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]JobHandler),
	}
}

func (r *Registry) Register(h JobHandler) error {
	if h == nil {
		return fmt.Errorf("handler is nil")
	}
	name := h.Name()
	if name == "" {
		return fmt.Errorf("handler name is empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.handlers[name]; ok {
		return fmt.Errorf("handler already registered: %s", name)
	}
	r.handlers[name] = h
	return nil
}

func (r *Registry) Get(name string) (JobHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	h, ok := r.handlers[name]
	return h, ok
}
