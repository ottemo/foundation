package flatrate

import (
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	if config := env.GetConfig(); config != nil {
		err := config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathGroup,
			Value:       nil,
			Type:        env.ConstConfigItemGroupType,
			Editor:      "",
			Options:     nil,
			Label:       "Flat Rate",
			Description: "static amount stipping method",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathEnabled,
			Value:       false,
			Type:        "bool",
			Editor:      "boolean",
			Options:     nil,
			Label:       "Enabled",
			Description: "enables/disables shipping method for storefront",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathAmount,
			Value:       10,
			Type:        "int",
			Editor:      "money",
			Options:     nil,
			Label:       "Amount",
			Description: "price of shipping",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathName,
			Value:       "Flat Rate",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "Name",
			Description: "shipping name displayed in checkout",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDays,
			Value:       0,
			Type:        "int",
			Editor:      "integer",
			Options:     nil,
			Label:       "Ship days",
			Description: "amount of days for shipping",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}
