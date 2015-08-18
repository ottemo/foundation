package giftcard

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultGiftcard)
	var _ checkout.InterfaceDiscount = instance
	checkout.RegisterDiscount(instance)

	db.RegisterOnDatabaseStart(setupDB)
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)

	app.OnAppStart(initListeners)
}

// DB preparations for current model implementation
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("code", db.ConstTypeID, true)
	collection.AddColumn("name", db.TypeWPrecision(db.ConstTypeVarchar, 150), false)
	collection.AddColumn("sku", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)

	collection.AddColumn("amount", db.ConstTypeMoney, false)

	collection.AddColumn("order_id", db.ConstTypeID, true)
	collection.AddColumn("visitor_id", db.ConstTypeID, false)

	collection.AddColumn("status", db.TypeWPrecision(db.ConstTypeVarchar, 50), false)
	collection.AddColumn("orders_used", db.ConstTypeJSON, false)
	collection.AddColumn("recipient_mailbox", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)

	return nil
}

// initListeners register event listeners
func initListeners() error {

	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)
	env.EventRegisterListener("order.proceed", orderProceedHandler)
	env.EventRegisterListener("order.rollback", orderRollbackHandler)

	return nil
}
