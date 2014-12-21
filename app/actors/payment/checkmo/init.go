package checkmo

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	checkout.RegisterPaymentMethod(new(CheckMoneyOrder))
	env.RegisterOnConfigStart(setupConfig)
}
