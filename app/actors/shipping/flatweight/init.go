package flatweight

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

func init() {
	// rates global is auto-populated via the config declaration
	env.RegisterOnConfigStart(setupConfig)

	i := new(ShippingMethod)
	if err := checkout.RegisterShippingMethod(i); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "30e9fc93-1841-4429-b5a1-c7c6cf7cd3b7", err.Error())
	}
}
