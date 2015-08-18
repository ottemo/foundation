package giftcard

import (
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "15859fac-8fc0-4fbf-a801-b9cacf70d356", "can't obtain config")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGiftCardSKU,
		Value:       "gift-card",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Gift cards SKU identifier",
		Description: "This value will be checked on presense in product SKU and it will be a flag for gift card operations",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGiftEmailSubject,
		Value:       "Your giftcard has arrived",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Email subject",
		Description: "This value will be in email subject",
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
		Label:       "Gift cards calculating position",
		Description: "This value used for using position to calculate it's possible applicable amount (Subtotal - 1, Shipping - 2, Grand total - 3)",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path: ConstConfigPathGiftEmailTemplate,
		Value: `Dear friend, you has received these gifts
		from {{.Visitor.name}}
<br />
You are free to use this gift card's at any time
<br />
<h3>Gift cards</h3><br />
Gift code: {{.GiftCard.Code}}, amount: ${{.GiftCard.Amount}}`,
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Gift card data send e-mail: ",
		Description: "contents of email will be sent to the specified address on success checkout",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
