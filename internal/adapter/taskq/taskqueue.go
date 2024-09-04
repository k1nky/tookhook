package taskq

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/k1nky/tookhook/internal/entity"
)

type Adapter struct {
	client      *asynq.Client
	server      *asynq.Server
	log         logger
	handler     entity.TaskHandlerFunc
	parentQueue string
}

func New(addr string, parentQueue string, log logger) *Adapter {
	return &Adapter{
		client: asynq.NewClient(asynq.RedisClientOpt{
			Addr: addr,
		}),
		server: asynq.NewServer(
			asynq.RedisClientOpt{
				Addr: addr,
			},
			asynq.Config{
				Concurrency: 10,
				Logger:      log,
			}),
		log:         log,
		parentQueue: parentQueue,
	}
}

func (a *Adapter) Process(ctx context.Context, handler entity.TaskHandlerFunc) error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(a.parentQueue, a.processTask)
	a.handler = handler
	if err := a.server.Start(mux); err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		a.server.Shutdown()
	}()
	return nil
}

func (a *Adapter) processTask(ctx context.Context, t *asynq.Task) error {
	qt := entity.QueueTask{
		Queue:   t.Type(),
		Payload: t.Payload(),
	}
	return a.handler(ctx, qt)
}

func (a *Adapter) Enqueue(ctx context.Context, queueTask *entity.QueueTask) error {
	t := asynq.NewTask(queueTask.Queue, queueTask.Payload)
	ti, err := a.client.EnqueueContext(ctx, t)
	if err != nil {
		a.log.Debugf("new task %s into %s", ti.ID, ti.Queue)
	}
	return err
}
