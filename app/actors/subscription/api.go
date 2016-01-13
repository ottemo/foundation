package subscription

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("subscriptions", api.ConstRESTOperationGet, APIListSubscriptions)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("subscription", api.ConstRESTOperationCreate, APICreateSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("subscriptional/checkout", api.ConstRESTOperationCreate, APICheckCheckoutSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("subscription/:subscriptionID", api.ConstRESTOperationGet, APIGetSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("subscription/:subscriptionID", api.ConstRESTOperationDelete, APIDeleteSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("subscription/:subscriptionID", api.ConstRESTOperationUpdate, APIUpdateSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("subscription/:subscriptionID/status/:status", api.ConstRESTOperationGet, APIUpdateSubscriptionStatus)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIListSubscriptions returns a list of subscriptions for visitor
//   - if "action" parameter is set to "count" result value will be just a number of list items
func APIListSubscriptions(context api.InterfaceApplicationContext) (interface{}, error) {

	sessionVisitorID := visitor.GetCurrentVisitorID(context)
	if sessionVisitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c73e39c9-dc23-463b-9792-a5d3f7e4d9dd", "You should log in first")
	}

	// if visitorID was specified - using this otherwise, taking current visitor
	visitorID := context.GetRequestArgument("visitorID")

	// check rights if it user we will search only for his subscriptions
	if err := api.ValidateAdminRights(context); err != nil {
		visitorID = sessionVisitorID
	}

	// list operation
	//---------------
	subscriptionCollectionModel, err := subscription.GetSubscriptionCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbCollection := subscriptionCollectionModel.GetDBCollection()
	if visitorID != "" {
		dbCollection.AddStaticFilter("visitor_id", "=", visitorID)
	}

	// filters handle
	models.ApplyFilters(context, dbCollection)

	// checking for a "count" request
	if context.GetRequestArgument("count") != "" {
		return dbCollection.Count()
	}

	// limit parameter handle
	subscriptionCollectionModel.ListLimit(models.GetListLimit(context))

	// extra parameter handle
	models.ApplyExtraAttributes(context, subscriptionCollectionModel)

	return subscriptionCollectionModel.List()
}

// APIGetSubscription return specified subscription information
//   - subscription id should be specified in "subscriptionID" argument
func APIGetSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	subscriptionID := context.GetRequestArgument("subscriptionID")
	if subscriptionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b626ec0a-a317-4b63-bd05-cc23932bdfe0", "subscription id should be specified")
	}

	subscriptionModel, err := subscription.LoadSubscriptionByID(subscriptionID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		if subscriptionModel.GetVisitorID() != visitor.GetCurrentVisitorID(context) {
			return nil, env.ErrorDispatch(err)
		}
	}

	result := subscriptionModel.ToHashMap()

	subscriptionCheckout, err := subscriptionModel.GetCheckout()
	if subscriptionCheckout != nil {
		result["checkout"] = subscriptionCheckout.ToHashMap()
	} else {
		result["checkout_error"] = err
	}

	return result, nil
}

// APIDeleteSubscription deletes existing purchase order
//   - subscription id should be specified in "subscriptionID" argument
func APIDeleteSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	subscriptionID := context.GetRequestArgument("subscriptionID")
	if subscriptionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "67bedbe8-7426-437b-9dbc-4840f13e619e", "subscription id should be specified")
	}

	subscriptionInstance, err := subscription.LoadSubscriptionByID(subscriptionID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		if subscriptionInstance.GetVisitorID() != visitor.GetCurrentVisitorID(context) {
			return nil, env.ErrorDispatch(err)
		}
	}

	// delete operation
	err = subscriptionInstance.Delete()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APICheckCheckoutSubscription provide check is current checkout allows to create new subscription
func APICheckCheckoutSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	// check visitor to be registered
	visitorID := visitor.GetCurrentVisitorID(context)
	if api.ValidateAdminRights(context) != nil && visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e6109c04-e35a-4a90-9593-4cc1f141a358", "you are not logined in")
	}

	currentCheckout, err := checkout.GetCurrentCheckout(context, false)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	customer := currentCheckout.GetVisitor()

	customerEmail := customer.GetEmail()
	if !utils.ValidEmailAddress(customerEmail) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "7c78ba76-647f-4e16-abd3-5e2d5afeb8cf", "Customer email invalid")
	}

	shippingAddress := currentCheckout.GetShippingAddress()
	if shippingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e8b6c4cd-123a-4ec4-b413-55e66def1652", "No shipping address")
	}

	billingAddress := currentCheckout.GetBillingAddress()
	if billingAddress == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b5ecb475-cf90-4d56-99e9-2dcc5a772c54", "No billing address")
	}

	currentCart := currentCheckout.GetCart()
	if currentCart.GetSubtotal() <= 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "054459d2-0b6b-4526-b0a7-92e7dfce43b4", "Cart with items should be provided")
	}

	paymentMethod := currentCheckout.GetPaymentMethod()
	if paymentMethod == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2c50d7c5-9f14-4ca2-8c79-1a14f068fe75", "Payment method not set")
	}

	if !paymentMethod.IsAllowed(currentCheckout) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8d2a66d2-5d47-43e3-bfaa-17ea0e796667", "Payment method not allowed")
	}

	if !paymentMethod.IsTokenable(currentCheckout) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "dd38d1d7-3921-4da1-8a2a-420067341511", "Payment method not support subsciptions")
	}

	// checking shipping method an shipping rates
	shippingMethod := currentCheckout.GetShippingMethod()
	if shippingMethod == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b847fd19-81e6-44fa-946b-fc4c7c45a38b", "Shipping method not set")
	}

	shippingRate := currentCheckout.GetShippingRate()
	if shippingRate == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "962291f9-3673-4ec6-b95a-a92a80c11eb4", "Shipping rate not set")
	}

	return "ok", nil
}

