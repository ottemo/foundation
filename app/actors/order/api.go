package order

import (
	"encoding/csv"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
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
	service.POST("orders/exportToCSV", api.IsAdmin(APIExportOrders))
	service.POST("orders/setStatus", api.IsAdmin(APIChangeOrderStatus))

	service.GET("order/:orderID", api.IsAdmin(APIGetOrder))
	service.PUT("order/:orderID", api.IsAdmin(APIUpdateOrder))
	service.DELETE("order/:orderID", api.IsAdmin(APIDeleteOrder))
	service.GET("order/:orderID/emailShipStatus", api.IsAdmin(APISendShipStatusEmail))
	service.GET("order/:orderID/emailOrderConfirmation", api.IsAdmin(APISendOrderConfirmationEmail))
	service.POST("order/:orderID/emailTrackingCode", api.IsAdmin(APIUpdateTrackingInfoAndSendEmail))

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
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fc3011c7-e58c-4433-b9b0-881a7ba005cf", "No order id found on request, orderID should be specified")
	}

	orderModel, err := order.LoadOrderByID(orderID)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4fd03907-4e6d-46d5-981e-8f858f1aa83f9", "system error loading id from db: "+utils.InterfaceToString(orderID))
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
	if err := models.ApplyFilters(context, orderCollectionModel.GetDBCollection()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "feb1e3b3-fb65-4e77-9524-d64f4cf574e8", err.Error())
	}

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return orderCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	if err := orderCollectionModel.ListLimit(models.GetListLimit(context)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3257c74b-a809-45ed-8863-14ee273d3f3b", err.Error())
	}

	// extra parameter handle
	if err := models.ApplyExtraAttributes(context, orderCollectionModel); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5e81c9be-1df8-436b-8c22-b9b29949dd1b", err.Error())
	}

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
		if err := orderModel.Set(attribute, value); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b563c756-92a7-4d08-b85b-d365c713cd69", err.Error())
		}
	}

	if err := orderModel.Save(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9a5c8363-baad-4060-a287-4d88b46878a6", err.Error())
	}

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

	if err := orderModel.Delete(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ad203565-3975-4f7c-acea-28d4a58be966", err.Error())
	}
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
	if err := orderCollection.GetDBCollection().AddFilter("status", "in", statusFilter); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "88c33afe-4341-470f-a9ab-da024927e5ea", err.Error())
	}

	descending := true
	if err := orderCollection.GetDBCollection().AddSort("created_at", descending); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c923b47c-fa6c-42ad-8c27-7c86688866bb", err.Error())
	}

	// filters handle
	if err := models.ApplyFilters(context, orderCollection.GetDBCollection()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2190f86e-1aef-4461-8107-37533c01d1ab", err.Error())
	}

	// extra parameter handle
	if err := models.ApplyExtraAttributes(context, orderCollection); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "16e512ae-3cfb-4d92-a671-0e6dea9446e3", err.Error())
	}

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

	err = orderModel.SendOrderConfirmationEmail()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return "Order confirmation email sent", nil
}

