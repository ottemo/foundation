package authorizenet

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/checkout"
)

// init makes package self-initialization routine
func init() {
	if err := checkout.RegisterPaymentMethod(new(DirectPostMethod)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c36db383-87db-45be-9cab-ce5c0967686c", err.Error())
	}
	if err := checkout.RegisterPaymentMethod(new(RestMethod)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c7d65f90-7ba4-48b9-96b9-31e9d4cd49d0", err.Error())
	}
	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
}
