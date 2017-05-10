package taxcloud

import (
	"github.com/ottemo/foundation/env"
)

func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "502049b0-cf64-4d2f-9b52-b7c95e290282", "can't obtain config")
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

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enabled",
		Description: "enables/disables TaxCloud integration",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
