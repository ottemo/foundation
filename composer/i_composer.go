package composer

import (
	"github.com/ottemo/foundation/env"
	"regexp"
	"strings"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/app/models"
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

				if matched, err := regexp.MatchString(utils.InterfaceToString(typeValue), composerUnit.GetType(typeName)); err != nil || !matched {
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

func (it *DefaultComposer) Validate(in interface{}, rule interface{}) (bool, error) {
	var result bool
	var err error

	// checking if in parameter is a ComposeUnit, then it should be processed
	if unit, ok := in.(InterfaceComposeUnit); ok {
		unitIn := make(map[string]interface{})

		// looking for arguments addressed to CompositeUnit and not for unit process result
		if mapRule, ok := rule.(map[string]interface{}); ok {
			for ruleKey, ruleValue := range mapRule {
				if strings.HasPrefix(ruleKey, ConstInPrefix) {
					key := strings.TrimPrefix(ruleKey, ConstInPrefix)
					unitIn[key] = ruleValue
					delete(mapRule, key)
				}

			}
		}

		// processing unit with it's arguments
		in, err = unit.Process(MakeComposeValue(unitIn), it)
		if err != nil {
			return false, err
		}
	}

	// unifying input (utils.InterfaceToArray(...) will do nothing if value already array)
	for _, ruleItem := range utils.InterfaceToArray(rule) {
		if utils.IsArray(ruleItem) {

			// case 1: in <- [...]
			result, err = it.Validate(in, ruleItem)
			if err != nil {
				env.LogError(err)
				result = false
			}

		} else if mapRule, ok := ruleItem.(map[string]interface{}); ok {

			// case 2: in <- {"key": ...}
			for ruleKey, ruleValue := range mapRule {

				// case 2.1: in <- {"$unit": value}
				if strings.HasPrefix(ruleKey, ConstUnitPrefix) {
					it.Validate(in, ruleValue)
				}

				// case 2.2: in <- {"key": value}
				if inAsMap, ok := in.(map[string]interface{}); ok {
					if inValue, present := inAsMap[ruleKey]; present {
						result = (inValue == ruleValue)
					} else {
						result = false
					}

				} else if inAsObject, ok := ruleValue.(models.InterfaceObject); ok {
					result = (inAsObject.Get(ruleKey) == ruleValue)
				} else {
					result = (in == ruleValue)
				}
			}

			if !result { break }

		} else {
			result = (in == rule)
		}

		if !result { break }
	}

	return result, err
}
