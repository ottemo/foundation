package shipstation

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
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
			context.SetResponseStatusNotFound()
			return "not enabled", nil
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
			context.SetResponseStatusForbidden()
			return "not authed", nil
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
	oResults := order.GetFullOrdersUpdatedBetween(startDate, endDate)

	// Get the order items
	var orderIds []string
	for _, orderResult := range oResults {
		orderIds = append(orderIds, orderResult.GetID())
	}

	oiResults := order.GetItemsForOrders(orderIds)

	// Assemble our response
	response := &Orders{}
	for _, orderResult := range oResults {
		responseOrder := buildItem(orderResult, oiResults)
		response.Orders = append(response.Orders, responseOrder)
	}

	return response, nil
}

// Convert an ottemo order and all possible orderitems into a shipstation order
func buildItem(oItem order.InterfaceOrder, allOrderItems []map[string]interface{}) Order {
	const outputDateFormat = "01/02/2006 15:04"

	// Base Order Details
	createdAt := utils.InterfaceToTime(oItem.Get("created_at"))
	updatedAt := utils.InterfaceToTime(oItem.Get("updated_at"))
	var customInfo = utils.InterfaceToMap(oItem.Get("custom_info"))
	var calculation = utils.InterfaceToMap(customInfo["calculation"])

	orderDetails := Order{
		OrderId:        oItem.GetID(),
		OrderNumber:    oItem.GetID(),
		OrderDate:      createdAt.Format(outputDateFormat),
		OrderStatus:    oItem.GetStatus(),
		LastModified:   updatedAt.Format(outputDateFormat),
		TaxAmount:      oItem.GetTaxAmount(),
		ShippingAmount: oItem.GetShippingAmount(),
		OrderTotal:     oItem.GetGrandTotal(),
	}

	// Customer Details
	orderDetails.Customer.CustomerCode = utils.InterfaceToString(oItem.Get("customer_email"))

	oBillAddress := oItem.GetBillingAddress()
	orderDetails.Customer.BillingAddress = BillingAddress{
		Name: oBillAddress.GetFirstName() + " " + oBillAddress.GetLastName(),
	}

	oShipAddress := oItem.GetShippingAddress()
	orderDetails.Customer.ShippingAddress = ShippingAddress{
		Name:       oShipAddress.GetFirstName() + " " + oShipAddress.GetLastName(),
		Address1:   oShipAddress.GetAddressLine1(),
		City:       oShipAddress.GetCity(),
		State:      oShipAddress.GetState(),
		PostalCode: oShipAddress.GetZipCode(),
		Country:    oShipAddress.GetCountry(),
	}

	var calculatedDiscounts float64
	var calculatedSubtotal float64

	// Order Items
	for _, oiItem := range allOrderItems {
		isThisOrder := oiItem["order_id"] == oItem.GetID()
		if !isThisOrder {
			continue
		}

		var oiItemPrice = utils.InterfaceToFloat64(oiItem["price"])
		var oiItemQty = utils.InterfaceToInt(oiItem["qty"])

		orderItem := OrderItem{
			Sku:       utils.InterfaceToString(oiItem["sku"]),
			Name:      utils.InterfaceToString(oiItem["name"]),
			Quantity:  oiItemQty,
			UnitPrice: oiItemPrice, // TODO: FORMAT?
		}
		orderDetails.Items = append(orderDetails.Items, orderItem)
		calculatedSubtotal += utils.InterfaceToFloat64(oiItemQty) * oiItemPrice

		if calculation != nil {
			var oiItemIdx = utils.InterfaceToString(oiItem["idx"])
			if oiItemCalculation := calculation[oiItemIdx]; oiItemCalculation != nil {
				var oiItemCalculationMap = utils.InterfaceToMap(oiItemCalculation)
				var oiItemDiscountedPrice = utils.InterfaceToFloat64(oiItemCalculationMap[checkout.ConstLabelGrandTotal]) / utils.InterfaceToFloat64(orderItem.Quantity)

				if utils.RoundPrice(oiItemPrice-oiItemDiscountedPrice) != 0 {
					orderItem := OrderItem{
						Sku:        "",
						Name:       "Discount on " + utils.InterfaceToString(oiItem["name"]),
						Quantity:   oiItemQty,
						UnitPrice:  utils.RoundPrice(oiItemDiscountedPrice - utils.InterfaceToFloat64(oiItem["price"])),
						Adjustment: true,
					}
					orderDetails.Items = append(orderDetails.Items, orderItem)
					calculatedDiscounts += (oiItemDiscountedPrice - utils.InterfaceToFloat64(oiItem["price"])) * utils.InterfaceToFloat64(oiItemQty)
				}
			}
		}
	}

	// apply whole order discount
	if calculation != nil {
		var calculatedGrandTotal = calculatedSubtotal + calculatedDiscounts + orderDetails.ShippingAmount + orderDetails.TaxAmount
		var orderDiscount = oItem.GetGrandTotal() - calculatedGrandTotal

		if utils.RoundPrice(orderDiscount) != 0 {
			orderItem := OrderItem{
				Sku:        "",
				Name:       "Discount on Order",
				Quantity:   1,
				UnitPrice:  utils.RoundPrice(orderDiscount),
				Adjustment: true,
			}
			orderDetails.Items = append(orderDetails.Items, orderItem)
		}
	}

	return orderDetails
}

// updateShipmentStatus An endpoint for shipstation to hit that will update the order with some shipment tracking info
// and then send off an email update
//
// - action :			The value will always be "shipnotify" when sending shipping notifications.
// - order_number :		This is the order's unique identifier.
// - carrier :			USPS, UPS, FedEx, DHL, Other, DHLGlobalMail, UPSMI, BrokersWorldWide, FedExInternationalMailService,
// 						CanadaPost, FedExCanada, OnTrac, Newgistics, FirstMile, Globegistics, LoneStar, Asendia,
// 						RoyalMail, APC, AccessWorldwide, AustraliaPost, DHLCanada, IMEX
// - service :			This will be the name of the shipping service that was used to ship the order.
// - tracking_number :	This is the tracking number for the package.
func updateShipmentStatus(context api.InterfaceApplicationContext) (interface{}, error) {
	const expectedAction = "shipnotify"

	action := context.GetRequestArgument("action")
	if action != expectedAction {
		context.SetResponseStatusBadRequest()
		return "unexpected action", nil
	}

	orderID := context.GetRequestArgument("order_number")
	carrier := context.GetRequestArgument("carrier")
	service := context.GetRequestArgument("service")
	trackingNumber := context.GetRequestArgument("tracking_number")

	orderModel, orderNotFound := order.LoadOrderByID(orderID)
	if orderNotFound != nil {
		context.SetResponseStatusBadRequest()
		return nil, nil
	}

	shippingInfo := utils.InterfaceToMap(orderModel.Get("shipping_info"))
	shippingInfo["carrier"] = carrier
	shippingInfo["service"] = service
	shippingInfo["tracking_number"] = trackingNumber
	shippingInfo["tracking_url"] = buildTrackingUrl(carrier, trackingNumber)

	orderModel.Set("shipping_info", shippingInfo)
	orderModel.Set("updated_at", time.Now())
	orderModel.SetStatus(order.ConstOrderStatusCompleted)
	err := orderModel.Save()

	if err != nil {
		context.SetResponseStatusBadRequest()
	} else {
		orderModel.SendShippingStatusUpdateEmail()
	}

	return nil, nil
}
