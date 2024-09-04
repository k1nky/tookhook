package hooker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/pkg/plugin"
)

type Service struct {
	rs  rulesStore
	pm  pluginmanager
	log logger
	q   taskqueue
}

func New(rs rulesStore, pm pluginmanager, log logger, q taskqueue) *Service {
	return &Service{
		rs:  rs,
		pm:  pm,
		log: log,
		q:   q,
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
			t := &entity.ForwardTaskPayload{
				Name:    h.Type,
				Options: h.AsPluginHandler().Options,
				Content: content,
			}
			payload, err := t.Payload()
			if err != nil {
				svc.log.Errorf("marshaling payload to %s failed: %v", h.Type, err)
				return err
			}
			if err := svc.q.Enqueue(ctx, &entity.QueueTask{
				Queue:   entity.ForwardQueueName,
				Payload: payload,
			}); err != nil {
				svc.log.Errorf("enqueue failed: %v", err)
				continue
			}
		}
	}
	return nil
}

func (svc *Service) Run(ctx context.Context) {
	svc.q.Process(ctx, svc.processQueueTask)
}

func (svc *Service) processQueueTask(ctx context.Context, qt entity.QueueTask) error {
	switch qt.Queue {
	case entity.ForwardQueueName:
		ftp := entity.ForwardTaskPayload{}
		if err := json.Unmarshal(qt.Payload, &ftp); err != nil {
			// 	return fmt.Errorf("%v: %w", err, asynq.SkipRetry)
			return err
		}
		fwd := svc.pm.Get(ftp.Name)
		_, err := fwd.Forward(ctx, plugin.Handler{
			Options: ftp.Options,
		}, ftp.Content)
		if err != nil {
			svc.log.Errorf("task processing failed: %v", err)
		}
		return err
	}
	return nil
}

func (svc *Service) Health(ctx context.Context) entity.Status {
	return entity.StatusOk
}
