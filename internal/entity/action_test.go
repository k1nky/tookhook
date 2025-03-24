package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActions_Execute(t *testing.T) {
	tests := []struct {
		name       string
		a          Actions
		data       []byte
		wantError  error
		wantResult []byte
	}{
		{
			name:       "Not match",
			data:       []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`),
			wantError:  ErrNotMatch,
			wantResult: []byte(`{"@timestamp": "2024-07-11T12:40:31.574Z", "level": "INFO", "port": 37628, "host.name": "hostname", "message": "Select 1", "type": "app_log"}`),
			a: Actions{
				&Action{
					On: "NOT_MATCH",
				},
			},
		},
		{
			name:      "Discard",
			data:      []byte(`{"name": "Name", "data": "My Data"}`),
			wantError: ErrDiscard,
			a: Actions{
				&Action{
					Type: ActionTypeDiscard,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.Compile()
			got, err := tt.a.Execute(tt.data)
			assert.ErrorIs(t, err, tt.wantError)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}

func TestActions_Compile(t *testing.T) {
	tests := []struct {
		name      string
		a         Actions
		wantError error
	}{
		{
			name: "Empty actions",
			a: Actions{
				&Action{},
			},
		},
		{
			name: "Invalid On value",
			a: Actions{
				&Action{
					On: ")",
				},
			},
			wantError: ErrCompile,
		},
		{
			name: "Valid On value",
			a: Actions{
				&Action{
					On: ".*",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Compile()
			assert.ErrorIs(t, got, tt.wantError)
		})
	}

}
