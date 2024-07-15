package hooker

import (
	"context"
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

func (svc *Service) Forward(ctx context.Context, name string, data []byte) error {
	rule, err := svc.store.GetIncomeHookByName(ctx, name)
	if err != nil {
		return err
	}
	if rule == nil {
		// TODO: not found error
		return nil
	}
	if rule.Disabled {
		svc.log.Debugf("rule %s skipped", rule.Income)
		return nil
	}
	for _, r := range rule.Outcome {
		if r.Disabled {
			svc.log.Debugf("reciever %s %s skipped", r.Type, r.Target)
			continue
		}
		fwd := svc.pm.Get(r.Type)
		if fwd != nil {
			content, err := r.Content(data)
			if err != nil {
				return err
			}
			if _, err := fwd.Forward(ctx, r, content); err != nil {
				svc.log.Errorf("send to %s %s failed: %v", r.Type, r.Target, err)
			}
		}
	}
	return nil
}
