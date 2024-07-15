package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReceiverContentWithTemplate(t *testing.T) {
	r := Receiver{
		Template: "{{ .message }}",
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
