package subscription

import (
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "1f7ecfb8-b5e3-4361-b066-42c088f6b350", "can't obtain config")
		return env.ErrorDispatch(err)
	}

	// Subscription config elements
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
		Label:       "Enable Subscriptions",
		Description: "",
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
		Type:        db.TypeArrayOf(db.ConstTypeID),
		Editor:      "product_selector",
		Options:     nil,
		Label:       "Applicable Products",
		Description: "",
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
		Label:       "Insufficient Funds Email: Subject",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path: subscription.ConstConfigPathSubscriptionEmailTemplate,
		Value: `Dear {{.Visitor.name}},
Yours subscription can't be processed because of an insufficient funds error from your credit card
provider, please create new subscription using valid credit card.`,
		Type:        env.ConstConfigTypeHTML,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Insufficient Funds Email: Body",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Cancellation email subject
	err = config.RegisterItem(env.StructConfigItem{
		Path:        subscription.ConstConfigPathSubscriptionCancelEmailSubject,
		Value:       "Subscription cancellation",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Cancellation Email: Subject",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Cancellation email body
	err = config.RegisterItem(env.StructConfigItem{
		Path:        subscription.ConstConfigPathSubscriptionCancelEmailTemplate,
		Value:       "Your subscription has been canceled",
		Type:        env.ConstConfigTypeHTML,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Cancellation Email: Body",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        subscription.ConstConfigPathSubscriptionStockEmailTemplate,
		Value:       "Subscription failure due to out of stock items.",
		Type:        env.ConstConfigTypeHTML,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Admin - Stock Warning Email: Body",
		Description: "",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
