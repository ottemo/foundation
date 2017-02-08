package rts

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine before app start
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
	db.RegisterOnDatabaseStart(setupDB)
	app.OnAppStart(initListners)
	app.OnAppStart(initSalesHistory)
	app.OnAppStart(initStatistic)
}

// DB preparations for current model implementation
func initListners() error {
	// env.EventRegisterListener("api.rts.visit", referrerHandler)
	env.EventRegisterListener("api.rts.visit", visitsHandler)
	env.EventRegisterListener("api.cart.addToCart", addToCartHandler)
	env.EventRegisterListener("api.checkout.visit", visitCheckoutHandler)
	env.EventRegisterListener("api.checkout.setPayment", setPaymentHandler)
	env.EventRegisterListener("checkout.success", purchasedHandler)
	env.EventRegisterListener("checkout.success", salesHandler)

	return nil
}

// DB preparations for current model implementation
func setupDB() error {

	collection, err := db.GetCollection(ConstCollectionNameRTSSalesHistory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("product_id", db.ConstTypeID, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "91c4509d-525c-4a95-8e18-a2c25ed5ab8c", err.Error())
	}
	if err := collection.AddColumn("created_at", db.ConstTypeDatetime, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "73523525-4cf9-4a28-86c3-16ed43a473c1", err.Error())
	}
	if err := collection.AddColumn("count", db.ConstTypeInteger, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "af107716-93c6-4981-8ec7-57c19fbca40f", err.Error())
	}

	collection, err = db.GetCollection(ConstCollectionNameRTSVisitors)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("day", db.ConstTypeDatetime, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "37dc56f1-b367-4245-a3d3-4a12e21a1727", err.Error())
	}
	if err := collection.AddColumn("visitors", db.ConstTypeInteger, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b1bc428b-0d48-47a4-bfd5-a72f26caafe0", err.Error())
	}
	if err := collection.AddColumn("total_visits", db.ConstTypeInteger, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "39e7a5cb-c730-4b21-a5b3-0b5a0efbb4d5", err.Error())
	}
	if err := collection.AddColumn("cart", db.ConstTypeInteger, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b220ffd7-ac72-4106-a96f-d931f3323dde", err.Error())
	}
	if err := collection.AddColumn("visit_checkout", db.ConstTypeInteger, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4fd48a58-6e72-4e08-97eb-2545c0b68817", err.Error())
	}
	if err := collection.AddColumn("set_payment", db.ConstTypeInteger, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9f5b68fe-7d10-44a4-a51f-4e925ecec299", err.Error())
	}
	if err := collection.AddColumn("sales", db.ConstTypeInteger, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8d4a385a-adfc-44cb-8cd6-bbe15bb4e355", err.Error())
	}
	if err := collection.AddColumn("sales_amount", db.ConstTypeFloat, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "54eab14d-f3c8-465f-be8a-11c7770d8197", err.Error())
	}

	return nil
}
