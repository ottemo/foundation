package paypal

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	checkout.RegisterPaymentMethod(new(Express))
	checkout.RegisterPaymentMethod(new(PayFlowAPI))
	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
}
