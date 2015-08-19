package composer

import(
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "composer"
	ConstErrorLevel  = env.ConstErrorLevelService
)

type DefaultComposer struct {
	units map[string]InterfaceComposeUnit
}