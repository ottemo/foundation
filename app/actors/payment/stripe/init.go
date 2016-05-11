package stripe

import (
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/checkout"
)

func init() {
	checkout.RegisterPaymentMethod(new(Payment))
	env.RegisterOnConfigStart(setupConfig)
}
