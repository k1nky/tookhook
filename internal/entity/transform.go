package entity

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"regexp"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TransformAction = string

const (
	TransformActionApply   TransformAction = ""
	TransformActionDiscard TransformAction = "discard"
)

type transform struct {
	on       *regexp.Regexp
	regexp   *regexp.Regexp
	template *template.Template
}

type Transform struct {
	transform
	Action TransformAction `yaml:"action"`
	// RegExp will be applied to the data before the template is executed.
	RegExp string `yaml:"regexp"`
	// Template will be applied to the data. The string must be formatted as a text/template package template.
	Template string `yaml:"template"`
	// On is a regexp, transformation will be applied if the regexp is matched.
	On string `yaml:"on"`
}

type Transforms []*Transform

func bultinFuncs() template.FuncMap {
	return template.FuncMap{
		"title": cases.Title(language.Und).String,
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
	if t.Action == TransformActionApply {
		if !isEmpty(t.Template) {
			templ := template.New("")
			templ.Funcs(bultinFuncs())
			if t.transform.template, err = templ.Parse(t.Template); err != nil {
				return fmt.Errorf("invalid template value: %w %w", err, ErrCompile)
			}
		}
		if t.transform.regexp, err = compileRegExp(t.RegExp); err != nil {
			return fmt.Errorf("invalid regexp value: %w %w", err, ErrCompile)
		}
	}
	if t.transform.on, err = compileRegExp(t.On); err != nil {
		return fmt.Errorf("invalid on value: %w %w", err, ErrCompile)
	}
	return nil
}

func (t Transform) Execute(data []byte) ([]byte, error) {
	if t.on != nil {
		if ok := t.on.Match(data); !ok {
			return data, fmt.Errorf("transform: %w", ErrNotMatch)
		}
	}
	if t.Action == TransformActionDiscard {
		return nil, fmt.Errorf("transform: %w", ErrDiscard)
	}
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
	return t.applyTemplateByJson(data)
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
	for _, t := range t {
		transformed, err = t.Execute(data)
		if errors.Is(err, ErrNotMatch) {
			continue
		}
		return
	}
	return
}
