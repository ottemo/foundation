package subscription

import (
	"github.com/ottemo/foundation/env"
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
		Path:        ConstConfigPathSubscription,
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
		Path:        ConstConfigPathSubscriptionEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Subscription",
		Description: `Enabled Subscription`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathSubscriptionEmailSubject,
		Value:       "Subscription confirm!",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Confirmation subject",
		Description: `Email subject for confirmation emails`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathSubscriptionConfirmationLink,
		Value:       "subscription?confirm={{subscriptionID}}",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Confirmation link",
		Description: `Part of confirmation link to storefront to procced subscription confirmation {{subscriptionID}} - will cahnged by it's real id`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path: ConstConfigPathSubscriptionEmailTemplate,
		Value: `Dear {{.Visitor.name}},
		Yours subscription order is can be subbmitted, confirm processing of subscription order
		<br />
		<br />
<h3>Duplicated order #{{.Order.increment_id}}: </h3><br />
Order summary<br />
Subtotal: ${{.Order.subtotal}}<br />
Tax: ${{.Order.tax_amount}}<br />
Shipping: ${{.Order.shipping_amount}}<br />
Discount: -${{ .Order.discount | printf "%.2f" }}<br />
Total: ${{.Order.grand_total}}<br />
<br />
<a href={{.Info.link}}>Confirm</a><br />`,
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Confirmation template",
		Description: "contents of confirmation email, it will be sented to cutomers one week before subscription date",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathSubscriptionSubmitEmailSubject,
		Value:       "Subscription proceed!",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Submit subject",
		Description: `Email subject for checkout proceed emails`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathSubscriptionSubmitEmailLink,
		Value:       "subscription?submit={{subscriptionID}}",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Submit link",
		Description: `Part of confirmation link to storefront to proceed checkout for subscription {{subscriptionID}} - will cahnged by it's real id`,
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path: ConstConfigPathSubscriptionSubmitEmailTemplate,
		Value: `Dear {{.Visitor.name}},
		Yours subscription date is comming soon please confirm processing of subscription order
		<br />
		<br />
<h3>Duplicated order #{{.Order.increment_id}}: </h3><br />
Order summary<br />
Subtotal: ${{.Order.subtotal}}<br />
Tax: ${{.Order.tax_amount}}<br />
Shipping: ${{.Order.shipping_amount}}<br />
Discount: -${{ .Order.discount | printf "%.2f" }}<br />
Total: ${{.Order.grand_total}}<br />
<br />
<a href={{.Info.link}}>Confirm</a><br />`,
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Submit template",
		Description: "contents of confirmation email, it will be sented to cutomers one week before subscription date",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
