package stripesubscription

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/stripesubscription"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
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

// APISubscription creates stripe subscription
func APISubscription(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Check visitor
	//----------------------------
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "287e280f-3567-4935-9cbf-0f2f1afd149e", "You should log in to subscribe")
	}
	visitorInstance, err := visitor.LoadVisitorByID(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Check shipping address
	//----------------------------
	reqShippingAddress := requestData["shipping_address"]
	if reqShippingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c91225f6-6752-4608-90a6-4119d008a25b", "Shipping address should be specified")
	}
	shippingAddress, err := checkout.ValidateAddress(reqShippingAddress)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Check billing address
	//----------------------------
	reqBillingAddress := requestData["billing_address"]
	if reqBillingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "838ecfa5-3e58-45ef-b52d-061abeeeedc8", "Billing address should be specified")
	}
	billingAddress, err := checkout.ValidateAddress(reqBillingAddress)
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

	// Check stripe api key
	//----------------------------
	stripeKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAPIKey))
	if stripeKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "057379c6-0664-465e-881f-082b6bafab48", "Stripe API key is empty")
	}

	// Process subscription
	//----------------------------
	stripeSubscriptionInstance, err := stripesubscription.GetStripeSubscriptionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stripe.Key = stripeKey

	// Create new customer for Stripe
	customerParams := &stripe.CustomerParams{
		Email: visitorInstance.GetEmail(),
	}
	customerParams.SetSource(ccToken)
	respStripeCustomer, err := customer.New(customerParams)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Create new subscription for Stripe
	respStripeSubscription, err := sub.New(&stripe.SubParams{
		Customer: respStripeCustomer.ID,
		Plan:     planID,
	})
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Save subscription
	//----------------------------
	stripeSubscriptionInstance.Set("visitor_id", visitorInstance.GetID())
	stripeSubscriptionInstance.Set("customer_name", visitorInstance.GetFullName())
	stripeSubscriptionInstance.Set("customer_email", visitorInstance.GetEmail())
	stripeSubscriptionInstance.Set("billing_address", billingAddress.ToHashMap())
	stripeSubscriptionInstance.Set("shipping_address", shippingAddress.ToHashMap())
	stripeSubscriptionInstance.Set("stripe_subscription_id", respStripeSubscription.ID)
	stripeSubscriptionInstance.Set("stripe_customer_id", respStripeCustomer.ID)
	stripeSubscriptionInstance.Set("next_payment_at", respStripeSubscription.PeriodEnd)
	stripeSubscriptionInstance.Set("price", respStripeSubscription.Plan.Amount)
	stripeSubscriptionInstance.Set("description", respStripeSubscription.Plan.Name)
	stripeSubscriptionInstance.Set("status", "confirmed")

	err = stripeSubscriptionInstance.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return stripeSubscriptionInstance.ToHashMap(), nil
}

// APIGetCoupon returns stripe coupon
func APIGetCoupon(context api.InterfaceApplicationContext) (interface{}, error) {
	return nil, nil
}

// APIListSubscriptions returns a list of stripe subscriptions
func APIListSubscriptions(context api.InterfaceApplicationContext) (interface{}, error) {
	// List operation
	stripeSubscriptionCollectionModel, err := stripesubscription.GetStripeSubscriptionCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	models.ApplyFilters(context, stripeSubscriptionCollectionModel.GetDBCollection())

	// Check for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return stripeSubscriptionCollectionModel.GetDBCollection().Count()
	}

	// Limit parameter handle
	stripeSubscriptionCollectionModel.ListLimit(models.GetListLimit(context))

	// Extra parameter handle
	models.ApplyExtraAttributes(context, stripeSubscriptionCollectionModel)

	return stripeSubscriptionCollectionModel.List()
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
