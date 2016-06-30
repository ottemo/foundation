package order

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// ------------------
// Internal functions
// ------------------

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	// Admin
	service.GET("orders/attributes", api.IsAdmin(APIListOrderAttributes))
	service.GET("orders", api.IsAdmin(APIListOrders))

	service.GET("order/:orderID", api.IsAdmin(APIGetOrder))
	service.PUT("order/:orderID", api.IsAdmin(APIUpdateOrder))
	service.DELETE("order/:orderID", api.IsAdmin(APIDeleteOrder))
	service.GET("order/:orderID/emailShipStatus", api.IsAdmin(APISendShipStatusEmail))
	service.GET("order/:orderID/emailOrderConfirmation", api.IsAdmin(APISendOrderConfirmationEmail))

	// Public
	service.GET("visit/orders", APIGetVisitorOrders)
	service.GET("visit/order/:orderID", APIGetVisitorOrder)

	return nil
}

// apiFindSpecifiedOrder tries for find specified order ID among request argumants
func apiFindSpecifiedOrder(context api.InterfaceApplicationContext) (order.InterfaceOrder, error) {

	// looking for specified order ID
	orderID := ""
	for _, key := range []string{"orderID", "order", "order_id"} {
		if value := context.GetRequestArgument(key); value != "" {
			orderID = value
		}
	}

	// returning error if order ID was not specified
	if orderID == "" {
		context.SetResponseStatusBadRequest()
		return orderID, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fc3011c7-e58c-4433-b9b0-881a7ba005cf", "No order id found on request, orderID should be specified")
	}

	orderModel, err := order.LoadOrderByID(orderID)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	return orderModel, nil
}

// -------------
// API functions
// -------------

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
	orderModel, err := apiFindSpecifiedOrder(context)
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

// APISendShipStatusEmail will send the visitor a shipping confirmation email
// - order id should be specified in "orderID" argument
func APISendShipStatusEmail(context api.InterfaceApplicationContext) (interface{}, error) {

	orderModel, err := apiFindSpecifiedOrder(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	err = orderModel.SendShippingStatusUpdateEmail()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return "Shipping status email sent", nil
}

// APIUpdateOrder update existing purchase order
//   - order id should be specified in "orderID" argument
func APIUpdateOrder(context api.InterfaceApplicationContext) (interface{}, error) {

	orderModel, err := apiFindSpecifiedOrder(context)
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
	orderModel, err := apiFindSpecifiedOrder(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderModel.Delete()
	return "Order deleted: " + orderModel.GetID(), nil
}

// APIGetVisitorOrder returns current visitor order details for specified order
//   - orderID should be specified in arguments
func APIGetVisitorOrder(context api.InterfaceApplicationContext) (interface{}, error) {

	orderModel, err := apiFindSpecifiedOrder(context)
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

// APIGetVisitorOrders returns list of orders related to current visitor
//   - visitorID is required, visitor must be logged in
func APIGetVisitorOrders(context api.InterfaceApplicationContext) (interface{}, error) {

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

// APISendOrderConfirmationEmail will send out an order confirmation email to the visitor specficied in the orderID
//   - orderID must be passed as a request argument
func APISendOrderConfirmationEmail(context api.InterfaceApplicationContext) (interface{}, error) {

	// loading the order model
	orderModel, err := apiFindSpecifiedOrder(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// preparing template object "Info"
	customInfo := make(map[string]interface{})
	customInfo["base_storefront_url"] = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStorefrontURL))

	// preparing template object "Visitor"
	visitor := make(map[string]interface{})
	visitor["first_name"] = orderModel.Get("customer_name")
	visitor["email"] = orderModel.Get("customer_email")

	// preparing template object "Order"
	order := orderModel.ToHashMap()
	order["payment_method_title"] = orderModel.GetPaymentMethod()
	order["shipping_method_title"] = orderModel.GetShippingMethod()

	// the dates in order should be converted to clients locale
	// TODO: the dates to locale conversion should not happens there - it should be either part of order helper or utilities routine over resulting map
	timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))

	// "created_at" date conversion
	if date, present := order["created_at"]; present {
		convertedDate, _ := utils.MakeTZTime(utils.InterfaceToTime(date), timeZone)
		if !utils.IsZeroTime(convertedDate) {
			order["created_at"] = convertedDate
		}
	}

	// order items extraction
	var items []map[string]interface{}
	for _, item := range orderModel.GetItems() {
		// the item options could also contain the date, which should be converted
		for key, value := range item.GetOptions() {
			if utils.IsAmongStr(key, "Date", "Delivery Date", "send_date", "Send Date", "date") {
				localizedDate, _ := utils.MakeTZTime(utils.InterfaceToTime(value), timeZone)
				if !utils.IsZeroTime(localizedDate) {
					value[key] = localizedDate
				}
			}
		}
		items = append(items, item.ToHashMap())
	}
	order["items"] = items

	// processing email template
	template := utils.InterfaceToString(env.ConfigGetValue(checkout.ConstConfigPathConfirmationEmail))
	confirmationEmail, err := utils.TextTemplate(template, map[string]interface{}{
		"Order":   order,
		"Visitor": visitor,
		"Info":    customInfo,
	})
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return "failure", env.ErrorDispatch(err)
	}

	// sending the email notification
	emailAddress := utils.InterfaceToString(visitor["email"])
	err = app.SendMail(emailAddress, "Order confirmation", confirmationEmail)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return "failure", env.ErrorDispatch(err)
	}

	return "Order confirmation email sent", nil
}
