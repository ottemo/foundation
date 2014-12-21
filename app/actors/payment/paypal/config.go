package paypal

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "681bd727d51e4bb6a2ba59855db7bee9", "can't obtain config")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGroup,
		Value:       nil,
		Type:        env.ConstConfigItemGroupType,
		Editor:      "",
		Options:     nil,
		Label:       "PayPal (Express)",
		Description: "see http://paypal.com",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathEnabled,
		Value:       false,
		Type:        "bool",
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enabled",
		Description: "enables/disables payment method in checkout",
		Image:       "",
	}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathTitle,
		Value:       "PayPal",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "Title",
		Description: "payment method name in checkout",
		Image:       "",
	}, func(value interface{}) (interface{}, error) {
		if utils.CheckIsBlank(value) {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "90107e36025949ffb74f3bb498fbed05", "can't be blank")
		}
		return value, nil
	})

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathNVP,
		Value:       "https://api-3t.sandbox.paypal.com/nvp",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "NVP Gateway",
		Description: "URL to send Name-Value Pair request",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGateway,
		Value:       "https://www.sandbox.paypal.com/webscr?cmd=_express-checkout",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "Gateway",
		Description: "PayPal gateway",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathUser,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "User",
		Description: "PayPal username",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathPass,
		Value:       "",
		Type:        "string",
		Editor:      "password",
		Options:     nil,
		Label:       "Password",
		Description: "PayPal password",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathSignature,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     nil,
		Label:       "Signature",
		Description: "PayPal signature",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathAction,
		Value:       "",
		Type:        "string",
		Editor:      "select",
		Options:     map[string]string{ConstPaymentActionSale: "Sale", ConstPaymentActionAuthorization: "Authorization"},
		Label:       "Signature",
		Description: "PayPal signature",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
