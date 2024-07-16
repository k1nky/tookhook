package entity

import (
	"bytes"
	"encoding/json"
	"html/template"
)

func ExecuteTemplateByJson(templ string, data []byte) ([]byte, error) {
	t := template.Must(template.New("").Parse(templ))
	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ExecuteTemplate(templ string, data any) ([]byte, error) {
	t := template.Must(template.New("").Parse(templ))

	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
