package emma

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

func init() {
	app.OnAppStart(appStart)
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)
}

func appStart() error {
	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)

	return nil
}
