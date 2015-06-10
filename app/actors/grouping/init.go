package grouping

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine before app start
func init() {
	app.OnAppStart(onAppStart)
	env.RegisterOnConfigStart(setupConfig)
}

// onAppStart makes module initialization on application startup
func onAppStart() error {

	rules, err := utils.DecodeJSONToArray(env.ConfigGetValue(ConstGroupingConfigPath))
	if err != nil {
		rules = make([]interface{}, 0)
	}
	currentRules = rules

	env.EventRegisterListener("api.cart.update", updateCartHandler)

	return nil
}
