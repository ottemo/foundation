package grouping

import (
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine before app start
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
	app.OnAppStart(initListners)
}

// DB preparations for current model implementation
func initListners() error {

	env.EventRegisterListener("api.cart.updatedCart", updateCartHandler)

	return nil
}


