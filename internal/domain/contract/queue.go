package contract

import (
	"context"

	"github.com/k1nky/tookhook/internal/domain/entity"
)

// Queue defines the interface for the task queue.
// This is implemented by the asynq adapter.
type Queue interface {
	// Enqueue adds a webhook task to the queue.
	Enqueue(ctx context.Context, task *entity.WebhookTask) error
	// Start begins processing tasks from the queue.
	// It blocks until the context is cancelled.
	Start(ctx context.Context, handler func(ctx context.Context, task *entity.WebhookTask) error) error
}
