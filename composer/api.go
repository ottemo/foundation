package composer

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("composer/unit/:unit", api.ConstRESTOperationGet, composerUnit)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("composer/units", api.ConstRESTOperationGet, composerUnits)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("composer/units/:namePattern", api.ConstRESTOperationGet, composerUnitSearch)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("composer", api.ConstRESTOperationGet, composerInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("composer/test/:name", api.ConstRESTOperationGet, composerTest)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("composer/type/:name", api.ConstRESTOperationGet, composerType)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("composer/check", api.ConstRESTOperationCreate, composerCheck)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

func composerCheck(context api.InterfaceApplicationContext) (interface{}, error) {
	input := map[string]interface{}{
		"a": 10,
		"b": 25,
		"c": "123",
		"d": map[string]interface{}{
			"id":    "0bf7939973984c67b6b56b1c098edfca",
			"name":  "Product1",
			"sku":   "PR-1",
			"price": 10.5,
		},
	}

	rules, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result, err := GetComposer().Check(input, rules)

	if err != nil {
		return "Validation fail", err
	} else if !result {
		return "Validation fail", nil
	}

	return result, nil
}

func composerType(context api.InterfaceApplicationContext) (interface{}, error) {
	//	modelName := context.GetRequestArgument("name")
	//	model, err := models.GetModel(modelName)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	object, ok := model.(models.InterfaceObject)
	//	if !ok {
	//		env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "24bb5e98-b5de-4d4a-a5dc-cf2573dae3dd", "Model " + modelName + " is not implements InterfaceObject")
	//	}
	//
	//	result := make(map[string]interface{})
	//
	//	for _, value := range object.GetAttributesInfo() {
	//		result[value.Attribute] = value.Type;
	//	}

	composer := registeredComposer

	typeName := context.GetRequestArgument("name")
	typeInfo := composer.GetType(typeName)
	if typeInfo == nil {
		env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "24bb5e98-b5de-4d4a-a5dc-cf2573dae3dd", "Type "+typeName+" is not defined")
	}

	result := make(map[string]interface{})

	for _, item := range typeInfo.ListItems() {
		result[item] = map[string]interface{}{
			"label":       typeInfo.GetLabel(item),
			"description": typeInfo.GetDescription(item),
			"type":        typeInfo.GetType(item),
		}
	}

	return result, nil
}

func composerTest(context api.InterfaceApplicationContext) (interface{}, error) {
	result := make(map[string]interface{})
	typeName := utils.InterfaceToString(api.GetArgumentOrContentValue(context, "name"))
	switch strings.ToLower(typeName) {
	case "product":
		return map[string]interface{}{
			"test":    "test",
			"sku":     "array",
			"name":    "string",
			"price":   "float",
			"related": "[]Id",
			"options": "[]Option",
		}, nil

	case "option":
		return map[string]interface{}{
			"key":   "string",
			"value": "[]Any",
		}, nil
	case "id":
		return map[string]interface{}{
			"Id": "int",
		}, nil
	case "any":
		return "baseTypes", nil
	case "test":
		return map[string]interface{}{
			"a": "[]Test",
			"b": "int",
		}, nil
	case "all":
		return utils.DecodeJSONToStringKeyMap(`{
						"Product": {
                            "sku": "boolean",
                            "name": "string",
                            "price": "float",
                            "related": "[]Id",
                            "options": "[]Option"
                        },
                        "Option": {
                            "key": "string",
                            "value": "[]Test"
                        },
                        "Id" : "int",

                        "Test": {
                            "a": "[]Test",
                            "b": "array"
                        },

                        "Any" : "baseTypes"
                        }
                        `)
	}

	return result, nil
}

func composerUnit(context api.InterfaceApplicationContext) (interface{}, error) {
	var result map[string]interface{}

	if composer := GetComposer(); composer != nil {
		unit := composer.GetUnit(context.GetRequestArgument("unit"))
		if unit != nil {
			result = make(map[string]interface{})

			for _, item := range unit.ListItems() {
				result[item] = map[string]interface{}{
					"label":       unit.GetLabel(item),
					"description": unit.GetDescription(item),
					"type":        unit.GetType(item),
					"required":    unit.IsRequired(item),
				}
			}
		}
	}

	return result, nil
}

func composerUnits(context api.InterfaceApplicationContext) (interface{}, error) {

	result := make(map[string]interface{})

	if composer := GetComposer(); composer != nil {
		for _, unit := range composer.ListUnits() {
			if unitName := unit.GetName(); unitName != "" {
				result[unitName] = map[string]interface{}{
					"name":        unit.GetName(),
					"label":       unit.GetLabel(ConstPrefixUnit),
					"description": unit.GetLabel(ConstPrefixUnit),
					"in_type":     unit.GetType(ConstPrefixUnit),
					"in_required": unit.IsRequired(ConstPrefixUnit),
				}
			}
		}
	}

	return result, nil
}

func composerUnitSearch(context api.InterfaceApplicationContext) (interface{}, error) {

	result := make(map[string]interface{})

	namePattern := context.GetRequestArgument("namePattern")
	typeFilter := context.GetRequestArguments()
	if _, present := typeFilter["namePattern"]; present {
		delete(typeFilter, "namePattern")
	}

	if composer := GetComposer(); composer != nil {
		for _, unit := range composer.SearchUnits(namePattern, typeFilter) {
			if unitName := unit.GetName(); unitName != "" {
				result[unitName] = map[string]interface{}{
					"name":        unit.GetName(),
					"label":       unit.GetLabel(ConstPrefixUnit),
					"description": unit.GetLabel(ConstPrefixUnit),
					"in_type":     unit.GetType(ConstPrefixUnit),
					"in_required": unit.IsRequired(ConstPrefixUnit),
				}
			}
		}
	}

	return result, nil
}

func composerInfo(context api.InterfaceApplicationContext) (interface{}, error) {

	result := map[string]interface{}{
		"item_prefix": map[string]interface{}{
			"unit": ConstPrefixUnit,
			"in":   ConstPrefixArg,
			"out":  "",
		},
	}

	if composer := GetComposer(); composer != nil {
		result["composer"] = composer.GetName()
		result["units_count"] = len(composer.ListUnits())
	}

	return result, nil
}
