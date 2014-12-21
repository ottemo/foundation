package usps

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine before app start
func init() {
	instance := new(USPS)

	checkout.RegisterShippingMethod(instance)

	env.RegisterOnConfigStart(setupConfig)
}
