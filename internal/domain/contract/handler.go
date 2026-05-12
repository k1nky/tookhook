// Package contract defines interfaces for the webhook processing service.
package contract

import (
	"context"
)

// HandlerExecutor defines the interface for executing a handler.
// Each handler type (~log, ~http, ~exec) must implement this interface.
type HandlerExecutor interface {
	// Execute processes the input data and returns the output.
	// The output is passed to the next handler in the chain.
	Execute(ctx context.Context, input []byte, opts map[string]any) ([]byte, error)
	// Name returns the handler type name (e.g., "~log", "~http").
	Name() string
}

// HandlerRegistry defines the interface for handler registry.
// It allows looking up handler executors by their type name.
type HandlerRegistry interface {
	// Get returns the handler executor for the given type.
	// Returns ErrHandlerNotFound if the handler is not registered.
	Get(handlerType string) (HandlerExecutor, error)
	// Register adds a handler executor to the registry.
	Register(executor HandlerExecutor)
}
