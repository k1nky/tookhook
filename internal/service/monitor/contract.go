package monitor

import (
	"context"

	"github.com/k1nky/tookhook/internal/entity"
)

//go:generate mockgen -source=contract.go -destination=mock/monitor.go -package=mock pluginmanager
type pluginmanager interface {
	Health(ctx context.Context) entity.PluginsStatus
}

//go:generate mockgen -source=contract.go -destination=mock/monitor.go -package=mock hookservice
type hookService interface {
	Health(ctx context.Context) entity.Status
}

//go:generate mockgen -source=contract.go -destination=mock/monitor.go -package=mock storage
type storage interface{}

type logger interface {
	Debugf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}
