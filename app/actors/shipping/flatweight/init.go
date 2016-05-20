package flatweight

import (
	"fmt"
	"github.com/ottemo/foundation/app"
	// "github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

func init() {

	// i := new(ShippingMethod)
	app.OnAppStart(onAppStart)
	// checkout.RegisterShippingMethod(i)

	env.RegisterOnConfigStart(setupConfig)
}

func onAppStart() error {
	fmt.Println("rates loaded from db", rates)
	configRates := configRates()
	fmt.Println("loaded this", configRates)
	_, err := validateAndApplyRates(configRates)
	if err != nil {
		fmt.Println("errord trying to validate", err)
		// eek we have bad data in the DB
		env.ErrorDispatch(env.ErrorNew(ConstErrorModule, 1, "todo", "Failed to create flat weight rates from config in DB"))
		rates = make(Rates, 0)
	}

	return nil
}
