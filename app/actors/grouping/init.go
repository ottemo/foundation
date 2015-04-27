package grouping

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine before app start
func init() {
	app.OnAppStart(initListners)
}

// init Listeners for current model
func initListners() error {

	env.EventRegisterListener("api.cart.updatedCart", updateCartHandler)

	return nil
}
