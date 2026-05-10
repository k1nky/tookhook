package rules

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/internal/entity/hooks"
	"github.com/k1nky/tookhook/pkg/httpclient"
	"github.com/k1nky/tookhook/pkg/thstrings/tstrings"
)

type Executer interface {
	Compile() (err error)
	Execute(data any) ([]byte, error)
}

type StageType string

const (
	StageTypeDefault  StageType = ""
	StageTypeTemplate StageType = "template"
	StageTypeHTTP     StageType = "http"
	StageTypeDiscard  StageType = "discard"
)

func NewStageType(s string) StageType {
	return StageType(strings.ToLower(s))
}

type Stage struct {
	Type              StageType       `yaml:"type"`
	On                tstrings.String `yaml:"on"`
	TemplateTransform *TemplateStage  `yaml:"template,omitempty"`
	HTTPTransform     *HTTPStage      `yaml:"http,omitempty"`
}

type Pipeline []*Stage

type TemplateStage struct {
	tstrings.String
}

type HTTPStage struct {
	URL          tstrings.String `yaml:"url"`
	TimeoutInSec uint            `yaml:"timeout"`
}

type DiscardStage struct{}

func (ts Stage) GetTransformer() Executer {
	switch ts.Type {
	case StageTypeTemplate:
		return ts.TemplateTransform
	case StageTypeHTTP:
		return ts.HTTPTransform
	case StageTypeDiscard:
		return &DiscardStage{}
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

func (s *Stage) Execute(data any) ([]byte, error) {
	return s.GetTransformer().Execute(data)
}

func (s *DiscardStage) Compile() error {
	return nil
}

func (s *DiscardStage) Execute(data any) ([]byte, error) {
	return nil, fmt.Errorf("transform: %w", entity.ErrDiscard)
}

func (ht *HTTPStage) Compile() (err error) {
	err = ht.URL.Compile()
	return
}

func (ht *HTTPStage) Execute(data any) ([]byte, error) {
	uri, err := ht.URL.Execute(data)
	if err != nil {
		return nil, err
	}
	r := httpclient.Request{
		Method:  http.MethodGet,
		URL:     string(uri),
		Timeout: time.Duration(ht.TimeoutInSec) * time.Second,
	}
	return httpclient.SendRequest(context.Background(), r, nil)
}

func (tp Pipeline) Compile() error {
	for _, v := range tp {
		if err := v.Compile(); err != nil {
			return err
		}
	}
	return nil
}

func (tp Pipeline) Execute(r *hooks.Hook) (transformed []byte, err error) {
	ps := &PipelineState{
		Hook: r,
		Data: r.RawBody,
	}
	for _, v := range tp {
		if ps.Data, err = v.Execute(ps); err != nil {
			return nil, err
		}
	}
	return ps.Data, nil
}

type PipelineState struct {
	*hooks.Hook
	Data []byte
}
