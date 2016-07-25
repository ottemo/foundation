package stripesubscription

import (
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "1f7ecfb8-b5e3-4361-b066-42c088f6b350", "can't obtain config")
		return env.ErrorDispatch(err)
	}

	// Subscription plans config
	//----------------------------
	err := config.RegisterItem(env.StructConfigItem{
		Path:    ConstConfigPathPlans,
		Value:   nil,
		Type:    env.ConstConfigTypeJSON,
		Editor:  "multiline_text",
		Options: nil,
		Label:   "Subscription plans",
		Description: `Subscription plans settings, pattern:
[
	{"id": "monthPlanStripeId", "name": "Month to Month", "price": 36, "deliveryCount": 1, "info": "Charged every month", "img": "/images/subscribe/1month.png"},
	...
]`,
		Image: "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
