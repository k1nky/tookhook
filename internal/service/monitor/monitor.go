package monitor

import (
	"context"

	"github.com/k1nky/tookhook/internal/entity"
)

type Service struct {
	pm        pluginmanager
	hookerSvc hookService
}

func New(pm pluginmanager, hookerSvc hookService) *Service {
	return &Service{
		pm:        pm,
		hookerSvc: hookerSvc,
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
