package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRulesValidateFailed(t *testing.T) {
	tests := []struct {
		name  string
		rules Rules
		want  error
	}{
		{
			name: "EmptyIncome",
			rules: Rules{
				Hooks: []Hook{
					{
						Income:   "",
						Handlers: []*Handler{{Type: "~log"}},
					},
				},
			},
			want: ErrEmptyValue,
		},
		{
			name: "EmptyOutcome",
			rules: Rules{
				Hooks: []Hook{
					{
						Income:   "test",
						Handlers: []*Handler{{Type: ""}},
					},
				},
			},
			want: ErrEmptyValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rules.Compile()
			assert.ErrorIs(t, got, tt.want)
		})
	}
}

func TestRulesValidateNoError(t *testing.T) {
	tests := []struct {
		name  string
		rules Rules
	}{
		{
			name: "NoError",
			rules: Rules{
				Hooks: []Hook{
					{
						Income:   "test",
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

func TestHandlerContentWithTemplate(t *testing.T) {
	h := Handler{
		Type:         "handler1",
		PreTransform: Transforms{&Transform{Template: "{{ .message }}"}},
	}
	h.Compile()
	data := []byte(`{"message": "Message", "text": "Text"}`)
	content, err := h.Content(data)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Message"), content)
}

func TestHandlerContentWithoutTemplate(t *testing.T) {
	h := Handler{}
	data := []byte(`My Message`)
	content, err := h.Content(data)
	assert.NoError(t, err)
	assert.Equal(t, []byte("My Message"), content)
}
