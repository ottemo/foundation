package composer

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/utils"
	"strings"
	"regexp"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultComposer)
	instance.units = make(map[string]InterfaceComposeUnit)

	registeredComposer = instance

	api.RegisterOnRestServiceStart(setupAPI)
	initBaseUnits()
}


func initBaseUnits() {

	action := func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
		if argValue, present := args[ConstPrefixArg]; present {
			return utils.Equals(in, argValue), nil
		}
		return false, nil
	}

	registeredComposer.RegisterUnit( &BasicUnit{
		Name: "eq",
		Type: map[string]string{
			ConstPrefixUnit: ConstTypeAny, // input type
			ConstPrefixArg: ConstTypeAny,  // operand type (unnamed argument is a key for rule right-side value if it is not a map)
			"": "bool",                    // output type
		},
		Label: map[string]string{ ConstPrefixUnit: "equals" },
		Description: map[string]string{ ConstPrefixUnit: "Checks if value equals to other value" },
		Action: action,
	})


	action = func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
		if argValue, present := args[ConstPrefixArg]; present {
			if utils.InterfaceToFloat64(in) > utils.InterfaceToFloat64(argValue) {
				return true, nil
			}
		}
		return false, nil
	}

	registeredComposer.RegisterUnit( &BasicUnit{
		Name: "gt",
		Type: map[string]string{
			ConstPrefixUnit: ConstTypeAny,
			ConstPrefixArg: ConstTypeAny,
			"": "bool",
		},
		Label: map[string]string{ ConstPrefixUnit: ">" },
		Description: map[string]string{ ConstPrefixUnit: "Checks if value if greather then other value" },
		Action: action,
	})

	action = func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
		if argValue, present := args[ConstPrefixArg]; present {
			if utils.InterfaceToFloat64(in) > utils.InterfaceToFloat64(argValue) {
				return true, nil
			}
		}
		return false, nil
	}

	registeredComposer.RegisterUnit( &BasicUnit{
		Name: "lt",
		Type: map[string]string{
			ConstPrefixUnit: ConstTypeAny,
			ConstPrefixArg: ConstTypeAny,
			"": "bool",
		},
		Label: map[string]string{ ConstPrefixUnit: "<" },
		Description: map[string]string{ ConstPrefixUnit: "Checks if value if lower then other value" },
		Action: action,
	})

	action = func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
		if argValue, present := args[ConstPrefixArg]; present {
			if strings.Contains(utils.InterfaceToString(in), utils.InterfaceToString(argValue)) {
				return true, nil
			}
		}
		return false, nil
	}

	registeredComposer.RegisterUnit( &BasicUnit{
		Name: "contains",
		Type: map[string]string{
			ConstPrefixUnit: "string",
			ConstPrefixArg: "string",
			"": "bool",
		},
		Label: map[string]string{ ConstPrefixUnit: "contains" },
		Description: map[string]string{ ConstPrefixUnit: "Checks if value containt other value" },
		Action: action,
	})

	action = func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
		if argValue, present := args[ConstPrefixArg]; present {
			if matched, err := regexp.MatchString(utils.InterfaceToString(argValue), utils.InterfaceToString(in)); err == nil {
				return matched, nil
			}
		}
		return false, nil
	}

	registeredComposer.RegisterUnit( &BasicUnit{
		Name: "regex",
		Type: map[string]string{
			ConstPrefixUnit: "string",
			ConstPrefixArg: "string",
			"": "bool",
		},
		Label: map[string]string{ ConstPrefixUnit: "regex" },
		Description: map[string]string{ ConstPrefixUnit: "Checks regular expression over value" },
		Action: action,
	})
}