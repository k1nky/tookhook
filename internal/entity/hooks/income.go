package hooks

import "fmt"

type ContentType int

type HookRequest struct {
	Meta    HookRequestMeta
	Content HookRequestBody
}

type HookRequestMeta struct {
	ID   uint64
	Name string
}

type HookRequestBody struct {
	Type string
	Body []byte
}

func (m HookRequestMeta) String() string {
	return fmt.Sprintf("%s[%d]", m.Name, m.ID)
}
