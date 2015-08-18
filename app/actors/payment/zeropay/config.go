package zeropay

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "a1d81bf1-7d62-4107-a83e-be4d10439c5d", "can't obtain config")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Zero Payment",
		Description: "payment method for orders with no payment amount",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enabled",
		Description: "enables/disables Zero payment method in checkout",
		Image:       "",
	}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathName,
		Value:       "Zero Amount Payment",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Title",
		Description: "Zero Amount Payment method name in checkout",
		Image:       "",
	}, func(value interface{}) (interface{}, error) {
		if utils.CheckIsBlank(value) {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "cfc4cb85-b769-414c-90fb-9be3fbe7fe98", "can't be blank")
		}
		return value, nil
	})

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
