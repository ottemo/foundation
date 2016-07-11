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

				productModel, err := product.GetProductModel()
				if err != nil {
					env.LogError(err)
				}

				if err = productModel.AddExternalAttributes(stockDelegate); err != nil {
					env.LogError(err)
				}

			} else {
				product.UnRegisterStock()

				productModel, err := product.GetProductModel()
				if err != nil {
					env.LogError(err)
				}

				if err = productModel.RemoveExternalAttributes(stockDelegate); err != nil {
					env.LogError(err)
				}
			}
			return boolValue, nil
		}

		err = config.RegisterItem(env.StructConfigItem{
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
	} else {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "fdc4f498-3d03-48a9-b51b-46aeae42edd1", "Unable to obtain configuration for Stock")
		return env.ErrorDispatch(err)
	}

	return nil
}
