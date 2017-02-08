package fedex

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	if err := checkout.RegisterShippingMethod(new(FedEx)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "619e374e-c247-4db5-81fd-2baf8dd6f9f6", err.Error())
	}
	env.RegisterOnConfigStart(setupConfig)
}
