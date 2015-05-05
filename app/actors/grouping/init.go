package grouping

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine before app start
func init() {
	app.OnAppStart(initListners)
	env.RegisterOnConfigStart(setupConfig)
}

// init Listeners for current model
func initListners() error {

	env.EventRegisterListener("api.cart.updatedCart", updateCartHandler)

	return nil
}
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "701e85e4-b63c-48f4-a990-673ba0ed6a2a", "can't obtain config")
	}

	// validateNewRules validate structure of new rules
	validateNewRules := func(newRulesValues interface{}) (interface{}, error) {

		newRules, err := utils.DecodeJSONToStringKeyMap(newRulesValues)
		if err != nil {
			return nil, err
		}

		// System error has occured will be if not Arrays in map, cant use default check method
		for _, groupInto := range newRules {
			groupInto := utils.InterfaceToArray(groupInto)

			for _, items := range groupInto {
				items := utils.InterfaceToArray(items)

				for _, item := range items {
					product, ok := item.(map[string]interface{})

					if !ok && product["pid"] == nil || product["qty"] == nil {
						return newRules, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0bc07b3d321", "need to scpecify qty pid in every product instance")
					}
				}
			}
		}

		return newRules, nil
	}

	// Grouping rules config setup
	//---------
	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstGroupingConfigPath,
		Value:       nil,
		Type:        env.ConstConfigTypeJSON,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Rules for grouping items",
		Description: `decribe products that will be grouped; example: { "group":[[{"pid":"pid1","qty":"n"},...],...], "into":[[{"options":{},"pid":"resultpid1","qty":"n"}],...]} `,
		Image:       "",
	}, env.FuncConfigValueValidator(validateNewRules))

	if err != nil {
		return env.ErrorDispatch(err)
	}
	return nil
}
