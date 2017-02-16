package zeropay

import (
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/checkout"
)

// init makes package self-initialization routine
func init() {
	if err := checkout.RegisterPaymentMethod(new(ZeroAmountPayment)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c2f3aaed-7a1e-4b43-bc31-40cdec749ac5", err.Error())
	}
	env.RegisterOnConfigStart(setupConfig)
}
