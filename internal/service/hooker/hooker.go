package hooker

import (
	"context"
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
)

const (
	ReceiverTypeLog = "!log"
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
	for _, r := range rule.Handlers {
		if r.Disabled {
			svc.log.Debugf("reciever %s %s skipped", r.Type)
			continue
		}
		content, err := r.Content(data)
		if err != nil {
			return err
		}
		if r.Type == ReceiverTypeLog {
			svc.log.Debugf(string(content))
			continue
		}
		fwd := svc.pm.Get(r.Type)
		if fwd != nil {
			if _, err := fwd.Forward(ctx, r.AsPluginReceiver(), content); err != nil {
				svc.log.Errorf("send to %s %s failed: %v", r.Type, err)
				return err
			}
		}
	}
	return nil
}

func (svc *Service) Health(ctx context.Context) entity.Status {
	return entity.StatusOk
}
