package checkout

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultCheckout)
	var _ checkout.InterfaceCheckout = instance
	if err := models.RegisterModel(checkout.ConstCheckoutModelName, instance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6dcce70e-7d29-46fb-a91a-6dc6d2836695", err.Error())
	}

	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
}
