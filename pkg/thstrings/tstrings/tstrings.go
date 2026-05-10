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
		"title": func(s string) string {
			return cases.Title(language.Und).String(s)
		},
		"toUpper":   strings.ToUpper,
		"toLower":   strings.ToLower,
		"trimSpace": strings.TrimSpace,
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
		"json": func(s string) any {
			var d any
			d = map[string]interface{}{}
			if err := json.Unmarshal([]byte(s), &d); err == nil {
				return d
			}
			d = []interface{}{}
			if err := json.Unmarshal([]byte(s), &d); err == nil {
				return d
			}
			return s
		},
		"match": regexp.MatchString,
		"ifMatch": func(pattern string, s string) string {
			if matched, err := regexp.MatchString(pattern, s); matched && err != nil {
				return s
			}
			return ""
		},
		"reReplaceAll": func(pattern, repl, text string) string {
			re := regexp.MustCompile(pattern)
			return re.ReplaceAllString(text, repl)
		},

		"reFindAll": func(pattern string, text string) []string {
			re := regexp.MustCompile(pattern)
			found := re.FindAllStringSubmatch(string(text), -1)
			if len(found) == 0 {
				return nil
			}
			return found[0]
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

func (ts *String) Execute(data any) ([]byte, error) {
	if ts.template == nil {
		return nil, thstrings.ErrEmptyTemplate
	}
	buf := bytes.NewBuffer(nil)
	if err := ts.template.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("execute: %w %w", err, thstrings.ErrFailedExecution)
	}
	return buf.Bytes(), nil

}

func (s *String) IsEmpty() bool {
	return s.Template == ""
}

func (ts *String) Match(data any) (bool, error) {
	if ts.template == nil {
		return true, nil
	}
	result, err := ts.Execute(data)
	return len(result) > 0, err
}
