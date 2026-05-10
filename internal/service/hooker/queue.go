package hooker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/internal/entity/hooks"
	"github.com/k1nky/tookhook/internal/entity/tasks"
)

func (svc *Service) Forward(ctx context.Context, r *hooks.Hook) error {
	rule := svc.rs.GetIncomeHookByName(ctx, r.Meta.Name)
	if rule == nil {
		return fmt.Errorf("hook %s: %w", r.Meta, entity.ErrNotFound)
	}
	if rule.Disabled {
		svc.log.Debugf("hook %s skipped", rule.Name)
		return nil
	}
	payload, err := json.Marshal(r)
	if err != nil {
		svc.log.Errorf("marshaling payload to %s failed: %v", r.Meta, err)
		return err
	}
	if err := svc.tq.Enqueue(ctx, &tasks.QueueTask{
		Queue:   tasks.HookQueueName,
		Payload: payload,
	}); err != nil {
		svc.log.Errorf("enqueue failed: %v", err)
		return err
	}
	return nil
}

func (svc *Service) processQueueTask(ctx context.Context, qt tasks.QueueTask) error {
	switch qt.Queue {
	case tasks.HookQueueName:
		r := hooks.Hook{}
		if err := json.Unmarshal(qt.Payload, &r); err != nil {
			return fmt.Errorf("%v: %w", err, entity.ErrSkipRetry)
		}
		if err := svc.processHook(ctx, &r); err != nil {
			return fmt.Errorf("%v: %w", err, entity.ErrSkipRetry)
		}
	case tasks.ForwardQueueName:
		t := tasks.ForwardTask{}
		if err := json.Unmarshal(qt.Payload, &t); err != nil {
			return fmt.Errorf("%v: %w", err, entity.ErrSkipRetry)
		}
		if err := svc.processForward(ctx, &t); err != nil {
			return err
		}
	}
	return nil
}
