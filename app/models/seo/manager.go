package seo

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/seo"
)

// Package global variables
var (
	registeredSEOEngine InterfaceSEOEngine

	seoModels = map[string]string{}
	seoPaths = map[string]string{}
)

// RegisterSEOType registers SEO type association in system
func RegisterSEOType(seoType string, apiPath string, modelName string) error {

	_, present1 := seoPaths[seoType]
	_, present2 := seoPaths[seoType]
	if present1 || present2 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "99b01aae-e0cb-4d27-b0cf-406888828e31", "Already registered")
	}

	seoPaths[seoType] = apiPath
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
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "99b01aae-e0cb-4d27-b0cf-406888828e31", "Already registered")
	}
	registeredSEOEngine = seoEngine

	return nil
}

// GetRegisteredSEOEngine returns currently using SEO engine or nil
func GetRegisteredSEOEngine() InterfaceSEOEngine {
	return registeredSEOEngine
}
