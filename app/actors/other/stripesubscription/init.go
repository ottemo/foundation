package stripesubscription

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/stripesubscription"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine before app start
func init() {
	stripeSubscriptionInstance := new(DefaultStripeSubscription)
	var _ stripesubscription.InterfaceStripeSubscription = stripeSubscriptionInstance
	models.RegisterModel(stripesubscription.ConstModelNameStripeSubscription, stripeSubscriptionInstance)

	stripeSubscriptionCollectionInstance := new(DefaultStripeSubscriptionCollection)
	var _ stripesubscription.InterfaceStripeSubscriptionCollection = stripeSubscriptionCollectionInstance
	models.RegisterModel(stripesubscription.ConstModelNameStripeSubscriptionCollection, stripeSubscriptionCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
}

// setupDB prepares system database for package usage
func setupDB() error {

	collection, err := db.GetCollection(ConstCollectionNameStripeSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("visitor_id", db.ConstTypeID, true)

	collection.AddColumn("customer_name", db.TypeWPrecision(db.ConstTypeVarchar, 100), true)
	collection.AddColumn("customer_email", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)

	collection.AddColumn("billing_address", db.ConstTypeJSON, false)
	collection.AddColumn("shipping_address", db.ConstTypeJSON, false)

	collection.AddColumn("stripe_subscription_id", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)
	collection.AddColumn("stripe_customer_id", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)
	collection.AddColumn("stripe_coupon", db.TypeWPrecision(db.ConstTypeVarchar, 100), false)
	collection.AddColumn("last_payment_info", db.ConstTypeJSON, false)
	collection.AddColumn("next_payment_at", db.ConstTypeDatetime, false)

	collection.AddColumn("price", db.ConstTypeMoney, false)

	collection.AddColumn("created_at", db.ConstTypeDatetime, false)
	collection.AddColumn("updated_at", db.ConstTypeDatetime, false)

	collection.AddColumn("description", db.TypeWPrecision(db.ConstTypeVarchar, 200), false)
	collection.AddColumn("info", db.ConstTypeJSON, false)
	collection.AddColumn("status", db.TypeWPrecision(db.ConstTypeVarchar, 100), true)

	return nil
}
