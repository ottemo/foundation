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
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/coupon"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/event"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/sub"
	"fmt"
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
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	stripeSubscriptionInstance, err := stripesubscription.GetStripeSubscriptionModel()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// Set stripe api key
	if err = setStripeAPIKey(); err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// Process visitor information
	//----------------------------
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		context.SetResponseStatusForbidden()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "287e280f-3567-4935-9cbf-0f2f1afd149e", "You should log in to subscribe")
	}
	visitorInstance, err := visitor.LoadVisitorByID(visitorID)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}
	stripeSubscriptionInstance.Set("visitor_id", visitorInstance.GetID())
	stripeSubscriptionInstance.Set("customer_name", visitorInstance.GetFullName())
	stripeSubscriptionInstance.Set("customer_email", visitorInstance.GetEmail())

	// Process shipping address
	//----------------------------
	shippingAddress, err := checkout.ValidateAddress(requestData["shipping_address"])
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}
	stripeSubscriptionInstance.Set("shipping_address", shippingAddress.ToHashMap())

	// Process billing address
	//----------------------------
	billingAddress, err := checkout.ValidateAddress(requestData["billing_address"])
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}
	stripeSubscriptionInstance.Set("billing_address", billingAddress.ToHashMap())

	// Process subscription plan
	//----------------------------
	planID := utils.InterfaceToString(requestData["plan_id"])
	if planID == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b153941b-27d3-4653-bc9f-61683b6047a9", "Subscription plan ID should be specified")
	}
	stripePlan, err := plan.Get(planID, nil)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "ba0c64db-e723-4c0f-a646-d06c7ca8f17c", err.Error())
	}
	stripeSubscriptionInstance.Set("description", stripePlan.Name)
	// Set price as plan amount
	price := float64(stripePlan.Amount / 100.)

	// Validate credit card token
	//----------------------------
	ccToken := utils.InterfaceToString(requestData["cc_token"])
	if ccToken == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a20e1235-c23c-40f0-ae5a-5475abf3427e", "Credit card token should be specified")
	}

	// Prepare parameters for Stripe subscription
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

	// Process coupon and adjust price
	if utils.InterfaceToString(requestData["coupon"]) != "" {
		stripeCoupon, err := getCoupon(utils.InterfaceToString(requestData["coupon"]))
		if err != nil {
			context.SetResponseStatusInternalServerError()
			return nil, env.ErrorDispatch(err)
		}

		if stripeCoupon.Valid != true {
			context.SetResponseStatusBadRequest()
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "10787e67-fa86-4728-8b1d-ca0f95c4c81c", "Coupon is not valid")
		}

		stripeSubscriptionInstance.Set("stripe_coupon", stripeCoupon.ID)
		//TODO: adjust price here
	}
	stripeSubscriptionInstance.Set("price", price)

	// Set notify on renewing flag
	if utils.InterfaceToBool(requestData["notify_on_renew"]) == true {
		stripeSubscriptionInstance.Set("notify_renew", true)
	} else {
		stripeSubscriptionInstance.Set("notify_renew", false)
	}
	stripeSubscriptionInstance.Set("renew_notified", false)

	// Set gift information
	if utils.InterfaceToBool(requestData["is_gift"]) == true {
		stripeSubscriptionInstance.Set("info", map[string]interface{}{"Gift": true})
	}

	stripeCustomer, err := customer.New(customerParams)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f3825627-1dd1-475f-a6f7-d3acaebe590e", err.Error())
	}
	stripeSub := stripeCustomer.Subs.Values[0]
	stripeSubscriptionInstance.Set("status", stripeSub.Status)
	stripeSubscriptionInstance.Set("stripe_customer_id", stripeCustomer.ID)
	stripeSubscriptionInstance.Set("stripe_subscription_id", stripeSub.ID)
	stripeSubscriptionInstance.Set("period_end", stripeSub.PeriodEnd)

	err = stripeSubscriptionInstance.Save()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// TODO: send email

	return stripeSubscriptionInstance.ToHashMap(), nil
}

// APIGetCoupon returns stripe coupon
func APIGetCoupon(context api.InterfaceApplicationContext) (interface{}, error) {
	// Validate params
	id := context.GetRequestArgument("id")
	if id == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2630545b-ae5f-4a62-8082-a0a457db0b57", "Coupon id should be specified")
	}

	// Set stripe api key
	if err := setStripeAPIKey(); err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	stripeCoupon, err := getCoupon(id)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	if stripeCoupon.Valid != true {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "372f6fdc-2194-4dcc-a668-9b2b69d82eac", "Coupon "+id+" can't be applied")
	}

	return stripeCoupon, nil
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
	// Validate request context
	id := context.GetRequestArgument("id")
	if id == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "945b9d1e-ab0a-4853-aaaf-2ae2a4478d9a", "Subscription id should be specified")
	}

	stripeSubscriptionModel, err := stripesubscription.LoadStripeSubscriptionByID(id)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	result := stripeSubscriptionModel.ToHashMap()

	return result, nil
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

	if err = setStripeAPIKey(); err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

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
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	if err = app.SendMail(email, emailSubject, emailToVisitor); err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return stripeSub, nil
}

// APIListVisitorSubscriptions returns a list of visitor's stripe subscriptions
func APIListVisitorSubscriptions(context api.InterfaceApplicationContext) (interface{}, error) {
	// Validate visitor
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "bc97068f-5666-4533-bea8-7a31d912cf83", "You should log in first")
	}

	// for showing subscriptions to a visitor, request is specific so handle it in different way from default List
	stripeSubscriptionCollection, err := stripesubscription.GetStripeSubscriptionCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbCollection := stripeSubscriptionCollection.GetDBCollection()
	dbCollection.AddStaticFilter("visitor_id", "=", visitorID)
	dbCollection.AddStaticFilter("status", "=", "active")
	models.ApplyFilters(context, dbCollection)

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return dbCollection.Count()
	}

	// limit parameter handle
	dbCollection.SetLimit(models.GetListLimit(context))

	subscriptions := stripeSubscriptionCollection.ListSubscriptions()
	var result []map[string]interface{}

	for _, subscriptionItem := range subscriptions {
		result = append(result, subscriptionItem.ToHashMap())
	}

	return result, nil
}

// APIProcessStripeEvent listens to Stripe events and makes appropriate updates to subscriptions
func APIProcessStripeEvent(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}
	fmt.Println(requestData)

	// Set stripe api key
	if err = setStripeAPIKey(); err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	fmt.Println(requestData)
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

// SetStripeAPIKey gets Stripe API key from config and sets it to the stripe
func setStripeAPIKey() error {
	if stripe.Key != "" {
		return nil
	}

	stripeKey := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathAPIKey))
	if stripeKey == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "66b29232-5f49-4c2f-baa6-4db95b931bbf", "Stripe API key is empty")
	}
	stripe.Key = stripeKey
	return nil
}

// getCoupon obtain coupon from Stripe by coupon id
func getCoupon(id string) (*stripe.Coupon, error) {
	stripeCoupon, err := coupon.Get(id, nil)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "91e36260-2b84-433b-8529-6ce58bf591e1", err.Error())
	}

	return stripeCoupon, nil
}
