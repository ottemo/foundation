package checkmo

import (
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/checkout"
)

// init makes package self-initialization routine
func init() {
	if err := checkout.RegisterPaymentMethod(new(CheckMoneyOrder)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "aa6367b6-3bf2-4843-915b-be8fbe3ffa6a", err.Error())
	}
	env.RegisterOnConfigStart(setupConfig)
}
