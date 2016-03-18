package shipstation

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupAPI() error {
	service := api.GetRestService()

	service.GET("shipstation", isEnabled(basicAuth(listOrders)))
	service.POST("shipstation", isEnabled(basicAuth(updateShipmentStatus)))

	return nil
}

func isEnabled(next api.FuncAPIHandler) api.FuncAPIHandler {
	return func(context api.InterfaceApplicationContext) (interface{}, error) {
		isEnabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathShipstationEnabled))

		if !isEnabled {
			// context.SetResponseStatusNotFound()
			// return "not enabled", nil
			return next(context) //TODO: REMOVE
		}

		return next(context)
	}
}

func basicAuth(next api.FuncAPIHandler) api.FuncAPIHandler {
	return func(context api.InterfaceApplicationContext) (interface{}, error) {

		authHash := utils.InterfaceToString(context.GetRequestSetting("Authorization"))
		username := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathShipstationUsername))
		password := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathShipstationPassword))

		isAuthed := func(authHash string, username string, password string) bool {
			// authHash := "Basic jalsdfjaklsdfjalksdjf"
			hashParts := strings.SplitN(authHash, " ", 2)
			if len(hashParts) != 2 {
				return false
			}

			decodedHash, err := base64.StdEncoding.DecodeString(hashParts[1])
			if err != nil {
				return false
			}

			userPass := strings.SplitN(string(decodedHash), ":", 2)
			if len(userPass) != 2 {
				return false
			}

			return userPass[0] == username && userPass[1] == password
		}

		if !isAuthed(authHash, username, password) {
			// context.SetResponseStatusForbidden()
			// return "not authed", nil
			return next(context) //TODO: REMOVE
		}

		return next(context)
	}
}

// Handler for getting a list of orders
// - XML formatted response
// - Should return any orders that were modified within the date range
//   regardless of the order status
func listOrders(context api.InterfaceApplicationContext) (interface{}, error) {
	context.SetResponseContentType("text/xml")

	// Our utils.InterfaceToTime doesn't handle this format well `01/23/2012 17:28`
	const parseDateFormat = "01/02/2006 15:04"
	const exportAction = "export"

	// The only action this endpoint accepts is "export"
	action := context.GetRequestArgument("action")
	if action != exportAction {
		return nil, nil
	}

	startArg := context.GetRequestArgument("start_date")
	endArg := context.GetRequestArgument("end_date")
	startDate, _ := time.Parse(parseDateFormat, startArg)
	endDate, _ := time.Parse(parseDateFormat, endArg)
	// page := context.GetRequestArgument("page") // we don't paginate currently

	// Get the orders
	orderQuery := getOrders(startDate, endDate)

	// Get the order items
	var orderIds []string
	for _, orderResult := range orderQuery {
		orderIds = append(orderIds, orderResult.GetID())
	}
	oiResults := getOrderItems(orderIds)

	// Assemble our response
	response := &Orders{}
	for _, orderResult := range orderQuery {
		responseOrder := buildItem(orderResult, oiResults)
		response.Orders = append(response.Orders, responseOrder)
	}

	return response, nil
}

// db query for getting all orders
func getOrders(startDate time.Time, endDate time.Time) []order.InterfaceOrder {
	oModel, _ := order.GetOrderCollectionModel()
	oModel.GetDBCollection().AddFilter("updated_at", ">=", startDate)
	oModel.GetDBCollection().AddFilter("updated_at", "<", endDate)
	result := oModel.ListOrders()

	return result
}

// db query for getting all relavent order items
func getOrderItems(orderIds []string) []map[string]interface{} {
	oiModel, _ := order.GetOrderItemCollectionModel()
	oiDB := oiModel.GetDBCollection()
	oiDB.AddFilter("order_id", "in", orderIds)
	oiResults, _ := oiDB.Load()
	// NOTE: If we could FromHashMap this into a struct i'd be happier
	// as is this is the only place where i'm forced to pass around an
	// ugly variable

	return oiResults
}

