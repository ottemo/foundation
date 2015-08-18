package coupon

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
		Path:        ConstConfigPathDiscounts,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Discounts",
		Description: "Discounts related options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathDiscountApplyPriority,
		Value:       2.10,
		Type:        env.ConstConfigTypeFloat,
		Editor:      "line_text",
		Options:     nil,
		Label:       "Discounts calculating position",
		Description: "This value used for using position to calculate it's possible applicable amount (Subtotal - 1, Shipping - 2, Grand total - 3)",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
