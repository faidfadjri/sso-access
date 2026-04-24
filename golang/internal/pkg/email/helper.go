package email

import (
	"bytes"
	"text/template"
)

func ParseTemplate(path string, data interface{}) (string, error) {
	t, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}