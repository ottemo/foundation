package composer

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultComposer)
	instance.units = make(map[string]InterfaceComposeUnit)

	composer = instance

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

	composer.RegisterUnit( &BasicUnit{
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

	composer.RegisterUnit( &BasicUnit{
		Name: "gt",
		Type: map[string]string{
			ConstPrefixUnit: ConstTypeAny,
			ConstPrefixArg: ConstTypeAny,
			"": "bool",
		},
		Label: map[string]string{  ConstPrefixUnit: "greather" },
		Description: map[string]string{ ConstPrefixUnit: "Checks if value if greather then other value" },
		Action: action,
	})

}