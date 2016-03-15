package composer

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/utils"
	"regexp"
	"strings"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultComposer)
	instance.units = make(map[string]InterfaceComposeUnit)
	instance.types = make(map[string]InterfaceComposeType)

	registeredComposer = instance

	api.RegisterOnRestServiceStart(setupAPI)
	app.OnAppStart(initModelTypes)
	initBaseUnits()
	initTest()
}

func initBaseUnits() {

	action := func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
		if argValue, present := args[ConstPrefixArg]; present {
			return utils.Equals(in, argValue), nil
		}
		return false, nil
	}

	registeredComposer.RegisterUnit(&BasicUnit{
		Name: "eq",
		Type: map[string]string{
			ConstPrefixUnit: ConstTypeAny, // input type
			ConstPrefixArg:  ConstTypeAny, // operand type (unnamed argument is a key for rule right-side value if it is not a map)
			"":              "bool",       // output type
		},
		Label:       map[string]string{ConstPrefixUnit: "equals"},
		Description: map[string]string{ConstPrefixUnit: "Checks if value equals to other value"},
		Action:      action,
	})

	action = func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
		if argValue, present := args[ConstPrefixArg]; present {
			if utils.InterfaceToFloat64(in) > utils.InterfaceToFloat64(argValue) {
				return true, nil
			}
		}
		return false, nil
	}

	registeredComposer.RegisterUnit(&BasicUnit{
		Name: "gt",
		Type: map[string]string{
			ConstPrefixUnit: ConstTypeAny,
			ConstPrefixArg:  ConstTypeAny,
			"":              "bool",
		},
		Label:       map[string]string{ConstPrefixUnit: ">"},
		Description: map[string]string{ConstPrefixUnit: "Checks if value if greather then other value"},
		Action:      action,
	})

	action = func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
		if argValue, present := args[ConstPrefixArg]; present {
			if utils.InterfaceToFloat64(in) > utils.InterfaceToFloat64(argValue) {
				return true, nil
			}
		}
		return false, nil
	}

	registeredComposer.RegisterUnit(&BasicUnit{
		Name: "lt",
		Type: map[string]string{
			ConstPrefixUnit: ConstTypeAny,
			ConstPrefixArg:  ConstTypeAny,
			"":              "bool",
		},
		Label:       map[string]string{ConstPrefixUnit: "<"},
		Description: map[string]string{ConstPrefixUnit: "Checks if value if lower then other value"},
		Action:      action,
	})

	action = func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
		if argValue, present := args[ConstPrefixArg]; present {
			if strings.Contains(utils.InterfaceToString(in), utils.InterfaceToString(argValue)) {
				return true, nil
			}
		}
		return false, nil
	}

	registeredComposer.RegisterUnit(&BasicUnit{
		Name: "contains",
		Type: map[string]string{
			ConstPrefixUnit: "string",
			ConstPrefixArg:  "string",
			"":              "bool",
		},
		Label:       map[string]string{ConstPrefixUnit: "contains"},
		Description: map[string]string{ConstPrefixUnit: "Checks if value containt other value"},
		Action:      action,
	})

	action = func(in interface{}, args map[string]interface{}, composer InterfaceComposer) (interface{}, error) {
		if argValue, present := args[ConstPrefixArg]; present {
			if matched, err := regexp.MatchString(utils.InterfaceToString(argValue), utils.InterfaceToString(in)); err == nil {
				return matched, nil
			}
		}
		return false, nil
	}

	registeredComposer.RegisterUnit(&BasicUnit{
		Name: "regex",
		Type: map[string]string{
			ConstPrefixUnit: "string",
			ConstPrefixArg:  "string",
			"":              "bool",
		},
		Label:       map[string]string{ConstPrefixUnit: "regex"},
		Description: map[string]string{ConstPrefixUnit: "Checks regular expression over value"},
		Action:      action,
	})
}

func initTest() error {
	testType := &BasicType{
		Name: "Test",
		Label: map[string]string{
			"a": "IntegerTest",
			"b": "FloatTest",
			"c": "StringTest",
			"d": "ProductTest",
		},
		Type: map[string]string{
			"a": "int",
			"b": "float",
			"c": "string",
			"d": "Product",
		},
		Description: map[string]string{
			"a": "Description for Test type",
		},
	}

	registeredComposer.RegisterType(testType)
	return nil
}

func initModelTypes() error {
	for modelName, modelInstance := range models.GetDeclaredModels() {
		if modelInstance == nil {
			continue
		}

		modelInstance, err := modelInstance.New()
		if err != nil || modelInstance == nil {
			continue
		}

		if objectInstance, ok := modelInstance.(models.InterfaceObject); ok {
			productType := &BasicType{
				Name:        modelName,
				Label:       make(map[string]string),
				Type:        make(map[string]string),
				Description: make(map[string]string),
			}

			for _, v := range objectInstance.GetAttributesInfo() {
				productType.Label[v.Attribute] = v.Label
				productType.Type[v.Attribute] = v.Type
				productType.Description[v.Attribute] = "Product field " + v.Label
			}

			registeredComposer.RegisterType(productType)
		}
	}

	return nil
}
