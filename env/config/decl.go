// Package config is a default implementation of InterfaceConfig declared in
// "github.com/ottemo/foundation/env" package
package config

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCollectionNameConfig = "config"

	ConstErrorModule = "env/config"
	ConstErrorLevel  = env.ConstErrorLevelService
)

// DefaultConfig is a default implementer of InterfaceConfig
type DefaultConfig struct {
	configValues     map[string]interface{}
	configTypes      map[string]string
	configValidators map[string]env.FuncConfigValueValidator
}
