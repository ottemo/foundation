package saleprice

import (
	"github.com/ottemo/foundation/app/models/discount/saleprice"
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
			Label:       "Sale Prices",
			Description: "Sale Prices system",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		validateEnabled := func(value interface{}) (interface{}, error) {
			boolValue := utils.InterfaceToBool(value)
			if boolValue {
				saleprice.RegisterSalePrice(new(DefaultSalePrice))

				productModel, err := product.GetProductModel()
				if err != nil {
					env.LogError(err)
				}

				if err = productModel.AddExternalAttributes(salePriceDelegate); err != nil {
					env.LogError(err)
				}

			} else {
				saleprice.UnRegisterSalePrice()

				productModel, err := product.GetProductModel()
				if err != nil {
					env.LogError(err)
				}

				if err = productModel.RemoveExternalAttributes(salePriceDelegate); err != nil {
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
			Description: "enables/disables sale price management",
			Image:       "",
		}, validateEnabled)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathSalePriceApplyPriority,
			Value:       1.10,
			Type:        env.ConstConfigTypeFloat,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Sale Price Priority",
			Description: "This value is used to determine when a sale price should be applied, (at Subtotal - 1, at Shipping - 2, at Grand total - 3)",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		validateEnabled(env.ConfigGetValue(ConstConfigPathEnabled))
	} else {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "de6ed851-1543-4b05-9058-fc098021578f", "Unable to obtain configuration for Sale Price")
		return env.ErrorDispatch(err)
	}

	return nil
}
