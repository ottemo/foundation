package stripesubscription

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/stripesubscription"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/event"
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
	service.PUT("stripe/subscription/:id/cancel", APICancelSubscription)

	// Visitor Account
	service.GET("visit/stripe/subscriptions", APIListVisitorSubscriptions)
	service.PUT("visit/stripe/subscription/:id/cancel", APICancelSubscription)

	// Stripe events
	service.POST("stripe/event", APIProcessStripeEvent)

	return nil
}

// APISubscription creates stripe subscription
func APISubscription(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Validate visitor
	//----------------------------
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "287e280f-3567-4935-9cbf-0f2f1afd149e", "You should log in to subscribe")
	}
	visitorInstance, err := visitor.LoadVisitorByID(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Validate shipping address
	//----------------------------
	reqShippingAddress := requestData["shipping_address"]
	if reqShippingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c91225f6-6752-4608-90a6-4119d008a25b", "Shipping address should be specified")
	}
	shippingAddress, err := checkout.ValidateAddress(reqShippingAddress)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Validate billing address
	//----------------------------
	reqBillingAddress := requestData["billing_address"]
	if reqBillingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "838ecfa5-3e58-45ef-b52d-061abeeeedc8", "Billing address should be specified")
	}
	billingAddress, err := checkout.ValidateAddress(reqBillingAddress)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Validate subscription plan ID
	//----------------------------
	planID := utils.InterfaceToString(requestData["plan_id"])
	if planID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b153941b-27d3-4653-bc9f-61683b6047a9", "Subscription plan ID should be specified")
	}

	// Validate credit card token
	//----------------------------
	ccToken := utils.InterfaceToString(requestData["cc_token"])
	if ccToken == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a20e1235-c23c-40f0-ae5a-5475abf3427e", "Credit card token should be specified")
	}

	// Validate stripe api key
	//----------------------------
	stripeKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAPIKey))
	if stripeKey == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "057379c6-0664-465e-881f-082b6bafab48", "Stripe API key is empty")
	}
	stripe.Key = stripeKey

	stripeSubscriptionInstance, err := stripesubscription.GetStripeSubscriptionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Create new customer and subscription on Stripe
	customerParams := &stripe.CustomerParams{
		Email: visitorInstance.GetEmail(),
		Plan:  planID,
		Shipping: &stripe.CustomerShippingDetails{
			Name:  shippingAddress.GetFirstName() + " " + shippingAddress.GetLastName(),
			Phone: shippingAddress.GetPhone(),
			Address: stripe.Address{
				Line1:   shippingAddress.GetAddressLine1(),
				Line2:   shippingAddress.GetAddressLine2(),
				City:    shippingAddress.GetCity(),
				State:   shippingAddress.GetState(),
				Zip:     shippingAddress.GetZipCode(),
				Country: shippingAddress.GetCountry(),
			},
		},
	}
	customerParams.SetSource(ccToken)
	stripeCustomer, err := customer.New(customerParams)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f3825627-1dd1-475f-a6f7-d3acaebe590e", err.Error())
	}

	// Save subscription
	//----------------------------
	stripeSubscriptionInstance.Set("visitor_id", visitorInstance.GetID())
	stripeSubscriptionInstance.Set("customer_name", visitorInstance.GetFullName())
	stripeSubscriptionInstance.Set("customer_email", visitorInstance.GetEmail())
	stripeSubscriptionInstance.Set("billing_address", billingAddress.ToHashMap())
	stripeSubscriptionInstance.Set("shipping_address", shippingAddress.ToHashMap())
	stripeSubscriptionInstance.Set("description", stripeCustomer.Subs.Values[0].Plan.Name)
	stripeSubscriptionInstance.Set("status", stripeCustomer.Subs.Values[0].Status)
	stripeSubscriptionInstance.Set("stripe_customer_id", stripeCustomer.ID)
	stripeSubscriptionInstance.Set("price", stripeCustomer.Subs.Values[0].Plan.Amount)
	stripeSubscriptionInstance.Set("period_end", stripeCustomer.Subs.Values[0].PeriodEnd)
	stripeSubscriptionInstance.Set("notify_renew", requestData["notify_on_renew"])
	stripeSubscriptionInstance.Set("renew_notified", false)
	if utils.InterfaceToBool(requestData["is_gift"]) == true {
		stripeSubscriptionInstance.Set("info", map[string]interface{}{"Gift": true})
	}

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

