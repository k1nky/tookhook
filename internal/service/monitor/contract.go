package monitor

import (
	"context"

	"github.com/k1nky/tookhook/internal/entity"
)

type pluginmanager interface {
	Health(ctx context.Context) entity.PluginsStatus
}

type hookService interface {
	Health(ctx context.Context) entity.Status
}
