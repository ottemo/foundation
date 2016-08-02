package stripesubscription

import (
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

	// Stripe Subscription config section
	//----------------------------
	err := config.RegisterItem(env.StructConfigItem{
		Path:  ConstConfigPathGroup,
		Label: "Stripe Subscription",
		Type:  env.ConstConfigTypeGroup,
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Enabled
	//----------------------------
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathEnabled,
		Value:       false,
		Type:        env.ConstConfigTypeBoolean,
		Editor:      "boolean",
		Options:     nil,
		Label:       "Enable",
		Description: "Enables/disables subscritions through Stripe",
		Image:       "",
	}, func(value interface{}) (interface{}, error) { return utils.InterfaceToBool(value), nil })

	// Stripe API key
	//----------------------------
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathAPIKey,
		Label:       "API Key",
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Description: "Your API Key will be located in your Stripe Dashboard.",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Stripe API key
	//----------------------------
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathAPIKey,
		Label:       "API Key",
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Description: "Your API Key will be located in your Stripe Dashboard.",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Subscription plans config
	//----------------------------
	err = config.RegisterItem(env.StructConfigItem{
		Path:    ConstConfigPathPlans,
		Value:   nil,
		Type:    env.ConstConfigTypeText,
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

	// Cancel subscription email subject
	//----------------------------
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathEmailCancelSubject,
		Value:       "Subscription was canceled",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Cancel Subscription Email: Subject",
		Description: "",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Cancel subscription email body
	//----------------------------
	err = config.RegisterItem(env.StructConfigItem{
		Path: ConstConfigPathEmailCancelTemplate,
		Value: `Dear {{.Visitor.name}},
Your subscription was canceled`,
		Type:        env.ConstConfigTypeHTML,
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Cancel Subscription Email: Body",
		Description: "",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Charge date
	//----------------------------
	err = config.RegisterItem(env.StructConfigItem{
		Path: ConstConfigPathChargeDate,
		Value: 25,
		Type:        env.ConstConfigTypeInteger,
		Editor:      "text",
		Options:     "",
		Label:       "Charge on Date",
		Description: "Users will be charged on this date of month",
		Image:       "",
	}, nil)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
