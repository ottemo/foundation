package rts

import (
	"github.com/ottemo/foundation/env"
)

func setupConfig() error {

	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "a0c79b8f-5782-40bf-bae9-f0108e38d344", "Error configuring rts module")
		return env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathCheckoutPath,
		Value:       "/checkout",
		Type:        env.ConstConfigTypeText,
		Editor:      "",
		Options:     nil,
		Label:       "Checkout page path",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
