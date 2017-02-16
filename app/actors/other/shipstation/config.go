package shipstation

import (
	"github.com/ottemo/foundation/env"
)

const (
	ConstConfigPathShipstation         = "general.shipstation"
	ConstConfigPathShipstationEnabled  = "general.shipstation.enabled"
	ConstConfigPathShipstationUsername = "general.shipstation.username"
	ConstConfigPathShipstationPassword = "general.shipstation.password"
)

func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "26abb8aa-253a-403e-bfc3-3852885373eb", "can't obtain config")
		return env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathShipstation,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Ship Station",
		Description: "",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathShipstationEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enabled",
		Description: "",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathShipstationUsername,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Username",
		Description: "",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathShipstationPassword,
		Value:       "",
		Type:        env.ConstConfigTypeSecret,
		Editor:      "password",
		Options:     "",
		Label:       "Password",
		Description: "",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
