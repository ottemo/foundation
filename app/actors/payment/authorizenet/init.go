package authorizenet

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	checkout.RegisterPaymentMethod(new(DirectPostMethod))
	checkout.RegisterPaymentMethod(new(RestAPI))
	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
}
