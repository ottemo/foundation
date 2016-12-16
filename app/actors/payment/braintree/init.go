package braintree

import (
	//"github.com/ottemo/foundation/api" // TODO
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/checkout"
)

// init makes package self-initialization routine
func init() {
	checkout.RegisterPaymentMethod(new(braintreeCCMethod))
	//api.RegisterOnRestServiceStart(setupAPI) // TODO
	env.RegisterOnConfigStart(setupConfig)
}
