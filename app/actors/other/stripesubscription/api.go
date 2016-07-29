package stripesubscription

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/stripesubscription"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
	"github.com/ottemo/foundation/app/models/visitor"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	// Stripe subscription checkout
	service.POST("stripe/subscription", APISubscription)
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
func APISubscription(context api.InterfaceApplicationContext) (interface{}, error) {
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "287e280f-3567-4935-9cbf-0f2f1afd149e", "You should log in to subscribe")
	}
	visitorModel, err := visitor.LoadVisitorByID(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stripeSubscriptionModel, err := stripesubscription.GetStripeSubscriptionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Set shipping address
	//----------------------------
	shippingAddress := requestData["shipping_address"]
	if shippingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c91225f6-6752-4608-90a6-4119d008a25b", "Shipping address should be specified")
	}
	shippingAddressModel, err := checkout.ValidateAddress(shippingAddress)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = stripeSubscriptionModel.Set("shipping_address", shippingAddressModel.ToHashMap())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}


	// Set billing address
	//----------------------------
	billingAddress := requestData["billing_address"]
	if billingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "838ecfa5-3e58-45ef-b52d-061abeeeedc8", "Billing address should be specified")
	}
	billingAddressModel, err := checkout.ValidateAddress(billingAddress)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = stripeSubscriptionModel.Set("billing_address", billingAddressModel.ToHashMap())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Check subscription plan ID
	//----------------------------
	planID := utils.InterfaceToString(requestData["plan_id"])
	if planID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b153941b-27d3-4653-bc9f-61683b6047a9", "Subscription plan ID should be specified")
	}

	// Check credit card token
	//----------------------------
	ccToken := utils.InterfaceToString(requestData["cc_token"])
	if ccToken == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a20e1235-c23c-40f0-ae5a-5475abf3427e", "Credit card token should be specified")
	}

	stripe.Key = utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAPIKey))

	// Create customer on Stripe
	//----------------------------
	customerParams := &stripe.CustomerParams{
		Email: visitorModel.GetEmail(),
	}
	customerParams.SetSource(ccToken)
	customerForStripe, err := customer.New(customerParams)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = stripeSubscriptionModel.Set("stripe_customer_id", customerForStripe.ID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Create subscription on Stripe
	//----------------------------
	subscriptionForStripe, err := sub.New(&stripe.SubParams{
		Customer: customerForStripe.ID,
		Plan: planID,
	})
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = stripeSubscriptionModel.Set("stripe_subscription_id", subscriptionForStripe.ID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = stripeSubscriptionModel.Set("visitor_id", visitorModel.GetID())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = stripeSubscriptionModel.Set("customer_name", visitorModel.GetFullName())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = stripeSubscriptionModel.Set("customer_email", visitorModel.GetEmail())
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = stripeSubscriptionModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return stripeSubscriptionModel.ToHashMap(), nil
}

// APIGetCoupon returns stripe coupon
func APIGetCoupon(context api.InterfaceApplicationContext) (interface{}, error) {
	return nil, nil
}

// APIListSubscriptions returns a list of stripe subscriptions
func APIListSubscriptions(context api.InterfaceApplicationContext) (interface{}, error) {
	return nil, nil
}

// APIGetSubscription returns specified stripe subscription information
func APIGetSubscription(context api.InterfaceApplicationContext) (interface{}, error) {
	return nil, nil
}

// APIUpdateSubscription allows to update subscription status
func APIUpdateSubscription(context api.InterfaceApplicationContext) (interface{}, error) {
	return nil, nil
}

// APIListVisitorSubscriptions returns a list of visitor's stripe subscriptions
func APIListVisitorSubscriptions(context api.InterfaceApplicationContext) (interface{}, error) {
	return nil, nil
}

// APIUpdateVisitorSubscription allows visitor to cancel subscription
func APIUpdateVisitorSubscription(context api.InterfaceApplicationContext) (interface{}, error) {
	return nil, nil
}

// APIProcessStripeEvent listens to Stripe events and makes appropriate updates to subscriptions
func APIProcessStripeEvent(context api.InterfaceApplicationContext) (interface{}, error) {
	return nil, nil
}
