package utils

import (
	"bytes"
	"text/template"
)

// TextTemplate evaluates text template, returns error if not possible
func TextTemplate(templateContents string, context map[string]interface{}) (string, error) {

	textTemplate := template.New("textTemplate")
	textTemplate, err := textTemplate.Parse(templateContents)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	err = textTemplate.Execute(&result, context)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