// APICreateSubscription provide mechanism to create new subscription
// products are getted from specified cart or order or current cart
func APICreateSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	// check visitor rights
	// TODO: is there any + for admin?
	visitorID := visitor.GetCurrentVisitorID(context)
	if api.ValidateAdminRights(context) != nil && visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e6109c04-e35a-4a90-9593-4cc1f141a358", "you are not logined in")
	}

	// check request context
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// validating basic input (name, email, addresses)
	customerName := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "name", "customer_name"))
	if customerName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fcfd4ed9-13e7-443f-a2e0-a62d0aa47518", "Please specify customer name")
	}

	customerEmail := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "email", "customer_email"))
	if !utils.ValidEmailAddress(customerEmail) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e8b6c4cd-123a-4ec4-b413-55e66def1652", "Customer email invalid")
	}

	shippingAddress, err := checkout.ValidateAddress(utils.InterfaceToMap(utils.GetFirstMapValue(requestData, "shipping_address")))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	billingAddress, err := checkout.ValidateAddress(utils.InterfaceToMap(utils.GetFirstMapValue(requestData, "billing_address")))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cartID := utils.InterfaceToString(requestData["cart_id"])
	if cartID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "054459d2-0b6b-4526-b0a7-92e7dfce43b4", "Cart with items should be provided")
	}

	creditCardID := utils.InterfaceToString(requestData["credit_card_id"])
	if creditCardID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "39d5f023-f3f5-44d5-8d3b-2225af8ae0d7", "Saved credit card should be provided")
	}

	specifiedShippingMethod := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "shipping_method", "shipppingMethod"))
	specifiedShippingMethodRate := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "shipppingRate", "shipping_rate"))

	if specifiedShippingMethod == "" || specifiedShippingMethodRate == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "882141c2-634e-416b-9b5d-4b8cf9bcecb7", "Shipping method and rates can't be blank")
	}

	// retrieving and validating given subscription date
	subscriptionDateValue := utils.GetFirstMapValue(requestData, "date", "action_date", "billing_date")
	if subscriptionDateValue == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "43873ddc-a817-4216-aa3c-9b004d96a539", "Subscription Date can't be blank")
	}

	// retrieving and validating given subscription period
	subscriptionPeriodValue := utils.GetFirstMapValue(requestData, "period", "recurrence_period")
	if subscriptionPeriodValue == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cf33b877-97ab-4177-a529-3b1225c37fbd", "Subscription Period can't be blank")
	}

	subscriptionInstance, err := subscription.GetSubscriptionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = subscriptionInstance.SetPeriod(utils.InterfaceToInt(subscriptionPeriodValue))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = subscriptionInstance.SetActionDate(utils.InterfaceToTime(subscriptionDateValue))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	subscriptionInstance.Set("name", customerName)
	subscriptionInstance.Set("email", customerEmail)
	subscriptionInstance.Set("visitor_id", visitorID)
	subscriptionInstance.SetShippingAddress(shippingAddress)
	subscriptionInstance.SetBillingAddress(billingAddress)

	// TODO: handle usage of transaction from order or use credit card provided by visitor
	creditCardInstance, err := visitor.LoadVisitorCardByID(creditCardID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if creditCardInstance.GetVisitorID() != visitorID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4ed389a7-a9d6-40d1-9ff7-6128a95f3979", "Credit Card not found")
	}

	subscriptionInstance.SetCreditCard(creditCardInstance)

	providedCart, err := cart.LoadCartByID(cartID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if providedCart.GetVisitorID() != visitorID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "7646b717-08e3-4c23-a713-b3ff204b1cf0", "Given Cart not found")
	}

	err = providedCart.ValidateCart()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// saving current active cart as a new and make inactive
	if providedCart.IsActive() {

		providedCart.SetID("")
		err = providedCart.Deactivate()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		err = providedCart.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	subscriptionInstance.Set("cart_id", providedCart.GetID())

	// checking shipping method an shipping rates
	shippingMethod := checkout.GetShippingMethodByCode(specifiedShippingMethod)
	if shippingMethod == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b847fd19-81e6-44fa-946b-fc4c7c45a38b", "Shipping method not found")
	}

	checkoutInstance, _ := subscriptionInstance.GetCheckout()

	if !shippingMethod.IsAllowed(checkoutInstance) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "52815f06-8dea-4c46-b6a4-b1c00d07c7a0", "Shipping method not allowed")
	}

	subscriptionInstance.SetShippingMethod(shippingMethod)

	ratesFound := false
	for _, shippingRate := range shippingMethod.GetRates(checkoutInstance) {
		if shippingRate.Code == specifiedShippingMethodRate {
			ratesFound = true
			subscriptionInstance.SetShippingRate(shippingRate)
		}
	}

	if !ratesFound {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "446d2a21-bcb2-417b-a96c-5011a14289d8", "Shipping rates were not found")
	}

	subscriptionInstance.SetStatus(ConstSubscriptionStatusSuspended)

	err = subscriptionInstance.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return subscriptionInstance.ToHashMap(), nil
}

