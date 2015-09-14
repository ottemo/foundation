package composer

import(
	"github.com/ottemo/foundation/env"
)

var (
	composer InterfaceComposer
)

// Package global constants
const (
	ConstUnitPrefix = "$"

	ConstInPrefix = "@"
	ConstInItem = ConstInPrefix

	ConstOutPrefix = ""
	ConstOutItem = ConstOutPrefix

	ConstUnitLabelItem = ConstOutPrefix
	ConstUnitDescriptionItem = ConstOutPrefix

	ConstTypeAny = "*"
	ConstTypeValidate = "validate"

	ConstErrorModule = "composer"
	ConstErrorLevel  = env.ConstErrorLevelService
)


type FuncUnitAction func(in map[string]interface{}, composer InterfaceComposer) (map[string]interface{}, error)
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



func MakeComposeValue(inValue interface{}) map[string] interface{} {
	if typedValue, ok := inValue.(map[string]interface{}); !ok {
		return typedValue
	}
	return map[string]interface{} { ConstInItem: inValue }
}

func MakeInKey(name string) string {
	return ConstInPrefix + name
}

func MakeOutKey(name string) string {
	return ConstOutPrefix + name
}