package subscription

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine before app start
func init() {
	subscriptionInstance := new(DefaultSubscription)
	var _ subscription.InterfaceSubscription = subscriptionInstance
	if err := models.RegisterModel(subscription.ConstModelNameSubscription, subscriptionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7f2293e4-1f6c-4fed-8410-83f99f145e65", err.Error())
	}

	subscriptionCollectionInstance := new(DefaultSubscriptionCollection)
	var _ subscription.InterfaceSubscriptionCollection = subscriptionCollectionInstance
	if err := models.RegisterModel(subscription.ConstModelNameSubscriptionCollection, subscriptionCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c7513220-d5e5-4770-87b4-6405fced1c9d", err.Error())
	}

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
	app.OnAppStart(onAppStart)
}

// setupDB prepares system database for package usage
func setupDB() error {

	collection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := collection.AddColumn("visitor_id", db.ConstTypeID, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "36a415d0-08ea-4ea0-8205-cfbb690ab10a", err.Error())
	}
	if err := collection.AddColumn("order_id", db.ConstTypeID, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1864f423-be12-4b70-b910-366420cfbfcf", err.Error())
	}

	if err := collection.AddColumn("items", db.TypeArrayOf(db.ConstTypeJSON), true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3726acfb-afeb-4a1a-a8fe-f251bdfa1196", err.Error())
	}

	if err := collection.AddColumn("customer_email", db.TypeWPrecision(db.ConstTypeVarchar, 100), true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2a4ce601-2f55-46c7-a79b-09860c97fb92", err.Error())
	}
	if err := collection.AddColumn("customer_name", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ffd385b2-d16d-456e-a226-3a11122cd789", err.Error())
	}

	if err := collection.AddColumn("shipping_address", db.ConstTypeJSON, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1f39d1a1-d8b9-4b8d-8079-a6bc984684f8", err.Error())
	}
	if err := collection.AddColumn("billing_address", db.ConstTypeJSON, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "be736068-2e21-4572-94f9-4190675fdde9", err.Error())
	}

	if err := collection.AddColumn("shipping_method", db.TypeWPrecision(db.ConstTypeVarchar, 100), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "22f061f4-e69b-4167-b52f-88877c93a0bb", err.Error())
	}
	if err := collection.AddColumn("shipping_rate", db.ConstTypeJSON, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ac231f31-2f64-466a-977d-7205794b2717", err.Error())
	}

	if err := collection.AddColumn("payment_instrument", db.ConstTypeJSON, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a8b657dc-4faf-476b-a634-42b1f942325b", err.Error())
	}

	if err := collection.AddColumn("period", db.ConstTypeInteger, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "98bc7437-a59f-4056-a833-5fa0a5663b1a", err.Error())
	}
	if err := collection.AddColumn("status", db.TypeWPrecision(db.ConstTypeVarchar, 50), false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "726d93f5-d0cd-463f-9e02-51d1d41ac456", err.Error())
	}

	if err := collection.AddColumn("action_date", db.ConstTypeDatetime, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "236c2336-f857-404d-bf70-0bf00f051850", err.Error())
	}
	if err := collection.AddColumn("last_submit", db.ConstTypeDatetime, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "029666db-8126-4f92-a6bb-758569b1cb34", err.Error())
	}

	if err := collection.AddColumn("created_at", db.ConstTypeDatetime, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0c41472f-38c5-4ff9-83f8-6563afbe93b3", err.Error())
	}
	if err := collection.AddColumn("updated_at", db.ConstTypeDatetime, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bd4a1231-5eee-40c4-ad4f-edeba80f621e", err.Error())
	}

	if err := collection.AddColumn("info", db.ConstTypeJSON, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "10cfb39e-7f3b-44fd-9159-a852cdd86943", err.Error())
	}

	return nil
}

// onAppStart makes module initialization on application startup
func onAppStart() error {

	products := make([]string, 0)
	productsValue := utils.InterfaceToArray(env.ConfigGetValue(subscription.ConstConfigPathSubscriptionProducts))

	for _, value := range productsValue {
		if productID := utils.InterfaceToString(value); productID != "" {
			products = append(products, productID)
		}
	}

	subscriptionProducts = products

	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)
	env.EventRegisterListener("product.getOptions", getOptionsExtend)

	// process order creation every one hour
	if scheduler := env.GetScheduler(); scheduler != nil {
		if err := scheduler.RegisterTask(ConstSchedulerTaskName, placeOrders); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "10cfb39e-7f3b-44fd-9159-a852cdd86943", err.Error())
		}
		if _, err := scheduler.ScheduleRepeat(
			subscription.GetSubscriptionCronExpr(
				subscription.GetSubscriptionPeriodValue(
					utils.InterfaceToString(
						env.ConfigGetValue(subscription.ConstConfigPathSubscriptionExecutionTime)))),
			ConstSchedulerTaskName,
			nil); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "30717db6-0e30-4656-8698-531c5c4d60f6", err.Error())
		}
	}

	return nil
}
