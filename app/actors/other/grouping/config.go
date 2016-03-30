package grouping

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "b2c1c442-36b9-4994-b5d1-7c948a7552bd", "can't obtain config")
	}

	// validateNewRules validate structure of new rules
	validateNewRules := func(newRulesValues interface{}) (interface{}, error) {

		// taking rules as array
		if newRulesValues != "" {
			rules, err := utils.DecodeJSONToArray(newRulesValues)
			if err != nil {
				return nil, err
			}

			// checking rules array
			for _, rule := range rules {
				ruleItem := utils.InterfaceToMap(rule)

				if !utils.KeysInMapAndNotBlank(ruleItem, "group", "into") {
					return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7912df05-8ea7-451e-83bd-78e9e201378e", "keys 'group' and 'into' should be not null")
				}

				// checking product specification arrays
				for _, groupingValue := range []interface{}{ruleItem["group"], ruleItem["into"]} {
					groupingElement := utils.InterfaceToArray(groupingValue)

					for _, productValue := range groupingElement {
						productElement := utils.InterfaceToMap(productValue)

						if !utils.KeysInMapAndNotBlank(productElement, "pid", "qty") {
							return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6b9deedd-39d1-46b0-9157-9b8d96bda858", "keys 'qty' and 'pid' should be not null")
						}
					}
				}
			}

			currentRules = rules
		}

		return newRulesValues, nil
	}

	// grouping rules config setup
	//----------------------------
	err := config.RegisterItem(env.StructConfigItem{
		Path:    ConstGroupingConfigPath,
		Value:   ``,
		Type:    env.ConstConfigTypeJSON,
		Editor:  "multiline_text",
		Options: "",
		Label:   "Rules for grouping items",
		Description: `Rules must be in JSON format:
[
	{
		"group": [{ "pid": "id1", "qty": 1 }, ...],
		"into":  [{ "pid": "id2", "qty": 1, "options": {"color": "red"}, ...]
	}, ...
]`,
		Image: "",
	}, env.FuncConfigValueValidator(validateNewRules))

	if err != nil {
		return env.ErrorDispatch(err)
	}
	return nil
}