// APIUpdateSubscription change status of subscription to suspended, and allow to change date and period
//   - subscription id should be specified in "subscriptionID" argument
// TODO: update actions
func APIUpdateSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	subscriptionID := context.GetRequestArgument("subscription_id")
	if subscriptionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4e8f9873-9144-42ae-b119-d1e95bb1bbfd", "subscription id should be specified")
	}

	// check request context
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// TODO: is it allowed to update date?
	// validate new params for subscription
	requestedParams := make(map[string]interface{})
	subscriptionDateValue := utils.GetFirstMapValue(requestData, "date", "action_date", "billing_date")

	if subscriptionDateValue != nil {
		subscriptionDate := utils.InterfaceToTime(subscriptionDateValue)

		// here is requirements for subscription date and day
		if err := validateSubscriptionDate(subscriptionDate); err != nil {
			return nil, env.ErrorDispatch(err)
		}

		requestedParams["action_date"] = subscriptionDate
	}

	subscriptionPeriodValue := utils.GetFirstMapValue(requestData, "period", "recurrence_period", "recurring")

	if subscriptionPeriodValue != nil {
		subscriptionPeriod := utils.InterfaceToInt(subscriptionPeriodValue)

		if subscriptionPeriod < 1 || subscriptionPeriod > 3 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "85f539fa-89fe-4ad8-b171-3b66910bad3f", "subscription recurrence period cannot be only 1, 2 or 3 monthes")
		}

		requestedParams["period"] = subscriptionPeriod
	}

	subscriptionEmailValue := utils.GetFirstMapValue(requestData, "email", "customer_email")

	if subscriptionEmailValue != nil {
		subscriptionEmail := utils.InterfaceToString(subscriptionEmailValue)
		if !utils.ValidEmailAddress(subscriptionEmail) {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "bfb0b402-6e89-43ee-8c69-ae104f389e70", "given email is not valid")
		}

		requestedParams["email"] = subscriptionEmail
	}

	subscriptionNameValue := utils.GetFirstMapValue(requestData, "name", "customer_name")
	if subscriptionNameValue != nil {
		requestedParams["name"] = utils.InterfaceToString(subscriptionNameValue)
	}

	subscriptionShippingAddressValue := utils.GetFirstMapValue(requestData, "address", "shipping_address")
	if subscriptionShippingAddressValue != nil {
		subscriptionShippingAddress := utils.InterfaceToMap(subscriptionShippingAddressValue)

		requiredAddressFields := []string{"first_name", "last_name", "address_line1", "country", "city", "zip_code"}
		if !utils.KeysInMapAndNotBlank(subscriptionShippingAddress, requiredAddressFields) {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "08a18398-a836-4a2b-a63a-4fac077407bb", "shipping address fields not all")
		}

		requestedParams["shipping_address"] = subscriptionShippingAddress
	}

	subscriptionCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	subscriptionCollection.AddFilter("_id", "=", subscriptionID)

	dbRecords, err := subscriptionCollection.Load()

	if len(dbRecords) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "59e4ab86-7726-4e3d-bec8-7ef5bf0ebbbf", "subscription not found")
	}

	subscription := utils.InterfaceToMap(dbRecords[0])

	visitorID := visitor.GetCurrentVisitorID(context)
	if api.ValidateAdminRights(context) != nil && visitorID != subscription["visitor_id"] {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cec5f9c7-1034-4c19-b4b0-251052255570", "you are not logined in")
	}

	for key, value := range requestedParams {
		subscription[key] = value
	}

	_, err = subscriptionCollection.Save(subscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return subscription, nil
}

// APIUpdateSubscriptionStatus change status of subscription to suspended, confirmed or canceled
//   - subscription id and new status should be specified in "subscription_id" and "status" arguments
func APIUpdateSubscriptionStatus(context api.InterfaceApplicationContext) (interface{}, error) {

	visitorID := visitor.GetCurrentVisitorID(context)

	// check request context
	//---------------------
	subscriptionID := context.GetRequestArgument("subscription_id")
	if subscriptionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4e8f9873-9144-42ae-b119-d1e95bb1bbfd", "subscription id should be specified")
	}

	// load subscription by id
	subscriptionInstance, err := subscription.LoadSubscriptionByID(subscriptionID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check is current visitor was a creator of subscription or it's admin
	if api.ValidateAdminRights(context) != nil && visitorID != subscriptionInstance.GetVisitorID() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5d438cbd-60a3-44af-838f-bddf4e19364e", "subscription not found")
	}

	requestedStatus := context.GetRequestArgument("status")
	err = subscriptionInstance.SetStatus(requestedStatus)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}
