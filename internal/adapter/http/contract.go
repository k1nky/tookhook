package http

import (
	"context"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/internal/entity/hooks"
)

type logger interface {
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}

//go:generate mockgen -source=contract.go -destination=mock/hooker.go -package=mock hookService
type hookService interface {
	Forward(ctx context.Context, r *hooks.Hook) error
}

//go:generate mockgen -source=contract.go -destination=mock/hooker.go -package=mock rulesService
type rulesService interface {
	Load(ctx context.Context) error
}

//go:generate mockgen -source=contract.go -destination=mock/hooker.go -package=mock monitorService
type monitorService interface {
	Status(ctx context.Context) entity.ServiceStatus
}
