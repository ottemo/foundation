package giftcard

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "15859fac-8fc0-4fbf-a801-b9cacf70d356", "Unable to obtain configuration for Gift Cards")
		return env.ErrorDispatch(err)
	}

	// giftCardSkuElement

	giftCardSkuElementSet := func(value interface{}) (interface{}, error) {
		newValue := utils.InterfaceToString(value)
		if newValue == "" {
			err := env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0a625584-1f91-4416-86bb-b25d6b37c70d", "can't be empty string")
			return value, env.ErrorDispatch(err)
		}
		checkout.GiftCardSkuElement = newValue
		return newValue, nil
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGiftCardSKU,
		Value:       "gift-card",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "GiftCard SKU Identifier",
		Description: "This value represents the product SKU for GiftCards and will act as a flag for gift card operations",
		Image:       "",
	}, giftCardSkuElementSet)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGiftEmailSubject,
		Value:       "Your GiftCard has Arrived",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Email Subject",
		Description: "This value will appear in the recipient email subject line",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGiftCardApplyPriority,
		Value:       3.10,
		Type:        env.ConstConfigTypeFloat,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Gift Card Priority",
		Description: "This value is used to determine when a gift card should be applied, (at Subtotal - 1, at Shipping - 2, at Grand total - 3)",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path: ConstConfigPathGiftEmailTemplate,
		Value: `Dear {{.Recipient.Name}}, a gift card has been purchased on your behalf
			by {{.Buyer.Name}}
			<br />
			You are free to use this gift card at any time upon checkout.
			<br />
			<h3>Gift Cards</h3><br />
			Unique Code: {{.GiftCard.Code}}, Value: ${{.GiftCard.Amount}}`,
		Type:        env.ConstConfigTypeHTML,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Gift Card Recipeient E-mail: ",
		Description: "Email sent to the recipient address upon successful checkout",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
