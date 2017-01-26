package rts

import (
	"github.com/ottemo/foundation/env"
)

func setupConfig() error {

	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "6b78d38a-35c5-4aa2-aec1-eaa16830ff61", "Error configuring Emma module")
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
