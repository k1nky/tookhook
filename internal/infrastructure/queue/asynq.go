// Package queue provides implementations of the task queue.
package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/k1nky/tookhook/internal/domain/entity"
)

const (
	// TaskTypeWebhook is the task type for webhook processing.
	TaskTypeWebhook = "webhook"
	// DefaultMaxRetry is the default maximum number of retries for a task.
	DefaultMaxRetry = 20
)

// AsynqQueue implements the Queue interface using asynq.
type AsynqQueue struct {
	client      *asynq.Client
	server      *asynq.Server
	mux         *asynq.ServeMux
	concurrency int
}

// NewAsynqQueue creates a new asynq-based queue.
func NewAsynqQueue(addr string, db int, concurrency int) *AsynqQueue {
	redisOpt := asynq.RedisClientOpt{
		Addr: addr,
		DB:   db,
	}

	return &AsynqQueue{
		client: asynq.NewClient(redisOpt),
		server: asynq.NewServer(redisOpt, asynq.Config{
			Concurrency: concurrency,
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				// Exponential backoff: 10s, 20s, 40s, ...
				return time.Duration(10*(1<<uint(n))) * time.Second
			},
		}),
		mux:         asynq.NewServeMux(),
		concurrency: concurrency,
	}
}

// Enqueue adds a webhook task to the queue.
func (q *AsynqQueue) Enqueue(ctx context.Context, task *entity.WebhookTask) error {
	payload, err := task.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	asynqTask := asynq.NewTask(TaskTypeWebhook, payload, asynq.MaxRetry(DefaultMaxRetry))
	_, err = q.client.EnqueueContext(ctx, asynqTask)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	return nil
}

// Start begins processing tasks from the queue.
func (q *AsynqQueue) Start(ctx context.Context, handler func(ctx context.Context, task *entity.WebhookTask) error) error {
	q.mux.HandleFunc(TaskTypeWebhook, func(ctx context.Context, t *asynq.Task) error {
		task, err := entity.UnmarshalWebhookTask(t.Payload())
		if err != nil {
			// Don't retry if the payload is malformed
			return fmt.Errorf("failed to unmarshal task: %w: %v", err, asynq.SkipRetry)
		}

		return handler(ctx, task)
	})

	go func() {
		<-ctx.Done()
		q.server.Shutdown()
	}()

	if err := q.server.Start(q.mux); err != nil {
		return fmt.Errorf("failed to start queue server: %w", err)
	}

	return nil
}

// Close closes the queue connections.
func (q *AsynqQueue) Close() error {
	return q.client.Close()
}
