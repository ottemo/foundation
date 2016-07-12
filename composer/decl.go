package composer

import (
	"github.com/ottemo/foundation/env"
)

var (
	registeredComposer InterfaceComposer
)

// Package global constants
const (
	ConstPrefixUnit = "*"
	ConstPrefixArg  = "@"
	ConstPrefixOut  = "#"

	ConstTypeAny      = "any"
	ConstTypeValidate = "validate"

	ConstErrorModule = "composer"
	ConstErrorLevel  = env.ConstErrorLevelService
)

type FuncUnitAction func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error)
type FuncTypeValidator func(item string, inType string) bool

type DefaultComposer struct {
	units map[string]InterfaceComposeUnit
	types map[string]InterfaceComposeType
}

type BasicUnit struct {
	Name string

	Value       map[string]interface{}
	Type        map[string]string
	Label       map[string]string
	Description map[string]string

	Required map[string]bool

	Validator FuncTypeValidator
	Action    FuncUnitAction
}

type BasicType struct {
	Name string

	Type        map[string]string
	Label       map[string]string
	Description map[string]string
}
