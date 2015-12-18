package subscription

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
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

	err = api.GetRestService().RegisterAPI("subscription/:subscriptionID/submit", api.ConstRESTOperationGet, APISubmitSubscription)
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

	return subscriptionModel.ToHashMap(), nil
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

	// delete operation
	err = subscriptionModel.Delete()
	if err != nil {
		return nil, env.ErrorDispatch(err)
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

	creationConfirmed := false
	if confirmValue, present := requestData["confirmed"]; present {
		creationConfirmed = utils.InterfaceToBool(confirmValue)
	}

	// validating basic input (name, email, addresses)
	customerName := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "name", "customer_name"))
	if customerName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fcfd4ed9-13e7-443f-a2e0-a62d0aa47518", "please specify customer name")
	}

	customerEmail := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "email", "customer_email"))
	if !utils.ValidEmailAddress(customerEmail) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e8b6c4cd-123a-4ec4-b413-55e66def1652", "customer email invalid")
	}

	shippingAddress, err := checkout.ValidateAddress(utils.InterfaceToMap(utils.GetFirstMapValue(requestData, "shipping_address")))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	billingAddress, err := checkout.ValidateAddress(utils.InterfaceToMap(utils.GetFirstMapValue(requestData, "billing_address")))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	subscriptionInstance, err := subscription.GetSubscriptionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// retrieving and validating given subscription date
	subscriptionDateValue := utils.GetFirstMapValue(requestData, "date", "action_date", "billing_date")
	if subscriptionDateValue == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "43873ddc-a817-4216-aa3c-9b004d96a539", "subscription Date can't be blank")
	}

	// retrieving and validating given subscription period
	subscriptionPeriodValue := utils.GetFirstMapValue(requestData, "period", "recurrence_period")
	if subscriptionPeriodValue == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cf33b877-97ab-4177-a529-3b1225c37fbd", "subscription Period can't be blank")
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
	subscriptionInstance.SetShippingAddress(shippingAddress)
	subscriptionInstance.SetBillingAddress(billingAddress)

	orderID, _ := requestData["order_id"]
	cartID, _ := requestData["cart_id"]

	// TODO: handle usage of transaction from order or use credit card proviede by visitor
	creditCardInstance, err := visitor.GetVisitorCardModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	creditCardData, present := requestData["credit_card"]
	if present {

		err = creditCardInstance.FromHashMap(utils.InterfaceToMap(creditCardData))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	// shipping method and payment method handling (no validation for is allowed - it require a checkout object?)
	// TODO: usage of shipping method
	//specifiedShippingMethod := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "shipping_method"))

	currentCart, err := cart.GetCurrentCart(context, true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentSession := context.GetSession()

	err = currentCart.ValidateCart()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cartID = currentCart.GetID()

	err = currentCart.Deactivate()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = currentCart.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentSession.Set(cart.ConstSessionKeyCurrentCart, nil)

	subscriptionCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := map[string]interface{}{
		"visitor_id":       visitorID,
		"order_id":         utils.InterfaceToString(orderID),
		"cart_id":          cartID,
		"email":            customerEmail,
		"name":             customerName,
		"shipping_address": shippingAddress,
		"billing_address":  billingAddress,
		"status":           ConstSubscriptionStatusSuspended,
		"action":           ConstSubscriptionActionUpdate,
	}

	// for orders that not have transaction by default we set action value to Update or Create
	// that means on subscription Date they will need to proceed checkout
	if creationConfirmed && orderID != nil {
		result["action"] = ConstSubscriptionActionSubmit
	}

	if orderID == nil && cartID != "" {
		result["action"] = ConstSubscriptionActionCreate
	}

	_, err = subscriptionCollection.Save(result)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return result, nil
}

// APIUpdateSubscription change status of subscription to suspended, and allow to change date and period
//   - subscription id should be specified in "subscriptionID" argument
func APIUpdateSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	subscriptionID := context.GetRequestArgument("subscriptionID")
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
//   - subscription id and new status should be specified in "subscriptionID" and "status" arguments
func APIUpdateSubscriptionStatus(context api.InterfaceApplicationContext) (interface{}, error) {

	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6845852a-e484-4f18-adcc-f8e166838c09", "you are not logined in")
	}

	// check request context
	//---------------------
	subscriptionID := context.GetRequestArgument("subscriptionID")
	if subscriptionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4e8f9873-9144-42ae-b119-d1e95bb1bbfd", "subscription id should be specified")
	}

	requestedStatus := context.GetRequestArgument("status")
	if requestedStatus != ConstSubscriptionStatusSuspended && requestedStatus != ConstSubscriptionStatusConfirmed && requestedStatus != ConstSubscriptionStatusCanceled {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3b7d17c3-c5fa-4369-a039-49bafec2fb9d", "new subscription status should be one of allowed")
	}

	// load subscription by id
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

	// check is current visitor was a creator of subscription or it's admin
	if api.ValidateAdminRights(context) != nil && visitorID != subscription["visitor_id"] {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5d438cbd-60a3-44af-838f-bddf4e19364e", "you are not logined in")
	}

	currentStatus := utils.InterfaceToString(subscription["status"])
	if currentStatus == requestedStatus {
		return "ok", nil
	}

	// in case subscription was canceled we would check it's date
	if currentStatus == ConstSubscriptionStatusCanceled {
		currentDate := utils.InterfaceToTime(subscription["action_date"])
		if currentDate.Before(nextAllowedCreationDate()) {
			subscription["action_date"] = nextCreationDate
		}
	}

	subscription["status"] = requestedStatus

	_, err = subscriptionCollection.Save(subscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APISubmitSubscription give current session new checkout and card from subscription  and try to proceed it
//   - subscription id should be specified in "subscriptionID" argument
func APISubmitSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	subscriptionID := context.GetRequestArgument("subscriptionID")
	if subscriptionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "027e7ef9-b202-475b-a242-02e2d0d74ce6", "subscription id should be specified")
	}

	subscriptionCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	subscriptionCollection.AddFilter("_id", "=", subscriptionID)

	dbRecords, err := subscriptionCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if len(dbRecords) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "aecec9b2-3b02-40cb-b163-39cb03b53252", "subscription not found")
	}

	subscription := utils.InterfaceToMap(dbRecords[0])

	if utils.InterfaceToString(subscription["status"]) != ConstSubscriptionStatusConfirmed {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "153ee2dd-3e3f-42ac-b669-1d15ec741547", "subscription not confirmed")
	}

	subscriptionDate := utils.InterfaceToTime(subscription["action_date"])
	subscriptionAction := utils.InterfaceToString(subscription["action"])

	currentDay := time.Now().Truncate(ConstTimeDay)

	// when someone try to submit subscription before available date (means submitting email wasn't sented yet)
	if currentDay.Before(subscriptionDate.Truncate(ConstTimeDay)) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "747ae177-a295-4029-b1dc-4abcce319d7b", "subscription can't be submited yet")
	}

	subscriptionOrderID := subscription["order_id"]
	subscriptionCartID := subscription["cart_id"]

	if subscriptionOrderID == nil && subscriptionCartID == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "780975f0-c24d-452c-ae43-cfbef64b9a1a", "this subscription can't be submited (no cart and order in)")
	}

	shippingAddress := utils.InterfaceToMap(subscription["shipping_address"])

	// obtain user current cart and checkout for future operations
	currentSession := context.GetSession()

	currentCheckout, err := checkout.GetCurrentCheckout(context, true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentCart := currentCheckout.GetCart()

	// Duplicating order and set new checkout and cart to current session with redirect to checkout if need to update info
	if utils.InterfaceToString(subscriptionOrderID) != "" && subscriptionAction != ConstSubscriptionActionCreate {
		orderModel, err := order.LoadOrderByID(utils.InterfaceToString(subscriptionOrderID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// update cart and checkout object for current session
		duplicateCheckout, err := orderModel.DuplicateOrder(nil)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		checkoutInstance, ok := duplicateCheckout.(checkout.InterfaceCheckout)
		if !ok {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3788a54b-6ef6-486f-9819-c85e34ff43c5", "order can't be typed")
		}

		err = checkoutInstance.Set("ShippingAddress", shippingAddress)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// rewrite current checkout and cart by newly created from duplicate order
		err = checkoutInstance.SetSession(currentSession)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		err = checkoutInstance.SetInfo("subscription", subscriptionID)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// rewrite current curt with duplicated
		err = currentCart.Deactivate()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		err = currentCart.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		duplicateCart := checkoutInstance.GetCart()

		err = duplicateCart.SetSessionID(currentSession.GetID())
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		err = duplicateCart.Activate()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		err = duplicateCart.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if subscription["action"] == ConstSubscriptionActionSubmit {
			_, err := checkoutInstance.Submit()
			if err != nil {
				currentSession.Set(cart.ConstSessionKeyCurrentCart, duplicateCart.GetID())
				currentSession.Set(checkout.ConstSessionKeyCurrentCheckout, checkoutInstance)
				return api.StructRestRedirect{Result: "ok", Location: app.GetStorefrontURL("checkout")}, env.ErrorDispatch(err)
			}
		} else {
			currentSession.Set(cart.ConstSessionKeyCurrentCart, duplicateCart.GetID())
			currentSession.Set(checkout.ConstSessionKeyCurrentCheckout, checkoutInstance)
			return api.StructRestRedirect{Result: "ok", Location: app.GetStorefrontURL("checkout")}, nil
		}

		// We need to set for user his subscription cart and add to checkout info subscription to handle it on success
	} else {
		err = currentCheckout.SetInfo("subscription", subscriptionID)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		err = currentCheckout.Set("ShippingAddress", shippingAddress)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		subscriptionCart, err := cart.LoadCartByID(utils.InterfaceToString(subscriptionCartID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		for _, cartItem := range currentCart.GetItems() {
			currentCart.RemoveItem(cartItem.GetIdx())
		}

		for _, cartItem := range subscriptionCart.GetItems() {
			_, err = currentCart.AddItem(cartItem.GetProductID(), cartItem.GetQty(), cartItem.GetOptions())
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}

		err = currentCart.Save()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		return api.StructRestRedirect{Result: "ok", Location: app.GetStorefrontURL("checkout")}, nil
	}

	// in case of instant checkout submit
	subscriptionNextDate := subscriptionDate.AddDate(0, utils.InterfaceToInt(subscription["period"]), 0)
	subscription["action_date"] = subscriptionNextDate
	subscription["last_submit"] = currentDay
	subscription["status"] = ConstSubscriptionStatusSuspended

	_, err = subscriptionCollection.Save(subscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}
