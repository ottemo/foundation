package coupon

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultDiscount)
	var _ checkout.InterfaceDiscount = instance
	checkout.RegisterDiscount(instance)

	db.RegisterOnDatabaseStart(setupDB)
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)

	app.OnAppStart(initListeners)
}

// setupDB prepares system database for package usage
func setupDB() error {

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		env.ErrorDispatch(err)
	}

	collection.AddColumn("code", db.ConstTypeVarchar, true)
	collection.AddColumn("name", db.ConstTypeVarchar, false)
	collection.AddColumn("amount", db.ConstTypeDecimal, false)
	collection.AddColumn("percent", db.ConstTypeDecimal, false)
	collection.AddColumn("times", db.ConstTypeInteger, false)
	collection.AddColumn("since", db.ConstTypeDatetime, false)
	collection.AddColumn("until", db.ConstTypeDatetime, false)
	collection.AddColumn("limits", db.ConstTypeJSON, false)

	return nil
}

// initListeners register event listeners
func initListeners() error {

	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)
	return nil
}
