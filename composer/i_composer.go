package composer

import (
	"github.com/ottemo/foundation/env"
	"regexp"
	"strings"
	"github.com/ottemo/foundation/utils"
)

func (it *DefaultComposer) RegisterUnit(unit InterfaceComposeUnit) error {
	unitName := unit.GetName()

	if _, present := it.units[unitName]; !present {
		it.units[unitName] = unit
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "57d471d0-1fe0-40bf-999d-96ef803f62fa", "unit already registered")
	}

	return nil
}

func (it *DefaultComposer) UnRegisterUnit(unit InterfaceComposeUnit) error {
	unitName := unit.GetName()

	if _, present := it.units[unitName]; present {
		delete(it.units, unitName)
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3537c93c-4f22-466a-8c76-da47373a26ba", "unit not exists")
	}

	return nil
}

func (it *DefaultComposer) GetUnit(name string) InterfaceComposeUnit {
	if unit, present := it.units[name]; present {
		return unit
	}
	return nil
}

func (it *DefaultComposer) ListUnits() []InterfaceComposeUnit {
	return it.SearchUnits("", nil)
}

func (it *DefaultComposer) SearchUnits(namePattern string, typeFilter map[string]interface{}) []InterfaceComposeUnit {
	var result []InterfaceComposeUnit

	if namePattern != "" && !(strings.HasPrefix(namePattern, "^") || strings.HasSuffix(namePattern, "$")) {
		namePattern = "^" + namePattern + "$"
	}

	nameRegex, err := regexp.Compile(namePattern)
	if err != nil {
		env.LogError(env.ErrorDispatch(err))
		return result
	}

	for _, composerUnit := range it.units {
		if nameRegex.MatchString(composerUnit.GetName()) {

			ok := true
			for typeName, typeValue := range typeFilter {
				if !regexp.MatchString(typeValue, composerUnit.GetType(typeName)) {
					ok = false
					break
				}
			}

			if ok {
				result = append(result, composerUnit)
			}
		}
	}

	return result
}

func (it *DefaultComposer) Process(in interface{}, rules map[string]interface{}) (bool, error) {

	var result interface{}
	var err error

	// checking if in parameter is a ComposeUnit, then it should be processed
	if unit, ok := in.(InterfaceComposeUnit); ok {
		unitIn := make(map[string]interface{})

		// looking for arguments addressed to CompositeUnit and not for unit process result
		for ruleKey, ruleValue := range rules {
			if strings.HasPrefix(ruleKey, ConstInPrefix) {
				key := strings.TrimPrefix(ruleKey, ConstInPrefix)
				unitIn[key] = ruleValue
				delete(rules, key)
			}

		}

		// processing unit with it's arguments
		in, err = unit.Process(unitIn)
		if err != nil {
			return nil, err
		}
	}

	// unifying {"item": value} to {"item": [value]} as particular case of {"item": [...]}
	for _, rule := range utils.InterfaceToArray(rules) {
		var result interface{}
		var err error

		if utils.IsArray(rule) {
			// case 1: in <- [...]
			result, err = it.Process(in, rule)
			if err != nil {
				env.LogError(err)
				result = false
			}

		} else if mapRule, ok := rule.(map[string]interface{}); ok {

			// case 2: in <- {"key": ...}
			for ruleKey, ruleValue := range mapRule {

				// case 2.1: in <- {"$unit": ...}
				if strings.HasPrefix(ruleKey, ConstUnitPrefix) {
					if unit := it.GetUnit(strings.TrimPrefix(ruleKey, ConstUnitPrefix)); unit != nil {
						result, err == unit.Process(ruleValue)
						if err != nil {
							env.LogError(err)
							result = false
						}
					}
				}

				// case 2 - {"item": {"$unit": ...}} or {"item": {"item": ...}}
				if ruleValueMap, ok := ruleValue.(map[string]interface{}); ok {
					result = it.Process(in, ruleValueMap)
				}
			}

		}


	}

	return result, err
}