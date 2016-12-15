package braintree

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/api"
)

// init makes package self-initialization routine
func init() {
	checkout.RegisterPaymentMethod(new(BraintreePaymentMethod))
	checkout.RegisterPaymentMethod(new(BraintreePaypalPaymentMethod))
	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
}

