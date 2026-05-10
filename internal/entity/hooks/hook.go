package hooks

import (
	"encoding/json"
	"fmt"
)

type ContentEncoding int

const (
	ContentJSON ContentEncoding = iota
	ContentText
)

//go:generate easyjson hook.go
//easyjson:json
type Hook struct {
	Meta        Meta
	Args        map[string]string
	RawBody     []byte
	Encoding    ContentEncoding
	Headers     map[string]string
	decodedBody any
}

type Meta struct {
	ID   uint64
	Name string
}

func (m Meta) String() string {
	return fmt.Sprintf("%s[%d]", m.Name, m.ID)
}

func (h *Hook) Body() any {
	if h.decodedBody != nil {
		return h.decodedBody
	}
	if h.Encoding == ContentText {
		return string(h.RawBody)
	}
	var d any
	d = map[string]interface{}{}
	if err := json.Unmarshal(h.RawBody, &d); err == nil {
		return d
	}
	d = []interface{}{}
	if err := json.Unmarshal(h.RawBody, &d); err == nil {
		return d
	}
	return string(h.RawBody)
}
