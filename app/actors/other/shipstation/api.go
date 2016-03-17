package shipstation

import (
	"encoding/base64"
	"fmt"
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
	// service.POST("shipstation", updateShipmentStatus)

	return nil
}

func isEnabled(next api.FuncAPIHandler) api.FuncAPIHandler {
	return func(context api.InterfaceApplicationContext) (interface{}, error) {
		isEnabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathShipstationEnabled))

		if !isEnabled {
			// TODO: update status?
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
			// TODO: update status?
			return next(context) //TODO: REMOVE
			// return "not authed", nil
		}

		return next(context)
	}
}

// Your page should return data for any order that was modified between
// the start and end date, regardless of the order’s status.
func listOrders(context api.InterfaceApplicationContext) (interface{}, error) {
	context.SetResponseContentType("text/xml")

	const parseDateFormat = "01/02/2006 15:04"

	// action := context.GetRequestArgument("action") // only expecting "export"
	// page := context.GetRequestArgument("page")
	// Our utils.InterfaceToTime doesn't handle this format well `01/23/2012 17:28`
	startArg := context.GetRequestArgument("start_date")
	endArg := context.GetRequestArgument("end_date")
	startDate, _ := time.Parse(parseDateFormat, startArg)
	endDate, _ := time.Parse(parseDateFormat, endArg)

	fmt.Println(startDate, endDate)

	// Get the orders
	orderCollectionModel, _ := order.GetOrderCollectionModel()
	// orderCollectionModel.ListFilterAdd("updated_at", ">=", startDate)
	// orderCollectionModel.ListFilterAdd("updated_at", "<", endDate)
	orderCollectionModel.ListLimit(0, 5)
	orders := orderCollectionModel.ListOrders()

	// Assemble our response
	response := &Orders{}
	for _, order := range orders {
		responseOrder := buildItem(order)
		response.Orders = append(response.Orders, responseOrder)
	}

	return response, nil
}

func buildItem(order order.InterfaceOrder) Order {
	const outputDateFormat = "01/02/2006 15:04"

	// Base Order Details
	createdAt := utils.InterfaceToTime(order.Get("created_at"))
	updatedAt := utils.InterfaceToTime(order.Get("updated_at"))

	orderDetails := Order{
		OrderId:        order.GetID(),
		OrderNumber:    order.GetID(),
		OrderDate:      createdAt.Format(outputDateFormat),
		OrderStatus:    order.GetStatus(),
		LastModified:   updatedAt.Format(outputDateFormat),
		OrderTotal:     order.GetSubtotal(), //TODO: DOUBLE CHECK THIS IS THE RIGHT ONE
		ShippingAmount: order.GetShippingAmount(),
	}


	// Customer Details
	oShipAddress := order.GetShippingAddress()
	oBillAddress := order.GetBillingAddress()

	customer := Customer{}
	customer.BillingAddress = BillingAddress{
		Name: oBillAddress.GetFirstName() + " " + oBillAddress.GetLastName(),
	}
	customer.ShippingAddress = ShippingAddress{
		Name:  oShipAddress.GetFirstName() + " " + oShipAddress.GetLastName(),
		Address1: oShipAddress.GetAddressLine1(),
		City: oShipAddress.GetCity(),
		State: oShipAddress.GetState(),
		Country: oShipAddress.GetCountry(),
	}

	orderDetails.Customer = customer


	// Order Items
	oItem := range order.GetItems() {
		orderItem := OrderItem{
			Sku: oItem.GetSku(),
			Name: oItem.GetName(),
			Quantity: oItem.GetQty(),
			UnitPrice: oItem.GetPrice(),
		}

		orderDetails.Items = append(orderDetails.Items, orderItem)
	}

	return orderDetails
}
