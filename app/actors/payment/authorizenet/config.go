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
			Type:        env.ConstConfigTypeGroup,
			Editor:      "",
			Options:     nil,
			Label:       "Authorize.Net (Direct Post)",
			Description: "see https://developer.authorize.net/api/dpm/",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMEnabled,
			Value:       false,
			Type:        env.ConstConfigTypeBoolean,
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
			Path:        ConstConfigPathDPMTitle,
			Value:       "Authorize.Net (Direct Post)",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Title",
			Description: "payment method name in checkout",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			if utils.CheckIsBlank(value) {
				err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "e3e1002e-02da-45dd-876a-2ffd1d7790d2", "can't be blank")
				return nil, env.ErrorDispatch(err)
			}
			return value, nil
		})

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMAction,
			Value:       ConstDPMActionAuthorizeOnly,
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "select",
			Options:     map[string]string{ConstDPMActionAuthorizeOnly: "Authorize Only", ConstDPMActionAuthorizeAndCapture: "Authorize and Capture"},
			Label:       "Action",
			Description: "specifies action on checkout submit",
			Image:       "",
		}, func(value interface{}) (interface{}, error) {
			stringValue := utils.InterfaceToString(value)
			if !utils.IsAmongStr(stringValue, ConstDPMActionAuthorizeOnly, ConstDPMActionAuthorizeAndCapture) {
				err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "91645509-37f3-4add-b692-2dafac332bd1", "should be "+ConstDPMActionAuthorizeOnly+" or "+ConstDPMActionAuthorizeAndCapture)
				return nil, env.ErrorDispatch(err)
			}
			return value, nil
		})

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMLogin,
			Value:       "",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "line_text",
			Options:     nil,
			Label:       "API Login ID",
			Description: "account login id",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMKey,
			Value:       "",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Transaction Key",
			Description: "account transaction key",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMGateway,
			Value:       "https://test.authorize.net/gateway/transact.dll",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Gateway",
			Description: "payment method gateway",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMTest,
			Value:       false,
			Type:        env.ConstConfigTypeBoolean,
			Editor:      "boolean",
			Options:     nil,
			Label:       "Test Mode",
			Description: "specifies test mode for payment method",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMDebug,
			Value:       false,
			Type:        env.ConstConfigTypeBoolean,
			Editor:      "boolean",
			Options:     nil,
			Label:       "Debug",
			Description: "specifies to write log on that payment method",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMCheckout,
			Value:       false,
			Type:        env.ConstConfigTypeBoolean,
			Editor:      "boolean",
			Options:     nil,
			Label:       "Custom checkout page",
			Description: "use the custom relay page on checkout",
			Image:       "",
		}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMReceiptURL,
			Value:       "",
			Type:        env.ConstConfigTypeText,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Relay response receipt redirect URL",
			Description: "URL for approved transactions, default value: StorefrontURL/account/order/:orderID",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMDeclineURL,
			Value:       "",
			Type:        env.ConstConfigTypeText,
			Editor:      "line_text",
			Options:     nil,
			Label:       "Relay response decline redirect URL",
			Description: "URL for declined transactions, default value: StorefrontURL/checkout",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMReceiptHTML,
			Value:       ConstDefaultReceiptTemplate,
			Type:        env.ConstConfigTypeText,
			Editor:      "multiline_text",
			Options:     nil,
			Label:       "Relay response receipt HTML",
			Description: "",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathDPMDeclineHTML,
			Value:       ConstDefaultDeclineTemplate,
			Type:        env.ConstConfigTypeText,
			Editor:      "multiline_text",
			Options:     nil,
			Label:       "Relay response decline HTML",
			Description: "",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}
	} else {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "25367d9b-18e6-4304-99c7-4df2eae5521c", "Unable to obtain configuration for AuthorizeNet")
		return env.ErrorDispatch(err)
	}

	return nil
}
