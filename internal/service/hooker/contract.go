package hooker

import (
	"context"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/pkg/plugin"
)

type storage interface {
	GetIncomeHookByName(ctx context.Context, name string) (*entity.Hook, error)
	ReadRules(ctx context.Context) error
}

type logger interface {
	Debugf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type pluginmanager interface {
	Get(name string) *plugin.GRPCClient
}
