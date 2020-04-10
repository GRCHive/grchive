package core

import (
	"bytes"
	htemplate "html/template"
	ttemplate "text/template"
)

func HtmlTemplateToString(tmpl *htemplate.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TextTemplateToString(tmpl *ttemplate.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
