package stock

import (
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	if config := env.GetConfig(); config != nil {
		err := config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathGroup,
			Value:       nil,
			Type:        env.ConstConfigTypeGroup,
			Editor:      "",
			Options:     nil,
			Label:       "Stock",
			Description: "Stock management system",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		validateEnabled := func(value interface{}) (interface{}, error) {
			boolValue := utils.InterfaceToBool(value)
			if boolValue {
				product.RegisterStock(new(DefaultStock))
			} else {
				product.UnRegisterStock()
			}
			return boolValue, nil
		}
		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathEnabled,
			Value:       false,
			Type:        env.ConstConfigTypeBoolean,
			Editor:      "boolean",
			Options:     nil,
			Label:       "Enabled",
			Description: "enables/disables stock management",
			Image:       "",
		}, validateEnabled)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		validateEnabled(env.ConfigGetValue(ConstConfigPathEnabled))
	}

	return nil
}
