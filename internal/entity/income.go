package entity

type ContentType int

type IncomeRequest struct {
	Type ContentType
	Body []byte
}

const (
	ContentTypeJSON ContentType = iota
	ContentTypePlainText
)
