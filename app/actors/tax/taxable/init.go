package taxable

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/checkout"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultTaxable)

	if err := checkout.RegisterPriceAdjustment(instance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "475f02b7-88e2-468f-957f-f56cec01e643", err.Error())
	}
}

