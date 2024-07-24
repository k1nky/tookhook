package ruler

import (
	"context"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/pkg/plugin"
)

//go:generate mockgen -source=contract.go -destination=mock/rules.go -package=mock pluginmanager
type pluginmanager interface {
	Get(name string) plugin.Plugin
	Health(ctx context.Context) entity.PluginsStatus
}

//go:generate mockgen -source=contract.go -destination=mock/rules.go -package=mock storage
type storage interface {
	GetRules(ctx context.Context) (*entity.Rules, error)
}

type logger interface {
	Debugf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}