// APIExportOrders returns a list of orders in CSV format
//    - returns orders specified in url parameters
//    - must include at least one order
func APIExportOrders(context api.InterfaceApplicationContext) (interface{}, error) {
	var itemCSVRecords [][]string
	var orders []string

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// must have at least one order in request
	if !utils.KeysInMapAndNotBlank(requestData, "orders") {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f2602f73-7cae-4525-8405-9e470681c20e", "Specifiy a minimum of one order in the orders parameter.")
	}

	// look up the order items collection
	orderItemsCollectionModel, err := order.GetOrderItemCollectionModel()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}
	dbOrderItemsCollection := orderItemsCollectionModel.GetDBCollection()

	// look up the order collection
	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}
	dbOrderCollection := orderCollectionModel.GetDBCollection()
	if err := dbOrderCollection.AddSort("created_at", false); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d237815b-a87b-4c2a-a4e0-fa28e1ed6459", err.Error())
	}

	// load orders based on order IDs passed
	orders = utils.InterfaceToStringArray(utils.InterfaceToArray(requestData["orders"]))
	if orders != nil && len(orders) > 0 && !utils.IsInListStr("all", orders) {
		if err := dbOrderCollection.AddFilter("_id", "in", orders); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bb689aa7-3cab-4f8b-b504-165aa8b7f00a", err.Error())
		}
	}

	ordersRecords, err := dbOrderCollection.Load()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// error if no records returned from db
	if len(ordersRecords) == 0 {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "28eac91b-39ec-4034-b664-4004e940a6d1", "No orders were found.")
	}

	// build the field names for csv file
	for _, columnsHeaders := range orderFields {
		itemCSVRecords = append(itemCSVRecords, columnsHeaders)
	}

	// Webgility importer bombs out on some characters even if they are properly
	// escaped in a csv: . , \ / ( )
	regexCleaner, _ := regexp.Compile(`[.,\/\\()]`)

	for _, orderRecord := range ordersRecords {

		if err := dbOrderItemsCollection.ClearFilters(); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b88f2ac6-308c-486c-8e67-a014cf8e12ef", err.Error())
		}
		if err := dbOrderItemsCollection.AddFilter("order_id", "=", orderRecord["_id"]); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "257eef6d-b273-4d4e-95aa-c4fcdd6aed27", err.Error())
		}

		orderItemsRecords, err := dbOrderItemsCollection.Load()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		orderRecord := utils.InterfaceToMap(orderRecord)

		// convert discount to an absolute value to fix quickbooks bug
		discount := utils.InterfaceToFloat64(orderRecord["discount"])
		discount = math.Abs(discount)
		orderRecord["discount"] = discount
		shippingAddress := utils.InterfaceToMap(orderRecord["shipping_address"])
		billingAddress := utils.InterfaceToMap(orderRecord["billing_address"])

		for orderItemIndex, orderItem := range orderItemsRecords {

			// process an order and convert it to csv format
			for _, inputValues := range dataSet {
				var rowData []string

				// build the csv row from the order data set
				for _, value := range inputValues {
					cellValue := ""

					switch typedValue := value.(type) {
					case string:
						cellValue = typedValue
						switch {
						case strings.Index(cellValue, "$") == 0:
							cellValue = utils.InterfaceToString(orderRecord[strings.Replace(cellValue, "$", "", 1)])
							break

						case strings.HasPrefix(cellValue, "item."):
							cellValue = utils.InterfaceToString(orderItem[strings.Replace(cellValue, "item.", "", 1)])
							break

						case strings.HasPrefix(cellValue, "shipping."):
							addressKey := strings.Replace(cellValue, "shipping.", "", 1)
							cellValue = utils.InterfaceToString(shippingAddress[addressKey])
							cellValue = regexCleaner.ReplaceAllString(cellValue, "")
							break

						case strings.HasPrefix(cellValue, "billing."):
							addressKey := strings.Replace(cellValue, "billing.", "", 1)
							cellValue = utils.InterfaceToString(billingAddress[addressKey])
							cellValue = regexCleaner.ReplaceAllString(cellValue, "")
							break
						}
						break

					case func(record map[string]interface{}) string:

						cellValue = typedValue(orderRecord)
						break

					case func(int, map[string]interface{}) string:
						cellValue = typedValue(orderItemIndex, orderItem)
						break

					}
					// append cellValue to row
					rowData = append(rowData, cellValue)
				}
				// add row to table
				itemCSVRecords = append(itemCSVRecords, rowData)
			}
		}
	}

	// preparing csv writer
	csvWriter := csv.NewWriter(context.GetResponseWriter())

	exportFilename := "orders_export_" + time.Now().Format(time.RFC3339) + ".csv"

	if err := context.SetResponseContentType("text/csv"); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4fd36680-5527-410a-be51-c02f79b986aa", err.Error())
	}
	if err := context.SetResponseSetting("Content-disposition", "attachment;filename="+exportFilename); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0c937338-a89d-4980-9489-ef5f7508a81b", err.Error())
	}

	for _, csvRecord := range itemCSVRecords {
		if err := csvWriter.Write(csvRecord); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b0435377-138d-448e-b8b5-ed37d1f44549", err.Error())
		}
	}
	csvWriter.Flush()

	return "", nil
}

// APIChangeOrderStatus will change orders to the state included in the status request variable
//   - order ids should be specified in "IDs" argument
//   - status should be specified in "status" argument
func APIChangeOrderStatus(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	statusValue, present := requestData["status"]
	if !present {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3d00647d-505a-4092-b821-20dd8638e471", "missing argument in request: status")
	}
	status := utils.InterfaceToString(statusValue)

	orderIDsValue, present := requestData["order_id"]
	if !present {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4456c336-96f0-4b9b-a54a-ab0409645f64", "missing argument in request: order_id")
	}
	orderIDs := utils.InterfaceToArray(orderIDsValue)

	if err = updateOrderStatus(orderIDs, status); err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// change the order status and persist new status to the db
//    - status is the new order status to be saved
func updateOrderStatus(orderIDs []interface{}, status string) error {

	for _, orderID := range orderIDs {
		orderModel, err := order.LoadOrderByID(utils.InterfaceToString(orderID))
		if err != nil {
			return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8cb7a9cd-10fd-4a3b-9e5d-336075cd16e9", "error loading id from db: "+utils.InterfaceToString(orderID))
		}
		if err = orderModel.SetStatus(status); err != nil {
			return env.ErrorDispatch(err)
		}
		if err = orderModel.Save(); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// APIUpdateTrackingInfoAndSendEmail updates order with shipping tracking information and sends a shipping status email
// - carrier, tracking_number, tracking_url are required
func APIUpdateTrackingInfoAndSendEmail(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	orderModel, err := apiFindSpecifiedOrder(context)
	if err != nil {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorDispatch(err)
	}

	carrier := utils.InterfaceToString(requestData["carrier"])
	if carrier == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "29948535-f32a-4393-b5e8-c66092bbbe6d", "carrier should be specified")
	}
	trackingNumber := utils.InterfaceToString(requestData["tracking_number"])
	if trackingNumber == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "73db2a29-c22d-40a1-a845-fc4ac03b100a", "tracking number should be specified")
	}
	trackingURL := utils.InterfaceToString(requestData["tracking_url"])
	if trackingURL == "" {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "04865d8a-58b5-4296-8caa-c9dff49136a3", "tracking url should be specified")
	}

	shippingInfo := utils.InterfaceToMap(orderModel.Get("shipping_info"))
	shippingInfo["carrier"] = carrier
	shippingInfo["tracking_number"] = trackingNumber
	shippingInfo["tracking_url"] = trackingURL
	if err := orderModel.Set("shipping_info", shippingInfo); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c786a6dc-a0b9-486e-a7d5-0b922cd31902", err.Error())
	}

	err = orderModel.Save()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	if err := orderModel.SendShippingStatusUpdateEmail(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c258f6eb-97fc-4eb1-a91a-615d8abbe89a", err.Error())
	}

	return "ok", nil
}
