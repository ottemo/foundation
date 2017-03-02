package stripe

import (
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/checkout"
)

func init() {
	if err := checkout.RegisterPaymentMethod(new(Payment)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "ea29fa2a-f947-4e7f-aff0-b0965256c751", err.Error())
	}
	env.RegisterOnConfigStart(setupConfig)
}
