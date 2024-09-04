package entity

import (
	"context"
	"encoding/json"
)

const (
	ParentQueueName  = "th"
	ForwardQueueName = ParentQueueName + ":forward"
)

//go:generate easyjson task.go
//easyjson:json
type QueueTask struct {
	Queue   string
	Payload []byte
}

type ForwardTaskPayload struct {
	Name    string
	Options []byte
	Content []byte
}

type TaskHandlerFunc func(context.Context, QueueTask) error

func (ftp *ForwardTaskPayload) Payload() ([]byte, error) {
	return json.Marshal(ftp)
}
