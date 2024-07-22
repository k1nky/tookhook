package hooker

import (
	"context"
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
)

const (
	ReceiverTypeLog = "log"
)

type Service struct {
	store storage
	pm    pluginmanager
	log   logger
}

func New(store storage, pm pluginmanager, log logger) *Service {
	return &Service{
		store: store,
		pm:    pm,
		log:   log,
	}
}

func (svc *Service) Reload(ctx context.Context) error {
	err := svc.store.ReadRules(ctx)
	if err != nil {
		svc.log.Errorf("reload rules: %v", err)
	} else {
		svc.log.Debugf("reload rules: success")
	}
	return err
}

func (svc *Service) Forward(ctx context.Context, name string, data []byte) error {
	rule, err := svc.store.GetIncomeHookByName(ctx, name)
	if err != nil {
		return err
	}
	if rule == nil {
		return fmt.Errorf("income rule %s: %w", name, entity.ErrNotFound)
	}
	if rule.Disabled {
		svc.log.Debugf("rule %s skipped", rule.Income)
		return nil
	}
	for _, r := range rule.Outcome {
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
			}
		}
	}
	return nil
}

func (svc *Service) Health(ctx context.Context) entity.Status {
	return entity.StatusOk
}
