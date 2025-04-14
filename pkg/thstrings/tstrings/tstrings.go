package tstrings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/k1nky/tookhook/pkg/thstrings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type String struct {
	template *template.Template
	Template string
}

func bultinFuncs() template.FuncMap {
	return template.FuncMap{
		"title":     cases.Title(language.Und).String,
		"toUpper":   strings.ToUpper,
		"toLower":   strings.ToLower,
		"trimSpace": strings.TrimSpace,
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
		"match": regexp.MatchString,
		"reReplaceAll": func(pattern, repl, text string) string {
			re := regexp.MustCompile(pattern)
			return re.ReplaceAllString(text, repl)
		},
		"stringSlice": func(s ...string) []string {
			return s
		},
		"date": func(fmt string, t time.Time) string {
			return t.Format(fmt)
		},
		"tz": func(name string, t time.Time) (time.Time, error) {
			loc, err := time.LoadLocation(name)
			if err != nil {
				return time.Time{}, err
			}
			return t.In(loc), nil
		},
		"since": time.Since,
	}
}

func New(s string) *String {
	return &String{
		Template: s,
	}
}

func (ts String) MarshalYAML() (interface{}, error) {
	return ts.Template, nil
}

func (ts *String) UnmarshalYAML(unmarshal func(interface{}) error) error {

	return unmarshal(&ts.Template)
}

func (ts *String) Compile() (err error) {
	if !thstrings.IsEmpty(ts.Template) {
		templ := template.New("")
		templ.Funcs(bultinFuncs())
		if ts.template, err = templ.Parse(ts.Template); err != nil {
			return fmt.Errorf("invalid template value: %w %w", err, thstrings.ErrCompile)
		}
	}
	return nil
}

func (ts *String) unmarshalData(data []byte) (any, error) {
	var (
		d any = data
	)
	d = map[string]interface{}{}
	if err := json.Unmarshal(data, &d); err != nil {
		d = []interface{}{}
		if err := json.Unmarshal(data, &d); err != nil {
			return string(data), fmt.Errorf("execute: %w %w", err, thstrings.ErrFailedExecution)
		}
	}
	return d, nil
}

func (ts *String) Execute(data []byte) ([]byte, error) {
	if ts.template == nil {
		return data, nil
	}
	d, _ := ts.unmarshalData(data)
	buf := bytes.NewBuffer(nil)
	if err := ts.template.Execute(buf, d); err != nil {
		return nil, fmt.Errorf("execute: %w %w", err, thstrings.ErrFailedExecution)
	}
	return buf.Bytes(), nil

}

func (s *String) IsEmpty() bool {
	return s.Template == ""
}
