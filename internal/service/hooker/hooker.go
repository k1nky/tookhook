package hooker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/internal/entity/hooks"
	"github.com/k1nky/tookhook/internal/entity/tasks"
	"github.com/k1nky/tookhook/pkg/plugin"
)

type Service struct {
	rs  rulesStore
	pm  pluginmanager
	log logger
	tq  taskqueue
}

func New(rs rulesStore, pm pluginmanager, log logger, tq taskqueue) *Service {
	return &Service{
		rs:  rs,
		pm:  pm,
		log: log,
		tq:  tq,
	}
}

func (svc *Service) Forward(ctx context.Context, r *hooks.HookRequest) error {
	rule := svc.rs.GetIncomeHookByName(ctx, r.Meta.Name)
	if rule == nil {
		return fmt.Errorf("hook %s: %w", r.Meta, entity.ErrNotFound)
	}
	if rule.Disabled {
		svc.log.Debugf("hook %s skipped", rule.Income)
		return nil
	}
	// TODO: verify income request?
	t := tasks.HookTask{
		Hook: r.Meta,
		Data: r.Content.Body,
	}
	payload, err := t.Payload()
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

func (svc *Service) processHook(ctx context.Context, t *tasks.HookTask) error {
	hook := svc.rs.GetIncomeHookByName(ctx, t.Hook.Name)
	if hook == nil || hook.Disabled {
		svc.log.Debugf("hook %s skipped: not found or disabled", t.Hook)
		return nil
	}
	for _, h := range hook.Handlers {
		if h.Disabled {
			svc.log.Debugf("handler %s %s skipped", t.Hook, h.Type)
			continue
		}
		if !h.Match(t.Data) {
			svc.log.Debugf("handler %s %s skipped", t.Hook, h.Type)
			continue
		}
		content, err := h.Content(t.Data)
		if err != nil {
			if errors.Is(err, entity.ErrDiscard) {
				svc.log.Debugf("handler %s %s discarded", t.Hook, h.Type)
				continue
			}
			if !errors.Is(err, entity.ErrNotMatch) {
				svc.log.Errorf("handler %s %s: %v", t.Hook, h.Type, err)
				continue
			}
		}
		t := &tasks.ForwardTask{
			Name:    h.Type,
			Hook:    t.Hook,
			Options: h.AsPluginHandler().Options,
			Content: content,
		}
		payload, err := t.Payload()
		if err != nil {
			svc.log.Errorf("marshaling payload to %s/%s failed: %v", t.Hook, h.Type, err)
			continue
		}
		if err := svc.tq.Enqueue(ctx, &tasks.QueueTask{
			Queue:   tasks.ForwardQueueName,
			Payload: payload,
		}); err != nil {
			svc.log.Errorf("enqueue a forward task %s/%s failed: %v", t.Hook, h.Type, err)
			continue
		}
	}
	return nil
}

func (svc *Service) Run(ctx context.Context) {
	svc.tq.Process(ctx, svc.processQueueTask)
}

func (svc *Service) processForward(ctx context.Context, t *tasks.ForwardTask) error {
	if h := svc.pm.Get(t.Name); h != nil {
		svc.log.Debugf("call plugin %s for %s", t.Name, t.Hook)
		if _, err := h.Forward(ctx, plugin.Handler{
			Options: t.Options,
		}, t.Content); err != nil {
			return err
		}
	} else {
		svc.log.Errorf("plugin %s not found", t.Name)
	}
	return nil
}

func (svc *Service) processQueueTask(ctx context.Context, qt tasks.QueueTask) error {
	switch qt.Queue {
	case tasks.HookQueueName:
		t := tasks.HookTask{}
		if err := json.Unmarshal(qt.Payload, &t); err != nil {
			return fmt.Errorf("%v: %w", err, entity.ErrSkipRetry)
		}
		if err := svc.processHook(ctx, &t); err != nil {
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

func (svc *Service) Health(ctx context.Context) entity.Status {
	return entity.StatusOk
}
