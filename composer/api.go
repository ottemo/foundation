package composer

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
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

	return nil
}

func composerUnit(context api.InterfaceApplicationContext) (interface{}, error) {
	var result map[string]interface{}

	if composer := GetComposer(); composer != nil {
		unit := composer.GetUnit( context.GetRequestArgument("unit") )
		if unit != nil {
			result = make(map[string]interface{})

			for _, item := range unit.ListItems() {
				result[item] = map[string]interface{}{
					"label": unit.GetLabel(item),
					"description": unit.GetDescription(item),
					"type": unit.GetType(item),
					"required": unit.IsRequired(item),
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
				result[unitName] = map[string]interface{} {
					"name": unit.GetName(),
					"label": unit.GetLabel(ConstPrefixUnit),
					"description": unit.GetLabel(ConstPrefixUnit),
					"in_type": unit.GetType(ConstPrefixUnit),
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
				result[unitName] = map[string]interface{} {
					"name": unit.GetName(),
					"label": unit.GetLabel(ConstPrefixUnit),
					"description": unit.GetLabel(ConstPrefixUnit),
					"in_type": unit.GetType(ConstPrefixUnit),
					"in_required": unit.IsRequired(ConstPrefixUnit),
				}
			}
		}
	}

	return result, nil
}

func composerInfo(context api.InterfaceApplicationContext) (interface{}, error) {

	result := map[string]interface{} {
		"item_prefix": map[string]interface{} {
			"unit": ConstPrefixUnit,
			"in": ConstPrefixArg,
			"out": "",
		},
	}

	if composer := GetComposer(); composer != nil {
		result["composer"] = composer.GetName()
		result["units_count"] = len(composer.ListUnits())
	}

	return result, nil
}