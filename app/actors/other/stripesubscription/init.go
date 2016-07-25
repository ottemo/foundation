package stripesubscription

import (
	//"github.com/ottemo/foundation/app/models"
	//"github.com/ottemo/foundation/db"
	//"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	//"github.com/ottemo/foundation/app"
)

// init makes package self-initialization routine before app start
func init() {
	//subscriptionInstance := new(DefaultSubscription)
	//var _ subscription.InterfaceSubscription = subscriptionInstance
	//models.RegisterModel(subscription.ConstModelNameStripeSubscription, subscriptionInstance)

	//subscriptionCollectionInstance := new(DefaultSubscriptionCollection)
	//var _ subscription.InterfaceSubscriptionCollection = subscriptionCollectionInstance
	//models.RegisterModel(subscription.ConstModelNameStripeSubscriptionCollection, subscriptionCollectionInstance)

	//db.RegisterOnDatabaseStart(setupDB)
	//api.RegisterOnRestServiceStart(setupAPI)
	env.RegisterOnConfigStart(setupConfig)
	//app.OnAppStart(onAppStart)
}
