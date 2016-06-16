package order

import (
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/actors/discount/giftcard"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	// Admin
	service.GET("orders/attributes", api.IsAdmin(APIListOrderAttributes))
	service.GET("orders", api.IsAdmin(APIListOrders))

	service.GET("order/:orderID", api.IsAdmin(APIGetOrder))
	service.PUT("order/:orderID", api.IsAdmin(APIUpdateOrder))
	service.DELETE("order/:orderID", api.IsAdmin(APIDeleteOrder))
	service.GET("order/:orderID/emailShipStatus", api.IsAdmin(APIEmailShipStatus))
	service.GET("order/:orderID/emailOrderConfirmation", api.IsAdmin(APIEmailOrderConfirmation))

	// Public
	service.GET("visit/orders", APIGetVisitOrders)
	service.GET("visit/order/:orderID", APIGetVisitOrder)

	return nil
}

// APIListOrderAttributes returns a list of purchase order attributes
func APIListOrderAttributes(context api.InterfaceApplicationContext) (interface{}, error) {

	orderModel, err := order.GetOrderModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return orderModel.GetAttributesInfo(), nil
}

// APIListOrders returns a list of existing purchase orders
//   - if "action" parameter is set to "count" result value will be just a number of list items
func APIListOrders(context api.InterfaceApplicationContext) (interface{}, error) {

	// taking orders collection model
	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// filters handle
	models.ApplyFilters(context, orderCollectionModel.GetDBCollection())

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return orderCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	orderCollectionModel.ListLimit(models.GetListLimit(context))

	// extra parameter handle
	models.ApplyExtraAttributes(context, orderCollectionModel)

	return orderCollectionModel.List()
}

// APIGetOrder return specified purchase order information
//   - order id should be specified in "orderID" argument
func APIGetOrder(context api.InterfaceApplicationContext) (interface{}, error) {

	// pull order id off context
	orderModel, err := orderByID(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	result := orderModel.ToHashMap()
	if notes, present := utils.InterfaceToMap(result["shipping_info"])["notes"]; present {
		utils.InterfaceToMap(result["shipping_address"])["notes"] = notes
	}

	result["items"] = orderModel.GetItems()
	return result, nil
}

// APIEmailShipStatus will send the visitor a shipping confirmation email
// - order id should be specified in "orderID" argument
func APIEmailShipStatus(context api.InterfaceApplicationContext) (interface{}, error) {

	orderModel, err := orderByID(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = orderModel.SendShippingStatusUpdateEmail()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return "Email Sent", nil
}

// APIUpdateOrder update existing purchase order
//   - order id should be specified in "orderID" argument
func APIUpdateOrder(context api.InterfaceApplicationContext) (interface{}, error) {

	orderModel, err := orderByID(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// update the order data from request
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		orderModel.Set(attribute, value)
	}

	orderModel.Save()

	return orderModel.ToHashMap(), nil
}

// APIDeleteOrder deletes existing purchase order
//   - order id should be specified in "orderID" argument
func APIDeleteOrder(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	orderID := context.GetRequestArgument("orderID")
	if orderID == "" {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fc3011c7-e58c-4433-b9b0-881a7ba005cf", "order id should be specified")
	}

	// grab order ID off request
	orderModel, err := order.GetOrderModelAndSetID(orderID)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	orderModel.Delete()
	return "Order Deleted", nil
}

// APIGetVisitOrder returns current visitor order details for specified order
//   - orderID should be specified in arguments
func APIGetVisitOrder(context api.InterfaceApplicationContext) (interface{}, error) {

	orderModel, err := orderByID(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// allow anonymous visitors through if the session id matches
	if utils.InterfaceToString(orderModel.Get("session_id")) != context.GetSession().GetID() {
		// force anonymous visitors to log in if their session id does not match the one on the order
		visitorID := visitor.GetCurrentVisitorID(context)
		if visitorID == "" {
			return "No Visitor ID found, unable to process order request. Please log in first.", nil
		} else if utils.InterfaceToString(orderModel.Get("visitor_id")) != visitorID {
			context.SetResponseStatusBadRequest()
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c5ca1fdb-7008-4a1c-a168-9df544df9825", "There is a mis-match between the current Visitor ID and the Visitor ID on the order.")
		}
	}

	result := orderModel.ToHashMap()
	result["items"] = orderModel.GetItems()

	return result, nil
}

// APIGetVisitOrders returns list of orders related to current visitor
//   - visitorID is required, visitor must be logged in
func APIGetVisitOrders(context api.InterfaceApplicationContext) (interface{}, error) {

	// list operation
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return "No Visitor ID found, unable to process request.  Please log in first.", nil
	}

	orderCollection, err := order.GetOrderCollectionModel()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	err = orderCollection.ListFilterAdd("visitor_id", "=", visitorID)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	// We only return orders that are in these two states
	statusFilter := [2]string{order.ConstOrderStatusProcessed, order.ConstOrderStatusCompleted}
	orderCollection.GetDBCollection().AddFilter("status", "in", statusFilter)

	descending := true
	orderCollection.GetDBCollection().AddSort("created_at", descending)

	// filters handle
	models.ApplyFilters(context, orderCollection.GetDBCollection())

	// extra parameter handle
	models.ApplyExtraAttributes(context, orderCollection)

	result, err := orderCollection.List()

	return result, env.ErrorDispatch(err)
}

// APIEmailOrderConfirmation will send out an order confirmation email to the visitor specficied in the orderID
//   - orderID must be passed as a request argument
func APIEmailOrderConfirmation(context api.InterfaceApplicationContext) (interface{}, error) {

	// load the email template
	email := utils.InterfaceToString(env.ConfigGetValue(checkout.ConstConfigPathConfirmationEmail))

	// load the order
	orderModel, err := orderByID(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// set visitor
	visitor := make(map[string]interface{})
	visitor["first_name"] = orderModel.Get("customer_name")
	visitor["email"] = orderModel.Get("customer_email")

	order := hydrateOrder(orderModel)

	// set store url
	customInfo := make(map[string]interface{})
	customInfo["base_storefront_url"] = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStorefrontURL))

	// build the email
	confirmationEmail, err := utils.TextTemplate(email, map[string]interface{}{
		"Order":   order,
		"Visitor": visitor,
		"Info":    customInfo,
	})
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return "failure", env.ErrorDispatch(err)
	}

	// send the email
	emailAddress := utils.InterfaceToString(visitor["email"])
	err = app.SendMail(emailAddress, "Order confirmation", confirmationEmail)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return "failure", env.ErrorDispatch(err)
	}

	return "Email Sent", nil
}

