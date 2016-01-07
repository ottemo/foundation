package mailchimp

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

func init() {
	app.OnAppStart(appStart)
	env.RegisterOnConfigStart(setupConfig)
}

func appStart() {
	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)
}
