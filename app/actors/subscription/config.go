package subscription

import (
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "1f7ecfb8-b5e3-4361-b066-42c088f6b350", "can't obtain config")
	}

	// Trust pilot config elements
	//----------------------------

	err := config.RegisterItem(env.StructConfigItem{
		Path:        subscription.ConstConfigPathSubscription,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Subscription",
		Description: "Subscription settings",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        subscription.ConstConfigPathSubscriptionEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Subscription",
		Description: `Enabled Subscription`,
		Image:       "",
	}, env.FuncConfigValueValidator(func(newValue interface{}) (interface{}, error) {
		subscriptionEnabled = utils.InterfaceToBool(newValue)
		return newValue, nil
	}))

	if err != nil {
		return env.ErrorDispatch(err)
	}

	productsUpdate := func(newProductsValues interface{}) (interface{}, error) {

		// taking an array of product ids
		productsValue := utils.InterfaceToArray(newProductsValues)

		newProducts := make([]string, 0)
		for _, value := range productsValue {
			if productID := utils.InterfaceToString(value); productID != "" {
				newProducts = append(newProducts, productID)
			}
		}

		subscriptionProducts = newProducts
		return newProductsValues, nil
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        subscription.ConstConfigPathSubscriptionProducts,
		Value:       ``,
		Type:        env.ConstConfigTypeText,
		Editor:      "product_selector",
		Options:     nil,
		Label:       "Products",
		Description: `list of products that will be subscription`,
		Image:       "",
	}, env.FuncConfigValueValidator(productsUpdate))

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        subscription.ConstConfigPathSubscriptionEmailSubject,
		Value:       "Subscription",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "On fail subject",
		Description: `Email subject for emails`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path: subscription.ConstConfigPathSubscriptionEmailTemplate,
		Value: `Dear {{.Visitor.name}},
		Yours subscription can't be processed couse you have insufficient funds on Credit Card
		please create new subscription using valid credit card`,
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "On fail template",
		Description: "Email content to send to customers in case of failing on submit (insufficient funds on Credit Card, or some technical error)",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        subscription.ConstConfigPathSubscriptionStockEmailTemplate,
		Value:       `Items out of stock`,
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Submit template",
		Description: "contents of email that sented on out of stock error",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
