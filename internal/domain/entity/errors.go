package entity

import "errors"

var (
	// ErrEmptyEndpointName indicates that endpoint name is empty.
	ErrEmptyEndpointName = errors.New("endpoint name cannot be empty")
	// ErrEmptyChain indicates that chain has no handlers.
	ErrEmptyChain = errors.New("chain cannot be empty")
	// ErrEmptyHandlerType indicates that handler type is empty.
	ErrEmptyHandlerType = errors.New("handler type cannot be empty")
	// ErrEndpointNotFound indicates that endpoint was not found.
	ErrEndpointNotFound = errors.New("endpoint not found")
	// ErrHandlerNotFound indicates that handler type is not registered.
	ErrHandlerNotFound = errors.New("handler not found")
)
