package entity

import "encoding/json"

// TaskType represents the type of task in the queue.
const TaskTypeWebhook = "webhook"

// WebhookTask represents a webhook processing task in the queue.
type WebhookTask struct {
	// EndpointName is the name of the endpoint to process.
	EndpointName string `json:"endpoint_name"`
	// Payload is the original webhook request body.
	Payload []byte `json:"payload"`
	// Headers contains the original request headers.
	Headers map[string][]string `json:"headers,omitempty"`
	// ContentType is the content-type of the original request.
	ContentType string `json:"content_type,omitempty"`
	// ID is the unique request ID.
	ID string `json:"id"`
}

// Marshal serializes the webhook task to JSON bytes.
func (t *WebhookTask) Marshal() ([]byte, error) {
	return json.Marshal(t)
}

// UnmarshalWebhookTask deserializes JSON bytes to a WebhookTask.
func UnmarshalWebhookTask(data []byte) (*WebhookTask, error) {
	var task WebhookTask
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}
