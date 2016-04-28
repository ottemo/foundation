package stripe

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "972198e1-2a7a-4bd7-9da9-d4d1be525dba", "can't obtain config")
	}

	newConfigs := env.ConfigList{
		env.ConfigItem{
			Config: env.StructConfigItem{
				Path:  ConstConfigPathGroup,
				Label: "Stripe",
				Type:  env.ConstConfigTypeGroup,
			},
			Validator: nil,
		},
		env.ConfigItem{
			Config: env.StructConfigItem{
				Path:   ConstConfigPathEnabled,
				Label:  "Enabled",
				Type:   env.ConstConfigTypeBoolean,
				Editor: "boolean",
			},
			Validator: func(value interface{}) (interface{}, error) {
				return utils.InterfaceToBool(value), nil
			},
		},
		env.ConfigItem{
			Config: env.StructConfigItem{
				Path:   ConstConfigPathName,
				Label:  "Name in checkout",
				Value:  "Credit Card",
				Type:   env.ConstConfigTypeVarchar,
				Editor: "line_text",
			},
			Validator: func(value interface{}) (interface{}, error) {
				if utils.CheckIsBlank(value) {
					return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "cfc4cb85-b769-414c-90fb-9be3fbe7fe98", "can't be blank")
				}
				return value, nil
			},
		},
		env.ConfigItem{
			Config: env.StructConfigItem{
				Path:        ConstConfigPathAPIKey,
				Label:       "API Key",
				Value:       "",
				Type:        env.ConstConfigTypeVarchar,
				Editor:      "line_text",
				Description: "Your API Key will be located in your Stripe Dashboard.",
			},
			Validator: nil,
		},
	}

	for _, newConfig := range newConfigs {
		err := config.RegisterItem(newConfig.Config, newConfig.Validator)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}
