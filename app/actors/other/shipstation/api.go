package shipstation

import (
	"fmt"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func setupAPI() error {
	service := api.GetRestService()

	service.GET("shipstation", isEnabled(middleAuth(listOrders)))
	// service.POST("shipstation", updateShipmentStatus)

	return nil
}

func isEnabled(handler api.FuncAPIHandler) api.FuncAPIHandler {
	return func(context api.InterfaceApplicationContext) (interface{}, error) {
		isEnabled := utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathShipstationEnabled))

		if !isEnabled {
			// TODO: update status?
			return "not enabled", nil
		}

		return handler(context)
	}
}

func middleAuth(handler api.FuncAPIHandler) api.FuncAPIHandler {
	return func(context api.InterfaceApplicationContext) (interface{}, error) {
		isAuthed := true //TODO: REAL LOGIC

		if !isAuthed {
			// TODO: update status?
			return "not authed", nil
		}

		return handler(context)
	}
}

// Your page should return data for any order that was modified between
// the start and end date, regardless of the order’s status.
func listOrders(context api.InterfaceApplicationContext) (interface{}, error) {

	fmt.Println(isEnabled)
	// action=export
	// start_date
	// end_date
	// page _if we get this we should log that we need to enable paging_
	//
	return "ok", nil
}
