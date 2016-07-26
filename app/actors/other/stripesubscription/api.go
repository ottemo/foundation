package stripesubscription

import "github.com/ottemo/foundation/api"

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	// Stripe subscription checkout
	service.POST("stripe/subscribe", APISubscribe)
	service.GET("stripe/coupon/:id", APIGetCoupon)

	// Administrative
	service.GET("stripe/subscriptions", api.IsAdmin(APIListSubscriptions))
	service.GET("stripe/subscription/:id", api.IsAdmin(APIGetSubscription))
	service.PUT("stripe/subscription/:id", api.IsAdmin(APIUpdateSubscription))

	// Visitor Account
	service.GET("visit/stripe/subscriptions", APIListVisitorSubscriptions)
	service.PUT("visit/stripe/subscription/:id", APIUpdateVisitorSubscription)

	// Stripe events
	service.POST("stripe/subscriptions", APIProcessStripeEvent)

	return nil
}

// APISubscribe creates stripe subscription
func APISubscribe(context api.InterfaceApplicationContext) (interface{}, error) {

}

// APIGetCoupon returns stripe coupon
func APIGetCoupon(context api.InterfaceApplicationContext) (interface{}, error) {

}

// APIListSubscriptions returns a list of stripe subscriptions
func APIListSubscriptions(context api.InterfaceApplicationContext) (interface{}, error) {

}

// APIGetSubscription returns specified stripe subscription information
func APIGetSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

}

// APIUpdateSubscription allows to update subscription status
func APIUpdateSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

}

// APIListVisitorSubscriptions returns a list of visitor's stripe subscriptions
func APIListVisitorSubscriptions(context api.InterfaceApplicationContext) (interface{}, error) {

}

// APIUpdateVisitorSubscription allows visitor to cancel subscription
func APIUpdateVisitorSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

}

// APIProcessStripeEvent listens to Stripe events and makes appropriate updates to subscriptions
func APIProcessStripeEvent(context api.InterfaceApplicationContext) (interface{}, error) {

}
