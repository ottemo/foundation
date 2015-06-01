package trustpilot

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine before app start
func init() {
	app.OnAppStart(initListeners)
	env.RegisterOnConfigStart(setupConfig)
}

// onAppStart makes module initialization on application startup
func initListeners() error {

	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)

	return nil
}
