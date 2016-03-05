package shareasale

import "github.com/ottemo/foundation/env"

func setupConfig() error {

	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "6b78d38a-35c5-4aa2-aec1-eaa16830ff61", "Error configuring ShareASale module")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathShareASale,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "ShareASale",
		Description: "Share A Sale Settings",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathShareASaleEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Share A Sale Module Enabled",
		Description: "Enable Share A Sale integration(defaults to false)",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathShareASaleMerchantID,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Share A Sale Merchant ID",
		Description: "Enter your Merchant ID",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