// APICancelSubscription cancels subscription
func APICancelSubscription(context api.InterfaceApplicationContext) (interface{}, error) {
	// Validate params
	id := context.GetRequestArgument("id")
	if id == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9147deb7-66ee-4608-82ad-620190193edf", "Subscription id should be specified")
	}

	stripeSubscriptionInstance, err := stripesubscription.LoadStripeSubscriptionByID(id)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// Validate ownership
	isAdmin := api.ValidateAdminRights(context) == nil
	isOwner := utils.InterfaceToString(stripeSubscriptionInstance.Get("visitor_id")) == visitor.GetCurrentVisitorID(context)

	if !isAdmin && !isOwner {
		context.SetResponseStatusForbidden()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "ed71c6e1-744b-41e9-9b4b-95a4e6c1b75b", "Subscription ownership could not be verified")
	}

	// Set stripe api key
	stripeKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAPIKey))
	if stripeKey == "" {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "7be27002-9d50-44a5-9b8b-11d22d757f11", "Stripe API key is empty")
	}
	stripe.Key = stripeKey

	stripeSub, err := sub.Cancel(stripeSubscriptionInstance.GetStripeSubscriptionID(), nil)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a1f8f680-1c9d-4f96-bfdd-95aedd9ca4b0", err.Error())
	}

	stripeSubscriptionInstance.Set("status", stripeSub.Status)
	stripeSubscriptionInstance.Save()

	email := stripeSubscriptionInstance.GetCustomerEmail()
	emailTemplate := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmailCancelTemplate))
	emailSubject := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathEmailCancelSubject))

	if emailTemplate == "" {
		emailTemplate = `Dear {{.Visitor.name}},
Your subscription was canceled`
	}
	if emailSubject == "" {
		emailSubject = "Subscription Cancelation"
	}
	templateMap := map[string]interface{}{
		"Visitor": map[string]interface{}{"name": stripeSubscriptionInstance.Get("customer_name")},
		"Site":    map[string]interface{}{"url": app.GetStorefrontURL("")},
	}
	emailToVisitor, err := utils.TextTemplate(emailTemplate, templateMap)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if err = app.SendMail(email, emailSubject, emailToVisitor); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return stripeSub, nil
}

// APIListVisitorSubscriptions returns a list of visitor's stripe subscriptions
func APIListVisitorSubscriptions(context api.InterfaceApplicationContext) (interface{}, error) {
	return nil, nil
}

// APIProcessStripeEvent listens to Stripe events and makes appropriate updates to subscriptions
func APIProcessStripeEvent(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// Set stripe api key
	stripeKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAPIKey))
	if stripeKey == "" {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "66b29232-5f49-4c2f-baa6-4db95b931bbf", "Stripe API key is empty")
	}
	stripe.Key = stripeKey

	// Get stripe event
	eventID := utils.InterfaceToString(requestData["id"])
	if eventID == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c4a15c78-6c95-400c-b516-67d65679ccb1", "Event ID should be scpecified")
	}
	evt, err := event.Get(eventID, nil)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2630a2ef-536f-4fff-a9e1-e6c678e05a9a", err.Error())
	}

	// Handle event
	switch evt.Type {
	case "customer.subscription.deleted":
		err = eventCancelHandler(evt)
	case "customer.subscription.updated":
		err = eventUpdateHandler(evt)
	case "invoice.payment_failed", "invoice.payment_succeeded":
		err = eventPaymentHandler(evt)
	default:
		err = nil
	}

	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}
