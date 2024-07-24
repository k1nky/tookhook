package mock

import (
	"context"

	"github.com/k1nky/tookhook/pkg/plugin"
)

type MockPlugin struct {
	ValidateResult     error
	HealthResult       error
	ForwardResultData  []byte
	ForwardResultError error
}

func (m *MockPlugin) Validate(ctx context.Context, r plugin.Receiver) error {
	return m.ValidateResult
}

func (m *MockPlugin) Health(ctx context.Context) error {
	return m.HealthResult
}

func (m *MockPlugin) Forward(ctx context.Context, r plugin.Receiver, data []byte) ([]byte, error) {
	return m.ForwardResultData, m.ForwardResultError
}
