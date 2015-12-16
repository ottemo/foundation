package subscription

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
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

	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "21762af1-e352-4ef6-82b3-2fe9d55d6c36", "you are not logined in")
	}

	// making database request
	subscriptionCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	subscriptionCollection.AddFilter("visitor_id", "=", visitorID)

	dbRecords, err := subscriptionCollection.Load()

	return dbRecords, env.ErrorDispatch(err)
}

// APIGetSubscription return specified subscription information
//   - subscription id should be specified in "subscriptionID" argument
func APIGetSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	subscriptionID := context.GetRequestArgument("subscriptionID")
	if subscriptionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "99b0a49b-9fe4-4f64-9879-bf5a45ff5ac7", "subscription id should be specified")
	}

	subscriptionCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	subscriptionCollection.AddFilter("_id", "=", subscriptionID)

	dbRecords, err := subscriptionCollection.Load()

	if len(dbRecords) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d724cdc3-5bb7-494b-9a8a-952fdc311bd0", "subscription not found")
	}

	return dbRecords[0], nil
}

// APIDeleteSubscription deletes existing purchase order
//   - subscription id should be specified in "subscriptionID" argument
func APIDeleteSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5d438cbd-60a3-44af-838f-bddf4e19364e", "you are not logined in")
	}
	// check request context
	//---------------------
	subscriptionID := context.GetRequestArgument("subscriptionID")
	if subscriptionID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "67bedbe8-7426-437b-9dbc-4840f13e619e", "subscription id should be specified")
	}

	subscriptionCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	subscriptionCollection.AddFilter("_id", "=", subscriptionID)

	dbRecords, err := subscriptionCollection.Load()

	if len(dbRecords) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6c9559d5-c0fe-4fa1-a07b-4e7b6ac1dad6", "subscription not found")
	}

	if api.ValidateAdminRights(context) == nil || visitorID == dbRecords[0]["visitor_id"] {
		subscriptionCollection.DeleteByID(subscriptionID)
	} else {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5d438cbd-60a3-44af-838f-bddf4e19364e", "you are not logined in")
	}

	return "ok", nil
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

	// validate new params for subscription
	requestedParams := make(map[string]interface{})
	subscriptionDateValue := utils.GetFirstMapValue(requestData, "date", "action_date", "billing_date")

	if subscriptionDateValue != nil {
		subscriptionDate := utils.InterfaceToTime(subscriptionDateValue)
		nextAllowedDate := nextAllowedCreationDate()
		// here is requirements for subscription date and day
		if subscriptionDate.Before(nextAllowedDate) {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c4881529-8b05-4a16-8cd4-6c79d0d79856", "subscription Date cannot be earlier then "+utils.InterfaceToString(nextAllowedDate))
		}

		if subscriptionDate.Day() != 15 && subscriptionDate.Day() != 1 {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "29c73d2f-0c85-4906-95b7-4812542e33a1", "schedule for either the 1st of the month or the 15th of the month")
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

// APICreateSubscription provide mechanism to create new subscription
func APICreateSubscription(context api.InterfaceApplicationContext) (interface{}, error) {

	// check visitor rights
	visitorID := visitor.GetCurrentVisitorID(context)
	if api.ValidateAdminRights(context) != nil && visitorID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e6109c04-e35a-4a90-9593-4cc1f141a358", "you are not logined in")
	}

	result := make(map[string]interface{})
	submittableOrder := false

	// check request context
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	subscriptionDateValue := utils.GetFirstMapValue(requestData, "date", "action_date", "billing_date")
	if subscriptionDateValue == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "43873ddc-a817-4216-aa3c-9b004d96a539", "subscription Date can't be blank")
	}

	subscriptionDate := utils.InterfaceToTime(subscriptionDateValue)
	nextAllowedDate := nextAllowedCreationDate()

	// here is requirements for subscription date and day
	if subscriptionDate.Before(nextAllowedDate) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c4881529-8b05-4a16-8cd4-6c79d0d79856", "subscription Date cannot be earlier then "+utils.InterfaceToString(nextAllowedDate))
	}

	if subscriptionDate.Day() != 15 && subscriptionDate.Day() != 1 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "29c73d2f-0c85-4906-95b7-4812542e33a1", "schedule for either the 1st of the month or the 15th of the month")
	}

	subscriptionPeriod := utils.GetFirstMapValue(requestData, "period", "recurrence_period", "recurring")
	if subscriptionPeriod == nil || utils.InterfaceToInt(subscriptionPeriod) < 1 {
		subscriptionPeriod = 1
	}

	if utils.InterfaceToInt(subscriptionPeriod) > 3 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "85f539fa-89fe-4ad8-b171-3b66910bad3f", "subscription recurrence period cannot be more than 3 month")
	}

	orderID, orderPresent := requestData["orderID"]
	cartID := ""

	customerName := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "name", "customer_name"))
	customerEmail := utils.InterfaceToString(utils.GetFirstMapValue(requestData, "email", "customer_email"))
	shippingAddress := utils.InterfaceToMap(utils.GetFirstMapValue(requestData, "address", "shipping_address"))

	currentCart, err := cart.GetCurrentCart(context, true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	currentSession := context.GetSession()

	if (!orderPresent || utils.InterfaceToString(orderID) == "") && len(currentCart.GetItems()) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d0484d9d-cb6d-48ed-be5f-f77fe19c6dca", "No items in cart or no order for subscription are specified")
	}

	// try to create new subscription with existing order
	if orderPresent && utils.InterfaceToString(orderID) != "" {

		orderModel, err := order.LoadOrderByID(utils.InterfaceToString(orderID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		orderVisitorID := utils.InterfaceToString(orderModel.Get("visitor_id"))
		cartID = utils.InterfaceToString(orderModel.Get("cart_id"))

		customerName = utils.InterfaceToString(orderModel.Get("customer_name"))
		customerEmail = utils.InterfaceToString(orderModel.Get("customer_email"))
		shippingAddress = orderModel.GetShippingAddress().ToHashMap()

		if api.ValidateAdminRights(context) != nil && visitorID != orderVisitorID {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4916bf20-e053-472e-98e1-bb28b7c867a1", "you are trying to use vicarious order")
		}

		// update cart and checkout object for current session
		duplicateCheckout, err := orderModel.DuplicateOrder(nil)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		checkoutInstance, ok := duplicateCheckout.(checkout.InterfaceCheckout)
		if !ok {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "946c3598-53b4-4dad-9d6f-23bf1ed6440f", "order can't be typed")
		}

		duplicateCart := checkoutInstance.GetCart()

		err = duplicateCart.Delete()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// check order for possibility to proceed automatically
		if paymentInfo := utils.InterfaceToMap(orderModel.Get("payment_info")); paymentInfo != nil {
			if _, present := paymentInfo["transactionID"]; present {
				submittableOrder = true
			}
		}
	} else {
		// Creation of subscription from current cart and given additional parameters
		if !utils.ValidEmailAddress(customerEmail) {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2a2b8eef-2168-492b-a59d-70d921005daf", "email address not valid")
		}

		requiredAddressFields := []string{"first_name", "last_name", "address_line1", "country", "city", "zip_code"}
		if !utils.KeysInMapAndNotBlank(shippingAddress, requiredAddressFields) {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "08a18398-a836-4a2b-a63a-4fac077407bb", "shipping address fields not all")
		}

		customerName = utils.InterfaceToString(shippingAddress["first_name"]) + " " + utils.InterfaceToString(shippingAddress["last_name"])

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
	}

	subscriptionCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result = map[string]interface{}{
		"visitor_id":       visitorID,
		"order_id":         utils.InterfaceToString(orderID),
		"cart_id":          cartID,
		"email":            customerEmail,
		"name":             customerName,
		"shipping_address": shippingAddress,
		"action_date":      subscriptionDate,
		"period":           utils.InterfaceToInt(subscriptionPeriod),
		"status":           ConstSubscriptionStatusSuspended,
		"action":           ConstSubscriptionActionUpdate,
	}

	// for orders that not have transaction by default we set action value to Update or Create
	// that means on subscription Date they will need to proceed checkout
	if submittableOrder && orderID != nil {
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
