// Package repository provides implementations of domain repositories.
package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/k1nky/tookhook/internal/domain/entity"
	"github.com/k1nky/tookhook/internal/infrastructure/config"
)

// FileRepository implements EndpointRepository using a YAML file.
type FileRepository struct {
	mu         sync.RWMutex
	endpoints  map[string]*entity.Endpoint
	configFile string
}

// NewFileRepository creates a new file repository with pre-loaded endpoints.
func NewFileRepository(configFile string) *FileRepository {
	fr := &FileRepository{
		endpoints:  make(map[string]*entity.Endpoint),
		configFile: configFile,
	}
	return fr
}

// GetByName returns the endpoint by name.
func (r *FileRepository) GetByName(_ context.Context, name string) (*entity.Endpoint, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ep, ok := r.endpoints[name]
	if !ok {
		return nil, entity.ErrEndpointNotFound
	}
	return ep, nil
}

// GetAll returns all endpoints.
func (r *FileRepository) GetAll(_ context.Context) ([]*entity.Endpoint, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*entity.Endpoint, 0, len(r.endpoints))
	for _, ep := range r.endpoints {
		result = append(result, ep)
	}
	return result, nil
}

// Reload reloads endpoints from the configuration file.
func (r *FileRepository) Reload(_ context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Load config from file
	cfg, err := config.Load(r.configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Convert to entities
	endpoints, err := cfg.ToEntities()
	if err != nil {
		return fmt.Errorf("failed to convert to entities: %w", err)
	}

	// Update endpoints
	r.endpoints = make(map[string]*entity.Endpoint)
	for i := range endpoints {
		ep := &endpoints[i]
		r.endpoints[ep.Name] = ep
	}

	return nil
}