// Convert an ottemo order and all possible orderitems into a shipstation order
func buildItem(oItem order.InterfaceOrder, allOrderItems []map[string]interface{}) Order {
	const outputDateFormat = "01/02/2006 15:04"

	// Base Order Details
	createdAt := utils.InterfaceToTime(oItem.Get("created_at"))
	updatedAt := utils.InterfaceToTime(oItem.Get("updated_at"))

	orderDetails := Order{
		OrderId:        oItem.GetID(),
		OrderNumber:    oItem.GetID(),
		OrderDate:      createdAt.Format(outputDateFormat),
		OrderStatus:    oItem.GetStatus(),
		LastModified:   updatedAt.Format(outputDateFormat),
		OrderTotal:     oItem.GetSubtotal(),       // TODO: DOUBLE CHECK THIS IS THE RIGHT ONE, AND FORMAT?
		ShippingAmount: oItem.GetShippingAmount(), // TODO: FORMAT?
	}

	// Customer Details
	orderDetails.Customer.CustomerCode = utils.InterfaceToString(oItem.Get("customer_email"))

	oBillAddress := oItem.GetBillingAddress()
	orderDetails.Customer.BillingAddress = BillingAddress{
		Name: oBillAddress.GetFirstName() + " " + oBillAddress.GetLastName(),
	}

	oShipAddress := oItem.GetShippingAddress()
	orderDetails.Customer.ShippingAddress = ShippingAddress{
		Name:     oShipAddress.GetFirstName() + " " + oShipAddress.GetLastName(),
		Address1: oShipAddress.GetAddressLine1(),
		City:     oShipAddress.GetCity(),
		State:    oShipAddress.GetState(),
		Country:  oShipAddress.GetCountry(),
	}

	// Order Items
	for _, oiItem := range allOrderItems {
		isThisOrder := oiItem["order_id"] == oItem.GetID()
		if !isThisOrder {
			continue
		}

		orderItem := OrderItem{
			Sku:       utils.InterfaceToString(oiItem["sku"]),
			Name:      utils.InterfaceToString(oiItem["name"]),
			Quantity:  utils.InterfaceToInt(oiItem["qty"]),
			UnitPrice: utils.InterfaceToFloat64(oiItem["price"]), // TODO: FORMAT?
		}

		orderDetails.Items = append(orderDetails.Items, orderItem)
	}

	return orderDetails
}

// action	The value will always be "shipnotify" when sending shipping notifications.
// order_number	This is the order's unique identifier.
// carrier	USPS, UPS, FedEx, DHL, Other, DHLGlobalMail, UPSMI, BrokersWorldWide, FedExInternationalMailService, CanadaPost, FedExCanada, OnTrac, Newgistics, FirstMile, Globegistics, LoneStar, Asendia, RoyalMail, APC, AccessWorldwide, AustraliaPost, DHLCanada, IMEX
// service	This will be the name of the shipping service that was used to ship the order.
// tracking_number	This is the tracking number for the package.
func updateShipmentStatus(context api.InterfaceApplicationContext) (interface{}, error) {
	const expectedAction = "shipnotify"

	action := context.GetRequestArgument("action")
	if action != expectedAction {
		context.SetResponseStatusBadRequest()
		return nil, nil
	}

	orderID := context.GetRequestArgument("order_number")
	carrier := context.GetRequestArgument("carrier")
	service := context.GetRequestArgument("service")
	trackingNumber := context.GetRequestArgument("tracking_number")

	// fmt.Println(action, orderID, carrier, service, trackingNumber)

	orderModel, orderNotFound := order.LoadOrderByID(orderID)
	if orderNotFound != nil {
		context.SetResponseStatusBadRequest()
		return nil, nil
	}

	shippingInfo := utils.InterfaceToMap(orderModel.Get("shipping_info"))
	shippingInfo["carrier"] = carrier
	shippingInfo["service"] = service
	shippingInfo["tracking_number"] = trackingNumber

	orderModel.Set("shipping_info", shippingInfo)
	orderModel.Set("updated_at", time.Now())
	orderModel.Save()

	return nil, nil
}