// orderByID will load the order from the database using the orderID from the application context
func orderByID(context api.InterfaceApplicationContext) (order.InterfaceOrder, error) {

	// load orderID
	orderID := context.GetRequestArgument("orderID")
	if orderID == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fc3011c7-e58c-4433-b9b0-881a7ba005cf", "order id should be specified")
	}

	orderModel, err := order.LoadOrderByID(orderID)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	return orderModel, nil
}

func hydrateOrder(orderModel order.InterfaceOrder) map[string]interface{} {

	var orderItems []map[string]interface{}

	order := orderModel.ToHashMap()
	timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
	giftCardSku := utils.InterfaceToString(env.ConfigGetValue(giftcard.ConstConfigPathGiftCardSKU))

	// convert date of order creation to store time zone
	if date, present := order["created_at"]; present {
		convertedDate, _ := utils.MakeTZTime(utils.InterfaceToTime(date), timeZone)
		if !utils.IsZeroTime(convertedDate) {
			order["created_at"] = convertedDate
		}
	}

	// load items into the order
	for _, item := range orderModel.GetItems() {
		options := make(map[string]interface{})

		for optionName, optionKeys := range item.GetOptions() {
			optionMap := utils.InterfaceToMap(optionKeys)

			if strings.Contains(item.GetSku(), giftCardSku) {

				// if we have a giftcard date, localize the date
				options[optionName] = ifGiftcardSetLocalDate(item.GetSku(), timeZone, optionMap)
			} else {

				// if not a date, just set the option value
				options[optionName] = optionMap["value"]
			}
		}

		orderItems = append(orderItems, map[string]interface{}{
			"name":    item.GetName(),
			"options": options,
			"sku":     item.GetSku(),
			"qty":     item.GetQty(),
			"price":   item.GetPrice()})
	}

	order["items"] = orderItems
	order["payment_method_title"] = orderModel.GetPaymentMethod()
	order["shipping_method_title"] = orderModel.GetShippingMethod()

	return order
}

func ifGiftcardSetLocalDate(optionName, timeZone string, optionMap map[string]interface{}) string {

	// make sure we are looking for a date
	if utils.IsAmongStr(optionName, "Date", "Delivery Date", "send_date", "Send Date", "date") {

		// localize the date
		giftcardDeliveryDate, _ := utils.MakeTZTime(utils.InterfaceToTime(optionMap["value"]), timeZone)

		// format the date if not zero
		if !utils.IsZeroTime(giftcardDeliveryDate) {
			// TODO: Should be "Monday Jan 2 15:04 (MST)" but we have a bug
			return giftcardDeliveryDate.Format("Monday Jan 2 15:04")
		}
	}

	return utils.InterfaceToString(optionMap["value"])
}
