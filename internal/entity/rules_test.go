package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestParseRules(t *testing.T) {
	rules := `
templates:
  jira: &jira
    - template: "{{ .message }}"
hooks:
  - income: test
    outcome:
      - type: pachca
        template: *jira
        options:
          chat: discussion/9913735

`
	r := Rules{}
	err := yaml.Unmarshal([]byte(rules), &r)
	assert.NoError(t, err)
	assert.NotNil(t, r)
}

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
						Income:  "",
						Outcome: []Receiver{{Type: "log"}},
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
						Income:  "test",
						Outcome: []Receiver{{Type: ""}},
					},
				},
			},
			want: ErrEmptyValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rules.Validate()
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
						Income:  "test",
						Outcome: []Receiver{{Type: "log"}},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rules.Validate()
			assert.NoError(t, got)
		})
	}
}

func TestReceiverContentWithTemplate(t *testing.T) {
	r := Receiver{
		Template: Templates{Template{Template: "{{ .message }}"}},
	}
	data := []byte(`{"message": "Message", "text": "Text"}`)
	content, err := r.Content(data)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Message"), content)
}

func TestReceiverContentWithoutTemplate(t *testing.T) {
	r := Receiver{}
	data := []byte(`My Message`)
	content, err := r.Content(data)
	assert.NoError(t, err)
	assert.Equal(t, []byte("My Message"), content)
}
