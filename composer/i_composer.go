package composer

import (
	"github.com/ottemo/foundation/env"
	"regexp"
	"strings"
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

func (it *DefaultComposer) Process(in interface{}, rules map[string]interface{}) (interface{}, error) {

	var result interface{}
	var err error

	// checking if input item is a Compose unit
	if unit, ok := in.(InterfaceComposeUnit); ok {
		unitIn := make(map[string]interface{})

		// gathering input for a unit
		for ruleKey, ruleValue := range rules {
			if strings.HasPrefix(ruleKey, ConstInPrefix) {
				key := strings.TrimPrefix(ruleKey, ConstInPrefix)
				unitIn[key] = ruleValue
				delete(unitIn, key)
			}

		}
		result, err = unit.Process(unitIn)
	}

	for ruleKey, ruleValue := range rules {
		currentIn := in

		// case 1 - {"$unit": ...}
		if strings.HasPrefix(ruleKey, ConstUnitPrefix) {
			if unit := it.GetUnit(strings.TrimPrefix(ruleKey, ConstUnitPrefix)); unit != nil {
				currentIn, err == unit.Process(currentIn)
				if err != nil {
					env.LogError(env.ErrorDispatch(err))
					result = false
				}
			}
		}

		// case 2 - {"item": {"$unit": ...}} or {"item": {"subItem": ...}}
		if ruleValue.(map[string]interface{}) {

		}

		// case 3 - {"item": []}
		// ...

		result = (ruleKey == ruleValue)
	}

	return result, err
}