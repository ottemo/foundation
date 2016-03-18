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

	err = api.GetRestService().RegisterAPI("composer/units/:names", api.ConstRESTOperationGet, composerUnits)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("composer", api.ConstRESTOperationGet, composerInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("composer/go-json", api.ConstRESTOperationGet, composerGoTypes)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// get type info {types: {}, units:{}, types_units_binding:{}}
	err = api.GetRestService().RegisterAPI("composer/types/:names", api.ConstRESTOperationGet, composerTypes)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// add Check!!!!

	return nil
}

func composerTypes(context api.InterfaceApplicationContext) (interface{}, error) {

	result := make(map[string]interface{})
	typesResult := make(map[string]interface{})
	unitsResult := make(map[string]interface{})
	bindingResult := make(map[string]interface{})

	composer := GetComposer()
	baseForAny := map[string]int{"string":1, "int":1, "float":1, "boolean":1}
	typeNames := strings.Split(context.GetRequestArgument("names"), ",")
	listUnits := composer.ListUnits()

	for _, typeName := range typeNames {
		// types definition
		if typeInfo := composer.GetType(typeName); typeInfo != nil {
			keyInfo := make(map[string]interface{})
			for _, item := range typeInfo.ListItems() {
				keyInfo[item] = map[string]interface{}{
					"label": 		typeInfo.GetLabel(item),
					"description":  typeInfo.GetDescription(item),
					"type":  		typeInfo.GetType(item),
				}
			}

			typesResult[typeName] = keyInfo
		}

		// units definition
		var binding []string
		for _, unitInfo := range listUnits {
			unitType := unitInfo.GetType(ConstPrefixUnit)

			if unitType == typeName || (baseForAny[typeName] == 1 && unitType == "any") {

				unitName := unitInfo.GetName()
				// binding definition
				binding = append(binding, unitName);

				if unitsResult[unitName] == nil {
					keyInfo := make(map[string]interface{})
					for _, item := range unitInfo.ListItems() {
						keyInfo[item] = map[string]interface{}{
							"label":       unitInfo.GetLabel(item),
							"description": unitInfo.GetDescription(item),
							"type":        unitInfo.GetType(item),
							"required":    unitInfo.IsRequired(item),
						}
					}

					unitsResult[unitName] = keyInfo
				}
			}
		}
		bindingResult[typeName] = binding;
	}

	result["types"] = typesResult
	result["units"] = unitsResult
	result["type_unit_binding"] = bindingResult

	return result, nil
}

func composerGoTypes(context api.InterfaceApplicationContext) (interface{}, error) {

	result := make(map[string]interface{})
	composer := GetComposer()

	for _, goType := range []string{
		utils.ConstDataTypeID,
		utils.ConstDataTypeBoolean,
		utils.ConstDataTypeVarchar,
		utils.ConstDataTypeText,
		utils.ConstDataTypeDecimal,
		utils.ConstDataTypeMoney,
		utils.ConstDataTypeDatetime,
		utils.ConstDataTypeJSON,
	} {
		typeInfo := composer.GetType(goType)
		for _, item := range typeInfo.ListItems() {
			result[item] = map[string]interface{}{
				"label": typeInfo.GetLabel(item),
				"desc":  typeInfo.GetDescription(item),
				"type":  typeInfo.GetType(item),
			}
		}
	}

	return result, nil
}

func composerUnits(context api.InterfaceApplicationContext) (interface{}, error) {

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
