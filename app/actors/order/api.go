package order

import (
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"

	"github.com/ottemo/foundation/app/actors/discount/giftcard"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	// Admin
	service.GET("orders/attributes", APIListOrderAttributes)
	service.GET("orders", APIListOrders)

	service.GET("order/:orderID", APIGetOrder)
	service.PUT("order/:orderID", APIUpdateOrder)
	service.DELETE("order/:orderID", APIDeleteOrder)
	service.POST("order/:orderID/sendConfirmation", APISendConfirmation)

	// service.POST("order", APICreateOrder)

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

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// taking orders collection model
	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
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

	var orderModel order.InterfaceOrder
	var err error

	if orderModel, err = getOrderID(context); err != nil {
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

// APIUpdateOrder update existing purchase order
//   - order id should be specified in "orderID" argument
func APIUpdateOrder(context api.InterfaceApplicationContext) (interface{}, error) {

	var orderModel order.InterfaceOrder
	var err error
	var orderID string

	if orderModel, err = getOrderID(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}
	orderID = orderModel.GetID()

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		orderModel.Set(attribute, value)
	}

	// TODO: shouldn't need this
	orderModel.SetID(orderID)
	orderModel.Save()

	return orderModel.ToHashMap(), nil
}

// APIDeleteOrder deletes existing purchase order
//   - order id should be specified in "orderID" argument
func APIDeleteOrder(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	orderID := context.GetRequestArgument("orderID")
	if orderID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fc3011c7-e58c-4433-b9b0-881a7ba005cf", "order id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	orderModel, err := order.GetOrderModelAndSetID(orderID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderModel.Delete()

	return "ok", nil
}

// APIGetOrder returns current visitor order details for specified order
//   - orderID should be specified in arguments
func APIGetVisitOrder(context api.InterfaceApplicationContext) (interface{}, error) {

	orderModel, err := order.LoadOrderByID(context.GetRequestArgument("orderID"))
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cc719b01-c1e4-4b69-9c89-5735f5c0d339", "Unable to retrieve an order associated with OrderID: "+context.GetRequestArgument("orderID")+".")
	}

	// allow anonymous visitors through if the session id matches
	if utils.InterfaceToString(orderModel.Get("session_id")) != context.GetSession().GetID() {
		// force anonymous visitors to log in if their session id does not match the one on the order
		visitorID := visitor.GetCurrentVisitorID(context)
		if visitorID == "" {
			return "No Visitor ID found, unable to process order request. Please log in first.", nil
		} else if utils.InterfaceToString(orderModel.Get("visitor_id")) != visitorID {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c5ca1fdb-7008-4a1c-a168-9df544df9825", "There is a mis-match between the current Visitor ID and the Visitor ID on the order.")
		}
	}

	result := orderModel.ToHashMap()
	result["items"] = orderModel.GetItems()

	return result, nil
}

// APIGetOrders returns list of orders related to current visitor
func APIGetVisitOrders(context api.InterfaceApplicationContext) (interface{}, error) {

	// list operation
	//---------------
	visitorID := visitor.GetCurrentVisitorID(context)
	if visitorID == "" {
		return "No Visitor ID found, unable to process request.  Please log in first.", nil
	}

	orderCollection, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = orderCollection.ListFilterAdd("visitor_id", "=", visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// We only return orders that are in these two states
	statusFilter := [2]string{order.ConstOrderStatusProcessed, order.ConstOrderStatusCompleted}
	orderCollection.GetDBCollection().AddFilter("status", "in", statusFilter)

	// filters handle
	models.ApplyFilters(context, orderCollection.GetDBCollection())

	// extra parameter handle
	models.ApplyExtraAttributes(context, orderCollection)

	result, err := orderCollection.List()

	return result, env.ErrorDispatch(err)
}

// APISendConfirmation will send out an order confirmation email to the visitor specficied in the orderID
func APISendConfirmation(context api.InterfaceApplicationContext) (interface{}, error) {

	var orderModel order.InterfaceOrder
	var orderItems []map[string]interface{}
	var err error
	var email, emailAddress, timeZone, giftCardSku string

	if orderModel, err = getOrderID(context); err != nil {
		return "failure", env.ErrorDispatch(err)
	}

	email = utils.InterfaceToString(env.ConfigGetValue(checkout.ConstConfigPathConfirmationEmail))
	timeZone = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))
	giftCardSku = utils.InterfaceToString(env.ConfigGetValue(giftcard.ConstConfigPathGiftCardSKU))

	// create visitor map
	visitor := make(map[string]interface{})
	visitor["first_name"] = orderModel.Get("customer_name")
	visitor["email"] = orderModel.Get("customer_email")

	// create order map
	order := orderModel.ToHashMap()

	// load order items
	for _, item := range orderModel.GetItems() {
		options := make(map[string]interface{})

		for optionName, optionKeys := range item.GetOptions() {
			optionMap := utils.InterfaceToMap(optionKeys)
			options[optionName] = optionMap["value"]

			// Giftcard's delivery date
			if strings.Contains(item.GetSku(), giftCardSku) {
				if utils.IsAmongStr(optionName, "Date", "Delivery Date", "send_date", "Send Date", "date") {
					// Localize and format the date
					giftcardDeliveryDate, _ := utils.MakeTZTime(utils.InterfaceToTime(optionMap["value"]), timeZone)
					if !utils.IsZeroTime(giftcardDeliveryDate) {
						//TODO: Should be "Monday Jan 2 15:04 (MST)" but we have a bug
						options[optionName] = giftcardDeliveryDate.Format("Monday Jan 2 15:04")
					}
				}
			}
		}

		orderItems = append(orderItems, map[string]interface{}{
			"name":    item.GetName(),
			"options": options,
			"sku":     item.GetSku(),
			"qty":     item.GetQty(),
			"price":   item.GetPrice()})
	}

	// convert date of order creation to store time zone
	if date, present := order["created_at"]; present {
		convertedDate, _ := utils.MakeTZTime(utils.InterfaceToTime(date), timeZone)
		if !utils.IsZeroTime(convertedDate) {
			order["created_at"] = convertedDate
		}
	}

	order["items"] = orderItems
	order["payment_method_title"] = orderModel.GetPaymentMethod()
	order["shipping_method_title"] = orderModel.GetShippingMethod()

	customInfo := make(map[string]interface{})
	customInfo["base_storefront_url"] = utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStorefrontURL))

	confirmationEmail, err := utils.TextTemplate(email, map[string]interface{}{
		"Order":   order,
		"Visitor": visitor,
		"Info":    customInfo,
	})
	if err != nil {
		return "failure", env.ErrorDispatch(err)
	}

	emailAddress = utils.InterfaceToString(visitor["email"])
	err = app.SendMail(emailAddress, "Order confirmation", confirmationEmail)
	if err != nil {
		return "failure", env.ErrorDispatch(err)
	}

	return "success", nil
}

// getOrderID will load the order from the database
func getOrderID(context api.InterfaceApplicationContext) (order.InterfaceOrder, error) {

	var orderID string
	var orderModel order.InterfaceOrder
	var err error

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// load orderID
	if orderID = context.GetRequestArgument("orderID"); orderID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fc3011c7-e58c-4433-b9b0-881a7ba005cf", "order id should be specified")
	}

	if orderModel, err = order.LoadOrderByID(orderID); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return orderModel, nil
}
