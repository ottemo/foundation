package taxcloud

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/checkout"
)

func init() {
	instance := new(DefaultTaxCloud)

	if err := checkout.RegisterPriceAdjustment(instance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ff476d06-012e-42bb-a04e-8f045a58a4d8", err.Error())
	}

	env.RegisterOnConfigStart(setupConfig)
}
