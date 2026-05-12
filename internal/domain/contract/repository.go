package contract

import (
	"context"

	"github.com/k1nky/tookhook/internal/domain/entity"
)

// EndpointRepository defines the interface for endpoint configuration storage.
// This is implemented by the file repository (YAML config).
type EndpointRepository interface {
	// GetByName returns the endpoint configuration for the given name.
	// Returns ErrEndpointNotFound if the endpoint does not exist.
	GetByName(ctx context.Context, name string) (*entity.Endpoint, error)
	// GetAll returns all configured endpoints.
	GetAll(ctx context.Context) ([]*entity.Endpoint, error)
	// Reload refreshes the configuration from the source.
	Reload(ctx context.Context) error
}
