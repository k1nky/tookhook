// Package monitor defines a service for monitoring the status of internal components of the tookhook server.
package monitor

import (
	"context"

	"github.com/k1nky/tookhook/internal/entity"
)

// Service for monitoring the status of internal components of the tookhook server.
type Service struct {
	hookerSvc hookService
	log       logger
	pm        pluginmanager
	store     storage
}

// Return new instance of service.
func New(pm pluginmanager, hookerSvc hookService, store storage, log logger) *Service {
	return &Service{
		hookerSvc: hookerSvc,
		log:       log,
		pm:        pm,
		store:     store,
	}
}

// Return the current status of the tookhook server.
func (svc *Service) Status(ctx context.Context) entity.ServiceStatus {
	status := svc.hookerSvc.Health(ctx)
	pluginsStatus := svc.pm.Health(ctx)
	for _, v := range pluginsStatus {
		if v == entity.StatusFailed {
			// mark the server as failed if there is any failed plugin
			status = entity.StatusFailed
			break
		}
	}
	return entity.ServiceStatus{
		Status:  status,
		Plugins: pluginsStatus,
	}
}
