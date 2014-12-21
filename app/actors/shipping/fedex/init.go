package fedex

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	checkout.RegisterShippingMethod(new(FedEx))
	env.RegisterOnConfigStart(setupConfig)
}
