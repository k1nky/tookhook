package rules

import (
	"testing"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestRules_Compile_Failed(t *testing.T) {
	tests := []struct {
		name  string
		rules Rules
		want  error
	}{
		{
			name: "EmptyIncome",
			rules: Rules{
				Endpoints: []Endpoint{
					{
						Name:     "",
						Handlers: []*Handler{{Type: "~log"}},
					},
				},
			},
			want: entity.ErrEmptyValue,
		},
		{
			name: "EmptyOutcome",
			rules: Rules{
				Endpoints: []Endpoint{
					{
						Name:     "test",
						Handlers: []*Handler{{Type: ""}},
					},
				},
			},
			want: entity.ErrEmptyValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rules.Compile()
			assert.ErrorIs(t, got, tt.want)
		})
	}
}

func TestRules_Compile_NoError(t *testing.T) {
	tests := []struct {
		name  string
		rules Rules
	}{
		{
			name: "NoError",
			rules: Rules{
				Endpoints: []Endpoint{
					{
						Name:     "test",
						Handlers: []*Handler{{Type: "log"}},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rules.Compile()
			assert.NoError(t, got)
		})
	}
}
