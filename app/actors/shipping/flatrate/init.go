package flatrate

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	instance := new(ShippingMethod)

	app.OnAppStart(onAppStart)
	checkout.RegisterShippingMethod(instance)

	env.RegisterOnConfigStart(setupConfig)
}

// onAppStart makes module initialization on application startup
func onAppStart() error {

	rules, err := utils.DecodeJSONToArray(env.ConfigGetValue(ConstConfigPathAdditionalRates))
	if err != nil {
		rules = make([]interface{}, 0)
	}

	additionalRates = rules

	return nil
}
