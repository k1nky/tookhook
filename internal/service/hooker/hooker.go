package hooker

import (
	"context"
	"errors"

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

func (svc *Service) processHook(ctx context.Context, r *hooks.Hook) error {
	hook := svc.rs.GetIncomeHookByName(ctx, r.Meta.Name)
	if hook == nil || hook.Disabled {
		svc.log.Debugf("hook %s skipped: not found or disabled", r.Meta)
		return nil
	}
	for _, h := range hook.Handlers {
		if h.Disabled {
			svc.log.Debugf("handler %s %s skipped", r.Meta, h.Type)
			continue
		}
		if matched, err := h.On.Match(r); !matched {
			svc.log.Debugf("handler %s %s skipped", r.Meta, h.Type)
			continue
		} else if err != nil {
			svc.log.Errorf("handler %s %s skipped due to: %s", r.Meta, h.Type, err)
			continue
		}
		content, err := h.Execute(r)
		if err != nil {
			if errors.Is(err, entity.ErrDiscard) {
				svc.log.Debugf("handler %s %s discarded", r.Meta, h.Type)
				continue
			}
			if !errors.Is(err, entity.ErrNotMatch) {
				svc.log.Errorf("handler %s %s: %v", r.Meta, h.Type, err)
				continue
			}
		}
		t := &tasks.ForwardTask{
			Name:    h.Type,
			Hook:    r.Meta,
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

func (svc *Service) Health(ctx context.Context) entity.Status {
	return entity.StatusOk
}
