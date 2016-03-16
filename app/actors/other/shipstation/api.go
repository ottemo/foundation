package shipstation

import (
	"encoding/base64"
	"strings"

	"github.com/ottemo/foundation/api"
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

	// const dateFormat = "01/02/2006 15:04"

	// action := context.GetRequestArgument("action")
	// page := context.GetRequestArgument("page")
	// startArg := context.GetRequestArgument("start_date")
	// endArg := context.GetRequestArgument("end_date")

	// // Our utils.InterfaceToTime doesn't handle this format well `01/23/2012 17:28`
	// startDate, startErr := time.Parse(dateFormat, startArg)
	// endDate, endErr := time.Parse(dateFormat, endArg)

	// if page != "" {
	// 	//TODO: LOG THAT WE ARE SURPRISED
	// }

	// if startErr != nil || endErr != nil {
	// 	//TODO: ERROR WITH INPUTS
	// }

	orders := &Orders{}
	orders.Orders = append(orders.Orders, Order{"Adam"})

	return orders, nil
}
