package builtin

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/k1nky/tookhook/internal/domain/entity"
)

// HandlerExecutor defines the interface for handler execution.
type HandlerExecutor interface {
	Execute(ctx context.Context, input []byte, opts map[string]any) ([]byte, error)
	Name() string
}

// Registry maintains a collection of registered handlers.
type Registry struct {
	mu       sync.RWMutex
	handlers map[string]HandlerExecutor
	logger   *slog.Logger
}

// NewRegistry creates a new handler registry with built-in handlers registered.
func NewRegistry(logger *slog.Logger) *Registry {
	r := &Registry{
		handlers: make(map[string]HandlerExecutor),
		logger:   logger,
	}

	// Register built-in handlers
	r.Register(NewLogHandler(logger))
	r.Register(NewHttpHandler(logger))
	r.Register(NewExecHandler(logger))
	r.Register(NewTextHandler(logger))
	r.Register(NewJsHandler(logger))

	return r
}

// Register adds a handler to the registry.
func (r *Registry) Register(h HandlerExecutor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[h.Name()] = h
}

// Get returns a handler by name.
func (r *Registry) Get(name string) (HandlerExecutor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	h, ok := r.handlers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", entity.ErrHandlerNotFound, name)
	}
	return h, nil
}

// Names returns all registered handler names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		names = append(names, name)
	}
	return names
}
