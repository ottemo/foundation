package braintree

import (
	"github.com/ottemo/foundation/env"

	"github.com/lionelbarrow/braintree-go"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	checkout.RegisterPaymentMethod(new(CreditCardMethod))
	env.RegisterOnConfigStart(setupConfig)

	app.OnAppStart(onAppStart)
}

func onAppStart() error {
	braintreeInstance = braintree.New(
		braintree.Environment(utils.InterfaceToString(env.ConfigGetValue(ConstGeneralConfigPathEnvironment))),
		utils.InterfaceToString(env.ConfigGetValue(ConstGeneralConfigPathMerchantID)),
		utils.InterfaceToString(env.ConfigGetValue(ConstGeneralConfigPathPublicKey)),
		utils.InterfaceToString(env.ConfigGetValue(ConstGeneralConfigPathPrivateKey)),
	)

	return nil
}
