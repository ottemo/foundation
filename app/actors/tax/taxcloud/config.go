package taxcloud

import (
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "6645b15f-e6cb-4dc3-9656-ead549a73d3c", "can't obtain config")
		return env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "TaxCloud",
		Description: "TaxCloud",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathAPILoginID,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "ApiLoginID",
		Description: "ApiLoginID from TaxCloud Websites Area",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathAPIKey,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "ApiKey",
		Description: "ApiKey from TaxCloud Websites Area",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathAPIKey,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "ApiKey",
		Description: "ApiKey from TaxCloud Websites Area",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	validateEnabled := func(value interface{}) (interface{}, error) {
		boolValue := utils.InterfaceToBool(value)
		if boolValue {
			productModel, err := product.GetProductModel()
			if err != nil {
				env.LogError(err)
			}

			if err = productModel.AddExternalAttributes(ticDelegate); err != nil {
				env.LogError(err)
			}

		} else {
			productModel, err := product.GetProductModel()
			if err != nil {
				env.LogError(err)
			}

			if err = productModel.RemoveExternalAttributes(ticDelegate); err != nil {
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
		Description: "enables/disables TaxCloud integration",
		Image:       "",
	}, validateEnabled)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
