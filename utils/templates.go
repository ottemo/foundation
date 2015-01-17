package utils

import (
	"bytes"
	"text/template"
)

var (
	templateFuncs = make(map[string]interface{})
)

// RegisterTemplateFunction registers custom function within text template processing
func RegisterTemplateFunction(name string, function interface{}) error {
	templateFuncs[name] = function
	return nil
}

// GetTemplateFunctions returns clone of templateFuncs (safe to manipulate)
func GetTemplateFunctions() map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range templateFuncs {
		result[key] = value
	}
	return result
}

// TextTemplate evaluates text template, returns error if not possible
func TextTemplate(templateContents string, context map[string]interface{}) (string, error) {

	textTemplate, err := template.New("TextTemplate").Funcs(templateFuncs).Parse(templateContents)
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
