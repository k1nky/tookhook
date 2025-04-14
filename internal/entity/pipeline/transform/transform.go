package transform

import (
	"context"
	"net/http"
	"strings"

	"github.com/k1nky/tookhook/pkg/httpclient"
	"github.com/k1nky/tookhook/pkg/thstrings/restrings"
	"github.com/k1nky/tookhook/pkg/thstrings/tstrings"
)

type Transformer interface {
	Compile() (err error)
	Execute(data []byte) ([]byte, error)
}

type StageType string

const (
	StageTypeDefault  StageType = ""
	StageTypeTemplate StageType = "template"
	StageTypeHTTP     StageType = "http"
)

func NewStageType(s string) StageType {
	return StageType(strings.ToLower(s))
}

type Stage struct {
	Type              StageType        `yaml:"type"`
	On                restrings.String `yaml:"on"`
	TemplateTransform *TemplateStage   `yaml:"template,omitempty"`
	HTTPTransform     *HTTPStage       `yaml:"http,omitempty"`
}

type Pipeline []*Stage

type TemplateStage struct {
	tstrings.String
}

type HTTPStage struct {
	URL tstrings.String `yaml:"url"`
}

func (ts Stage) GetTransformer() Transformer {
	switch ts.Type {
	case StageTypeTemplate:
		return ts.TemplateTransform
	case StageTypeHTTP:
		return ts.HTTPTransform
	}
	return nil
}

func (s *Stage) Compile() (err error) {
	if s.Type == StageTypeDefault {
		s.Type = StageTypeTemplate
	}
	if err = s.On.Compile(); err != nil {
		return
	}
	err = s.GetTransformer().Compile()
	return
}

func (s *Stage) Execute(data []byte) ([]byte, error) {
	return s.GetTransformer().Execute(data)
}

func (ht *HTTPStage) Compile() (err error) {
	err = ht.URL.Compile()
	return
}

func (ht *HTTPStage) Execute(data []byte) ([]byte, error) {
	uri, err := ht.URL.Execute(data)
	if err != nil {
		return nil, err
	}
	body, err := httpclient.SendRequest(context.Background(), http.MethodGet, string(uri), nil)
	if err != nil {
		return nil, err
	}

	return body, nil

}

func (tp Pipeline) Compile() error {
	for _, v := range tp {
		if err := v.Compile(); err != nil {
			return err
		}
	}
	return nil
}

func (tp Pipeline) Execute(data []byte) (transformed []byte, err error) {
	transformed = data
	for _, v := range tp {
		if transformed, err = v.Execute(transformed); err != nil {
			return nil, err
		}
	}
	return
}
