package braintree

import (
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/checkout"
)

// init makes package self-initialization routine
func init() {
	if err := checkout.RegisterPaymentMethod(new(CreditCardMethod)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "194b26f7-3399-4121-80f1-e24305708871", err.Error())
	}
	env.RegisterOnConfigStart(setupConfig)
}
