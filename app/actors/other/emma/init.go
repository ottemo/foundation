package emma

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/api"
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
