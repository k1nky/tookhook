package hooker

import (
	"context"
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
)

type Service struct {
	rs  rulesStore
	pm  pluginmanager
	log logger
}

func New(rs rulesStore, pm pluginmanager, log logger) *Service {
	return &Service{
		rs:  rs,
		pm:  pm,
		log: log,
	}
}

func (svc *Service) Forward(ctx context.Context, name string, data []byte) error {
	rule := svc.rs.GetIncomeHookByName(ctx, name)
	if rule == nil {
		return fmt.Errorf("income rule %s: %w", name, entity.ErrNotFound)
	}
	if rule.Disabled {
		svc.log.Debugf("rule %s skipped", rule.Income)
		return nil
	}
	for _, h := range rule.Handlers {
		if h.Disabled {
			svc.log.Debugf("handler %s %s skipped", h.Type)
			continue
		}
		if !h.Match(data) {
			svc.log.Debugf("handler %s %s skipped", h.Type)
			continue
		}
		content, err := h.Content(data)
		if err != nil {
			return err
		}
		fwd := svc.pm.Get(h.Type)
		if fwd != nil {
			if _, err := fwd.Forward(ctx, h.AsPluginHandler(), content); err != nil {
				svc.log.Errorf("send to %s %s failed: %v", h.Type, err)
				return err
			}
		}
	}
	return nil
}

func (svc *Service) Health(ctx context.Context) entity.Status {
	return entity.StatusOk
}
