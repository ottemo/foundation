package shareasale

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

func init() {
	app.OnAppStart(appStart)
	env.RegisterOnConfigStart(setupConfig)
}

func appStart() error {
	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)

	return nil
}
