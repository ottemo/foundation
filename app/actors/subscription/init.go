package subscription

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine before app start
func init() {
	subscriptionInstance := new(DefaultSubscription)
	var _ subscription.InterfaceSubscription = subscriptionInstance
	models.RegisterModel(subscription.ConstModelNameSubscription, subscriptionInstance)

	subscriptionCollectionInstance := new(DefaultSubscriptionCollection)
	var _ subscription.InterfaceSubscriptionCollection = subscriptionCollectionInstance
	models.RegisterModel(subscription.ConstModelNameSubscriptionCollection, subscriptionCollectionInstance)

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

	collection.AddColumn("visitor_id", db.ConstTypeID, true)
	collection.AddColumn("order_id", db.ConstTypeID, true)
	collection.AddColumn("cart_id", db.ConstTypeID, true)

	collection.AddColumn("email", db.TypeWPrecision(db.ConstTypeVarchar, 100), true)
	collection.AddColumn("name", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)

	collection.AddColumn("shipping_address", db.ConstTypeJSON, false)
	collection.AddColumn("billing_address", db.ConstTypeJSON, false)

	collection.AddColumn("shipping_method", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)
	collection.AddColumn("shipping_rate", db.ConstTypeJSON, false)

	collection.AddColumn("payment_instrument", db.ConstTypeJSON, false)

	collection.AddColumn("action_date", db.ConstTypeDatetime, true)
	collection.AddColumn("last_submit", db.ConstTypeDatetime, true)

	collection.AddColumn("created_at", db.ConstTypeDatetime, true)
	collection.AddColumn("updated_at", db.ConstTypeDatetime, true)

	collection.AddColumn("period", db.ConstTypeInteger, false)

	collection.AddColumn("status", db.TypeWPrecision(db.ConstTypeVarchar, 50), false)
	collection.AddColumn("action", db.TypeWPrecision(db.ConstTypeVarchar, 50), false)

	return nil
}

// onAppStart makes module initialization on application startup
func onAppStart() error {

	env.EventRegisterListener("checkout.success", checkoutSuccessHandler)
	env.EventRegisterListener("product.getOptions", getOptionsExtend)

	// process order creation every one hour
	if scheduler := env.GetScheduler(); scheduler != nil {
		scheduler.RegisterTask("subscriptionProcess", placeOrders)
		scheduler.ScheduleRepeat("* */1 * * *", "subscriptionProcess", nil)
	}

	return nil
}
