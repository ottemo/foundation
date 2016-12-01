package braintree

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "07fe3f67-d1d5-43e7-ace4-ce123c7f820d", "can't obtain config")
		return env.ErrorDispatch(err)
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:  ConstConfigPathGroup,
		Label: "Braintree",
		Type:  env.ConstConfigTypeGroup,
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:   ConstConfigPathEnabled,
		Label:  "Enabled",
		Type:   env.ConstConfigTypeBoolean,
		Editor: "boolean",
	}, func(value interface{}) (interface{}, error) {
		return utils.InterfaceToBool(value), nil
	})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:   ConstConfigPathName,
		Label:  "Name in checkout",
		Value:  ConstPaymentInternalName,
		Type:   env.ConstConfigTypeVarchar,
		Editor: "line_text",
	}, func(value interface{}) (interface{}, error) {
		if utils.CheckIsBlank(value) {
			err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "1bde8c7e-4b16-4f9d-9808-a7d520dcbc60", "can't be blank")
			return nil, env.ErrorDispatch(err)
		}
		return value, nil
	})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	//err = config.RegisterItem(env.StructConfigItem{
	//	Path:        ConstConfigPathAPIKey,
	//	Label:       "API Key",
	//	Value:       "",
	//	Type:        env.ConstConfigTypeVarchar,
	//	Editor:      "line_text",
	//	Description: "Your API Key will be located in your Stripe Dashboard.",
	//}, nil)
	//if err != nil {
	//	return env.ErrorDispatch(err)
	//}

	return nil
}

//// ConfigIsEnabled is a flag to enable/disable this payment module
//func (it Payment) ConfigIsEnabled() bool {
//	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEnabled))
//}
//
//// ConfigAPIKey is a method that returns the API Key from the db
//func (it Payment) ConfigAPIKey() string {
//	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAPIKey))
//}
//
//// ConfigNameInCheckout is a method that returns the payment method name to be used in checkout
//func (it Payment) ConfigNameInCheckout() string {
//	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathName))
//}

