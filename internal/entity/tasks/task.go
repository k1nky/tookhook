package tasks

import (
	"context"
	"encoding/json"

	"github.com/k1nky/tookhook/internal/entity/hooks"
)

const (
	ParentQueueName  = "th"
	ForwardQueueName = ParentQueueName + ":forward"
	HookQueueName    = ParentQueueName + ":hook"
)

//go:generate easyjson task.go
//easyjson:json
type QueueTask struct {
	Queue   string
	Payload []byte
}

//go:generate easyjson task.go
//easyjson:json
type ForwardTask struct {
	Name    string
	Hook    hooks.HookRequestMeta
	Options []byte
	Content []byte
}

//go:generate easyjson task.go
//easyjson:json
type HookTask struct {
	Hook hooks.HookRequestMeta
	Data []byte
}

type TaskHandlerFunc func(context.Context, QueueTask) error

func (ftp *ForwardTask) Payload() ([]byte, error) {
	return json.Marshal(ftp)
}

func (t *HookTask) Payload() ([]byte, error) {
	return json.Marshal(t)
}
