package composer

import(
	"github.com/ottemo/foundation/env"
)

var (
	composer InterfaceComposer
)

// Package global constants
const (
	ConstPrefixUnit = "$"
	ConstPrefixArg = "@"

	ConstTypeAny = "*"
	ConstTypeValidate = "validate"

	ConstErrorModule = "composer"
	ConstErrorLevel  = env.ConstErrorLevelService
)


type FuncUnitAction func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error)
type FuncTypeValidator func(item string, inType string) bool

type DefaultComposer struct {
	units map[string]InterfaceComposeUnit
}

type BasicUnit struct {
	Name  string

	Value map[string]interface{}
	Type  map[string]string
	Label map[string]string
	Description map[string]string

	Validator FuncTypeValidator
	Action FuncUnitAction
}