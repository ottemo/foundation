package stripesubscription

import (
	//"github.com/ottemo/foundation/db"
	//"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	//"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/stripesubscription"
)

// init makes package self-initialization routine before app start
func init() {
	stripeSubscriptionInstance := new(DefaultStripeSubscription)
	var _ stripesubscription.InterfaceStripeSubscription = stripeSubscriptionInstance
	models.RegisterModel(stripesubscription.ConstModelNameStripeSubscription, stripeSubscriptionInstance)

	stripeSubscriptionCollectionInstance := new(DefaultStripeSubscriptionCollection)
	var _ stripesubscription.InterfaceStripeSubscriptionCollection = stripeSubscriptionCollectionInstance
	models.RegisterModel(stripesubscription.ConstModelNameStripeSubscriptionCollection, stripeSubscriptionCollectionInstance)

	//db.RegisterOnDatabaseStart(setupDB)
	//api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
	//app.OnAppStart(onAppStart)
}
