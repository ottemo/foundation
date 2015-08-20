package composer

import(
	"github.com/ottemo/foundation/env"
)

var (
	composer InterfaceComposer
)

// Package global constants
const (
	ConstInPrefix = "@"
	ConstInKey = ConstInPrefix

	ConstOutPrefix = ""
	ConstOutKey = ConstOutPrefix

	ConstUnitDescriptionKey = ConstOutPrefix

	ConstTypeAny = "*"

	ConstErrorModule = "composer"
	ConstErrorLevel  = env.ConstErrorLevelService
)

type FuncUnitAction func(in map[string]interface{}) (map[string]interface{}, error)

type DefaultComposer struct {
	units map[string]InterfaceComposeUnit
}

type BasicUnit struct {
	Name  string
	Type  map[string]string
	Label map[string]string
	Description map[string]string

	Action FuncUnitAction
}