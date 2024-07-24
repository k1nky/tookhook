package monitor

import (
	"context"

	"github.com/k1nky/tookhook/internal/entity"
)

type Service struct {
	hookerSvc hookService
	log       logger
	pm        pluginmanager
	store     storage
}

func New(pm pluginmanager, hookerSvc hookService, store storage, log logger) *Service {
	return &Service{
		hookerSvc: hookerSvc,
		log:       log,
		pm:        pm,
		store:     store,
	}
}

func (svc *Service) Status(ctx context.Context) entity.ServiceStatus {
	status := svc.hookerSvc.Health(ctx)
	pluginsStatus := svc.pm.Health(ctx)
	for _, v := range pluginsStatus {
		if v == entity.StatusFailed {
			status = entity.StatusFailed
			break
		}
	}
	return entity.ServiceStatus{
		Status:  status,
		Plugins: pluginsStatus,
	}
}
