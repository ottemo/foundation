package authorizenet

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	if config := env.GetConfig(); config != nil {
		err := config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMGroup,
			Value:       nil,
			Type:        env.ConstConfigItemGroupType,
			Editor:      "",
			Options:     nil,
			Label:       "Authorize.Net (Direct Post)",
			Description: "see https://developer.authorize.net/api/dpm/",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMEnabled,
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

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMTitle,
			Value:       "Authorize.Net (Direct Post)",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "Title",
			Description: "payment method name in checkout",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			if utils.CheckIsBlank(value) {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "e3e1002e02da45dd876a2ffd1d7790d2", "can't be blank")
			}
			return value, nil
		})

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMAction,
			Value:       ConstDPMActionAuthorizeOnly,
			Type:        "string",
			Editor:      "select",
			Options:     map[string]string{ConstDPMActionAuthorizeOnly: "Authorize Only", ConstDPMActionAuthorizeAndCapture: "Authorize & Capture"},
			Label:       "Action",
			Description: "specifies action on checkout submit",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			stringValue := utils.InterfaceToString(value)
			if !utils.IsAmongStr(stringValue, ConstDPMActionAuthorizeOnly, ConstDPMActionAuthorizeAndCapture) {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "9164550937f34addb6922dafac332bd1", "should be "+ConstDPMActionAuthorizeOnly+" or "+ConstDPMActionAuthorizeAndCapture)
			}
			return value, nil
		})

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMLogin,
			Value:       "",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "API Login ID",
			Description: "account login id",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMKey,
			Value:       "",
			Type:        "string",
			Editor:      "password",
			Options:     nil,
			Label:       "Transaction Key",
			Description: "account transaction key",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMGateway,
			Value:       "https://test.authorize.net/gateway/transact.dll",
			Type:        "string",
			Editor:      "line_text",
			Options:     nil,
			Label:       "Gateway",
			Description: "payment method gateway",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMTest,
			Value:       false,
			Type:        "bool",
			Editor:      "boolean",
			Options:     nil,
			Label:       "Test Mode",
			Description: "specifies test mode for payment method",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMDebug,
			Value:       false,
			Type:        "bool",
			Editor:      "boolean",
			Options:     nil,
			Label:       "Debug",
			Description: "specifies to write log on that payment method",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMCheckout,
			Value:       false,
			Type:        "bool",
			Editor:      "boolean",
			Options:     nil,
			Label:       "Custom checkout page",
			Description: "use the custom relay page on checkout",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return err
		}
	}

	return nil
}
