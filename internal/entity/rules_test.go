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
        target: discussion/9913735        
        token: H7sBBoqmEPb6CMVBEgInuszyvYMSIZZ_K83Uhrvl0RQ

`
	r := Rules{}
	err := yaml.Unmarshal([]byte(rules), &r)
	assert.NoError(t, err)
	assert.NotNil(t, r)
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
