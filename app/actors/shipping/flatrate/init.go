package flatrate

import (
	"github.com/ottemo/commerce/app"
	"github.com/ottemo/commerce/app/models/checkout"
	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"
)

// init makes package self-initialization routine
func init() {
	instance := new(ShippingMethod)

	app.OnAppStart(onAppStart)
	if err := checkout.RegisterShippingMethod(instance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "66cfb31e-66ff-43e8-bcd3-ddff50af3249", err.Error())
	}

	env.RegisterOnConfigStart(setupConfig)
}

// onAppStart makes module initialization on application startup
func onAppStart() error {

	rules, err := utils.DecodeJSONToArray(env.ConfigGetValue(ConstConfigPathRates))
	if err != nil {
		rules = make([]interface{}, 0)
	}

	flatRates = rules

	return nil
}
