package core

import (
	"bytes"
	"html/template"
)

func TemplateToString(tmpl *template.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
