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


func (it *DefaultComposer) GetName() string {
	return "DefaultComposer"
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

func (it *DefaultComposer) SearchUnits(namePattern string, typeFilter map[string]string) []InterfaceComposeUnit {
	var result []InterfaceComposeUnit

	if namePattern != "" && !(strings.HasPrefix(namePattern, "^") || strings.HasSuffix(namePattern, "$")) {
		namePattern = "^" + namePattern + "$"
	}

	namePattern = strings.Replace(namePattern, "%", ".*", -1)

	nameRegex, err := regexp.Compile(namePattern)
	if err != nil {
		env.LogError(env.ErrorDispatch(err))
		return result
	}

	for _, composerUnit := range it.units {
		if nameRegex.MatchString(composerUnit.GetName()) {

			ok := true
			for typeName, typeValue := range typeFilter {

				if matched, err := regexp.MatchString(typeValue, composerUnit.GetType(typeName)); err != nil || !matched {
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

func (it *DefaultComposer) Check(in interface{}, rule interface{}) (bool, error) {
	var result bool
	var err error
	var stop bool

	// fmt.Printf("s: %v <- %v\n", in, rule)

	// subroutine used within 2 places
	applyRuleToUnit := func(unit InterfaceComposeUnit, rule interface{}) (interface{}, error, bool) {
		args := make(map[string]interface{})

		useAsResult := false

		// looking for arguments addressed to CompositeUnit and not for unit process result
		if mapRule, ok := rule.(map[string]interface{}); ok {
			nonArgsKeys := 0
			for ruleKey, ruleValue := range mapRule {
				if strings.HasPrefix(ruleKey, ConstPrefixArg) {
					argKey := strings.TrimPrefix(ruleKey, ConstPrefixArg)
					args[argKey] = ruleValue
				} else {
					nonArgsKeys++
				}
			}

			if nonArgsKeys == 0 {
				useAsResult = true
			}
		} else {
			if unit.GetType(ConstPrefixArg) != "" {
				args[ConstPrefixArg] = rule
				useAsResult = true
			}
		}

		// processing unit with it's arguments
		output, outErr := unit.Process(in, args, it)
		if outErr != nil {
			return false, outErr, true
		}

		return output, outErr, useAsResult
	}

	// checking if in parameter is a ComposeUnit, then it should be processed
	if unit, ok := in.(InterfaceComposeUnit); ok {
		in, err, stop = applyRuleToUnit(unit, rule)
		if err != nil {
			return false, err
		}
		if stop {
			return utils.InterfaceToBool(in), err
		}
	}

	// unifying input (utils.InterfaceToArray(...) will do nothing if value already array)
	for _, ruleItem := range utils.InterfaceToArray(rule) {

		// case 1: in interface{} <- [...]
		if utils.IsArray(ruleItem) {

			result, err = it.Check(in, ruleItem)
			if err != nil {
				env.LogError(err)
				result = false
			}

		// case 2: in interface{} <- {...}
		} else if mapRule, ok := ruleItem.(map[string]interface{}); ok {

			for ruleKey, ruleValue := range mapRule {

				// case 2.1: in interface{} <- {"$unit": value}
				if strings.HasPrefix(ruleKey, ConstPrefixUnit) {
					if unit, present := it.units[strings.TrimPrefix(ruleKey, ConstPrefixUnit)]; present {
						if out, outErr, stop := applyRuleToUnit(unit, ruleValue); err == nil {
							if !stop {
								result, err = it.Check(out, ruleValue)
							} else {
								result = utils.InterfaceToBool(out)
							}
						} else {
							err = outErr
							result = false
						}
					} else {
						err = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3537c93c-4f22-466a-8c76-da47373a26ba", "unit not exists")
						result = false
					}

				// case 2.2: in map[string]interface{} <- {"key": value}
				} else if inAsMap, ok := in.(map[string]interface{}); ok {
					if inValue, present := inAsMap[ruleKey]; present {
						result, err = it.Check(inValue, ruleValue)
					} else {
						result = false
					}

				// case 2.3: in InterfaceObject <- {"key": value}
				} else if inAsObject, ok := in.(models.InterfaceObject); ok {
					result, err = it.Check(inAsObject.Get(ruleKey), ruleValue)

				// case 2.4: in interface{} <- {"key": value}
				} else {
					result = utils.Equals(in, ruleValue)
				}

				if err != nil { result = false }
				if !result { break }
			}

		// case 3: in interface{} <- interface{}
		} else {
			result = utils.Equals(in, rule)
		}

		if err != nil { result = false }
		if !result { break }
	}

	// fmt.Printf("e: %v <- %v = %v, %v\n", in, rule, result, err)

	return result, err
}
