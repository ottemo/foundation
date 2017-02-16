package seo

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	registeredSEOEngine InterfaceSEOEngine

	seoModels   = map[string]string{}
	seoAPIPaths = map[string]string{}
)

// GetSEOTypeModel returns model name for a given seo type
func GetSEOTypeModel(seoType string) string {
	if value, present := seoModels[seoType]; present {
		return value
	}
	return ""
}

// GetSEOTypeAPIPath returns API path for a given seo type
func GetSEOTypeAPIPath(seoType string) string {
	if value, present := seoAPIPaths[seoType]; present {
		return value
	}
	return ""
}

// IsSEOType checks existence of given seo type
func IsSEOType(seoType string) bool {
	_, present1 := seoAPIPaths[seoType]
	_, present2 := seoModels[seoType]

	if present1 && present2 {
		return true
	}
	return false
}

// RegisterSEOType registers SEO type association in system
func RegisterSEOType(seoType string, apiPath string, modelName string) error {

	_, present1 := seoAPIPaths[seoType]
	_, present2 := seoModels[seoType]

	if present1 || present2 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "99b01aae-e0cb-4d27-b0cf-406888828e31", "Already registered")
	}

	seoAPIPaths[seoType] = apiPath
	seoModels[seoType] = modelName

	return nil
}

// UnRegisterSEOEngine removes currently using SEO engine from system
func UnRegisterSEOEngine() error {
	registeredSEOEngine = nil
	return nil
}

// RegisterSEOEngine registers given SEO engine in system
func RegisterSEOEngine(seoEngine InterfaceSEOEngine) error {
	if registeredSEOEngine != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a953a440-1c05-416a-8412-e01454e81c17", "Already registered")
	}
	registeredSEOEngine = seoEngine

	return nil
}

// GetRegisteredSEOEngine returns currently using SEO engine or nil
func GetRegisteredSEOEngine() InterfaceSEOEngine {
	return registeredSEOEngine
}
