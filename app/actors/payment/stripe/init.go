package stripe

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

func init() {
	checkout.RegisterPaymentMethod(new(Payment))
	env.RegisterOnConfigStart(setupConfig)
}
