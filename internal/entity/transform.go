package entity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"text/template"

	"github.com/k1nky/tookhook/pkg/httpclient"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TransformDataType string

const (
	DataJSON  TransformDataType = "json"
	DataPlain TransformDataType = "plain"
)

type transform struct {
	regexp   *regexp.Regexp
	template *template.Template
}

type Transform struct {
	transform
	// RegExp will be applied to the data before the template is executed.
	RegExp string `yaml:"regexp"`
	// Template will be applied to the data. The string must be formatted as a text/template package template.
	Template string            `yaml:"template"`
	DataType TransformDataType `yaml:"data_type"`
}

type Transforms []*Transform

func bultinFuncs() template.FuncMap {
	return template.FuncMap{
		"title": cases.Title(language.Und).String,
		"httpGet": func(url string) (string, error) {
			body, err := httpclient.SendRequest(context.Background(), http.MethodGet, url, nil)
			if err != nil {
				return "", err
			}

			return string(body), nil
		},
	}
}

func compileRegExp(expr string) (re *regexp.Regexp, err error) {
	if isEmpty(expr) {
		return nil, nil
	}
	if re, err = regexp.Compile(expr); err != nil {
		return nil, err
	}
	return
}

func (t *Transform) Compile() (err error) {
	if t.DataType == "" {
		t.DataType = DataJSON
	}
	if !isEmpty(t.Template) {
		templ := template.New("")
		templ.Funcs(bultinFuncs())
		if t.template, err = templ.Parse(t.Template); err != nil {
			return fmt.Errorf("invalid template value: %w %w", err, ErrCompile)
		}
	}
	if t.regexp, err = compileRegExp(t.RegExp); err != nil {
		return fmt.Errorf("invalid regexp value: %w %w", err, ErrCompile)
	}
	return nil
}

func (t *Transform) Execute(data []byte) ([]byte, error) {
	if t.template == nil {
		return data, nil
	}
	if t.regexp != nil {
		found := t.regexp.FindAllStringSubmatch(string(data), -1)
		if len(found) == 0 {
			return data, nil
		}
		return t.applyTemplate(found[0])
	}
	if t.DataType == DataJSON {
		return t.applyTemplateByJson(data)
	}

	return t.applyTemplate(data)

}

func (t Transforms) Compile() error {
	for _, v := range t {
		if err := v.Compile(); err != nil {
			return err
		}
	}
	return nil
}

func (t Transforms) Execute(data []byte) (transformed []byte, err error) {
	transformed = data
	for _, v := range t {
		if transformed, err = v.Execute(transformed); err != nil {
			return nil, err
		}
	}
	return
}

// applyTemplateByJson render the template with JSON `data`.
func (t Transform) applyTemplateByJson(data []byte) ([]byte, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("transform: %w %w", err, ErrFailedExecution)
	}

	buf := bytes.NewBuffer(nil)
	if err := t.template.Execute(buf, m); err != nil {
		return nil, fmt.Errorf("transform: %w %w", err, ErrFailedExecution)
	}
	return buf.Bytes(), nil
}

// applyTemplate render the template with `data`.
func (t Transform) applyTemplate(data any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := t.template.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("transform: %w %w", err, ErrFailedExecution)
	}
	return buf.Bytes(), nil
}
