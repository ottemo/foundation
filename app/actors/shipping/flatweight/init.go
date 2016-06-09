package flatweight

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

func init() {
	// rates global is auto-populated via the config declaration
	env.RegisterOnConfigStart(setupConfig)

	i := new(ShippingMethod)
	checkout.RegisterShippingMethod(i)
}
