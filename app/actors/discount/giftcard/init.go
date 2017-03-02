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
	var _ checkout.InterfacePriceAdjustment = instance

	if err := checkout.RegisterPriceAdjustment(instance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b1222356-9ec9-4071-8929-8476300a637d", err.Error())
	}

	freeShipping := new(Shipping)
	var _ checkout.InterfaceShippingMethod = freeShipping
	if err := checkout.RegisterShippingMethod(freeShipping); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7be6d453-a42e-4ec1-9fe5-1848c40271d8", err.Error())
	}

	db.RegisterOnDatabaseStart(setupDB)
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)

	app.OnAppStart(onAppStart)
}

// DB preparations for current model implementation
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameGiftCard)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("code", db.ConstTypeID, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fd0e3ffa-1cd8-4ffc-90e6-a7331b9f6bce", err.Error())
	}
	if err := collection.AddColumn("sku", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1c8a8a18-b7bb-4eae-8337-eb4b8d742515", err.Error())
	}

	if err := collection.AddColumn("amount", db.ConstTypeMoney, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2d391f8b-491e-4498-8537-a5690392cf28", err.Error())
	}

	if err := collection.AddColumn("order_id", db.ConstTypeID, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "668cb103-26a3-41c2-a7a9-0ceb4c194771", err.Error())
	}
	if err := collection.AddColumn("visitor_id", db.ConstTypeID, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dee210a7-7509-4793-9dab-9d93e3a2ff7e", err.Error())
	}

	if err := collection.AddColumn("status", db.TypeWPrecision(db.ConstTypeVarchar, 50), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1ef30ced-2c95-413a-a0f5-b97697751b34", err.Error())
	}
	if err := collection.AddColumn("orders_used", db.ConstTypeJSON, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "609148a7-49f1-4984-aa12-ec13a199fc52", err.Error())
	}

	if err := collection.AddColumn("name", db.TypeWPrecision(db.ConstTypeVarchar, 150), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c218c711-3fab-4115-a94b-cb81121dae0a", err.Error())
	}
	if err := collection.AddColumn("message", db.TypeWPrecision(db.ConstTypeVarchar, 150), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7df99934-7652-4388-8f0e-1bfc1aaa518f", err.Error())
	}

	if err := collection.AddColumn("recipient_mailbox", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "749c233d-aae9-4866-ba71-eef92eaa52a6", err.Error())
	}
	if err := collection.AddColumn("delivery_date", db.ConstTypeDatetime, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3ea6c6e1-1c37-4630-8774-4d92022eba5f", err.Error())
	}
	if err := collection.AddColumn("created_at", db.ConstTypeDatetime, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d5b22b68-91a5-4f00-8459-c5ef53a5d4bb", err.Error())
	}

	return nil
}

// onAppStart makes module initialization on application startup
func onAppStart() error {

	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)
	env.EventRegisterListener("order.proceed", orderProceedHandler)
	env.EventRegisterListener("order.rollback", orderRollbackHandler)

	if scheduler := env.GetScheduler(); scheduler != nil {
		if err := scheduler.RegisterTask("sendGiftCards", SendTask); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a3b21acc-4925-4ecc-b8b1-f7ca2f779b1a", err.Error())
		}
		if _, err := scheduler.ScheduleRepeat("0 8 * * *", "sendGiftCards", nil); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "06580361-37bf-44d9-bbe4-643fc6daad6c", err.Error())
		}
	}

	return nil
}
