package pachca

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessageSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(`{"message_id":"1"}`))
	}))
	defer ts.Close()
	c := NewPachca("token")
	c.Endpoint = ts.URL
	payload := MessagePayload{
		Message: Message{},
	}
	response, err := c.Send(payload)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"message_id":"1"}`), response)
}

func TestNewMessageFailed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"error":"error message"}`))
	}))
	defer ts.Close()
	c := NewPachca("token")
	c.Endpoint = ts.URL
	payload := MessagePayload{
		Message: Message{},
	}
	response, err := c.Send(payload)
	assert.Error(t, err)
	assert.Equal(t, []byte(`{"error":"error message"}`), response)
}

func TestNewMessageNoContent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()
	c := NewPachca("token")
	c.Endpoint = ts.URL
	payload := MessagePayload{
		Message: Message{},
	}
	response, err := c.Send(payload)
	assert.Error(t, err)
	assert.Equal(t, []byte{}, response)
}
