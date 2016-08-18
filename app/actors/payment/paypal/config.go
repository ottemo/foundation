package paypal

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "681bd727-d51e-4bb6-a2ba-59855db7bee9", "can't obtain config")
		return env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
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
		Path:        ConstConfigPathTitle,
		Value:       "PayPal",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Title",
		Description: "payment method name in checkout",
		Image:       "",
	}, func(value interface{}) (interface{}, error) {
		if utils.CheckIsBlank(value) {
			err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "90107e36-0259-49ff-b74f-3bb498fbed05", "can't be blank")
			return nil, env.ErrorDispatch(err)
		}
		return value, nil
	})

	if err != nil {
 		return env.ErrorDispatch(err)
	}

	// choose PayPal gateway according to the workflow mode
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathPayPalExpressGateway,
		Value:       ConstPaymentPayPalGatewaySandbox,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     map[string]string{
			ConstPaymentPayPalGatewaySandbox:	"Sandbox",
			ConstPaymentPayPalGatewayProduction:	"Production"},
		Label:       "Gateway mode",
		Description: "Change PayPal gateway according to the workflow mode",
		Image:       "",
	    }, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathUser,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
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
		Type:        env.ConstConfigTypeSecret,
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
		Type:        env.ConstConfigTypeVarchar,
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
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     map[string]string{ConstPaymentActionSale: "Sale", ConstPaymentActionAuthorization: "Authorization"},
		Label:       "Payment action",
		Description: "PayPal payment action, how you want to obtain payment.",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// register config values for PayPal Pro API
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathPayPalPayflowGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "PayPal Pro",
		Description: "see http://paypal.com",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathPayPalPayflowEnabled,
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
		Path:        ConstConfigPathPayPalPayflowTitle,
		Value:       "PayPal Pro",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Title",
		Description: "payment method name in checkout",
		Image:       "",
	}, func(value interface{}) (interface{}, error) {
		if utils.CheckIsBlank(value) {
			err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "90107e36-0259-49ff-b74f-3bb498fbed05", "can't be blank")
			return nil, env.ErrorDispatch(err)
		}
		return value, nil
	})

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathPayPalPayflowUser,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "User",
		Description: "Payflow user",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathPayPalPayflowVendor,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Vendor",
		Description: "Payflow unique, case sensitive merchant login",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathPayPalPayflowPass,
		Value:       "",
		Type:        env.ConstConfigTypeSecret,
		Editor:      "password",
		Options:     nil,
		Label:       "Password",
		Description: "Payflow password",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// choose host according to the workflow mode
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathPayPalPayFlowGateway,
		Value:       ConstPaymentPayPalGatewaySandbox,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     map[string]string{
			ConstPaymentPayPalGatewaySandbox:	"Sandbox",
			ConstPaymentPayPalGatewayProduction:	"Production"},
		Label:       "Gateway mode",
		Description: "Change PayPal gateway according to the workflow mode",
		Image:       "",
	    }, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathPayPalPayflowTokenable,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Credit card saving",
		Description: "enables/disables saving of credit card to visitor",
		Image:       "",
	}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
