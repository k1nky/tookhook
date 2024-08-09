package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"
)

type transform struct {
	on       *regexp.Regexp
	regexp   *regexp.Regexp
	template *template.Template
}

type Transform struct {
	transform
	// RegExp will be applied to the data before the template is executed.
	RegExp string `yaml:"regexp"`
	// Template will be applied to the data. The string must be formatted as a text/template package template.
	Template string `yaml:"template"`
	// On is a regexp, transformation will be applied if the regexp is matched.
	On string `yaml:"on"`
}

type Transforms []*Transform

func compileRegExp(expr string) (*regexp.Regexp, error) {
	if isEmpty(expr) {
		return nil, nil
	}
	return regexp.Compile(expr)
}

func (t *Transform) Compile() (err error) {
	if isEmpty(t.Template) {
		return fmt.Errorf("template value: %w", ErrEmptyValue)
	}
	if t.transform.template, err = template.New("").Parse(t.Template); err != nil {
		return err
	}
	if t.transform.on, err = compileRegExp(t.On); err != nil {
		return err
	}
	if t.transform.regexp, err = compileRegExp(t.RegExp); err != nil {
		return err
	}
	return nil
}

func (t Transform) Execute(data []byte) ([]byte, error) {
	if t.on != nil {
		if ok := t.on.Match(data); !ok {
			return data, nil
		}
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
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if err := t.template.Execute(buf, m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// applyTemplate render the template with `data`.
func (t Transform) applyTemplate(data any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := t.template.Execute(buf, data); err != nil {
		return nil, err
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

func (t Transforms) Execute(data []byte) ([]byte, error) {
	for _, t := range t {
		return t.Execute(data)
	}
	return data, nil
}
