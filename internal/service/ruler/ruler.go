// Package ruler defines the rules holder as a service.
package ruler

import (
	"context"

	"github.com/k1nky/tookhook/internal/entity"
)

type Service struct {
	log   logger
	pm    pluginmanager
	rules *entity.Rules
	store storage
}

// Returns new the rules holder service.
func New(pm pluginmanager, store storage, log logger) *Service {
	return &Service{
		log:   log,
		pm:    pm,
		store: store,
	}
}

// GetIncomeHookByName returns income hook definition by name.
func (svc *Service) GetIncomeHookByName(ctx context.Context, name string) *entity.Hook {
	hook := entity.Hook{}
	for _, v := range svc.rules.Hooks {
		if v.Income == name {
			hook = v
			return &hook
		}
	}
	// not found
	return nil
}

// Load loads the rules from a storage. The service stores the loaded rules in memory.
// Before setting new rules, the rules are validated.
func (svc *Service) Load(ctx context.Context) error {
	rules, err := svc.store.GetRules(ctx)
	if err != nil {
		svc.log.Errorf("load rules: %v", err)
		return err
	}
	if err := svc.Validate(ctx, rules); err != nil {
		svc.log.Errorf("load rules: %v", err)
		return err
	}
	svc.rules = rules
	svc.log.Debugf("reload rules: success")
	return nil
}

// Returns an error if specified rules is invalid.
func (svc *Service) Validate(ctx context.Context, rules *entity.Rules) error {
	if err := rules.Compile(); err != nil {
		return err
	}
	for _, hook := range rules.Hooks {
		for _, v := range hook.Handlers {
			p := svc.pm.Get(v.Type)
			if p == nil {
				continue
			}
			if err := p.Validate(ctx, v.AsPluginHandler()); err != nil {
				return err
			}
		}
	}
	return nil
}
