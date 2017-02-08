package usps

import (
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/checkout"
)

// init makes package self-initialization routine before app start
func init() {
	instance := new(USPS)

	if err := checkout.RegisterShippingMethod(instance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "16fcb0df-c6ed-4ae2-bb4b-53468e587685", err.Error())
	}

	env.RegisterOnConfigStart(setupConfig)
}
