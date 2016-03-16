package shipstation

import (
	"encoding/base64"
	"fmt"
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
			hashParts := strings.SplitN(authHash, " ", 2)
			if len(hashParts) != 2 {
				return false
			}

			b, err := base64.StdEncoding.DecodeString(hashParts[1])
			if err != nil {
				return false
			}

			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				return false
			}

			return pair[0] == username && pair[1] == password
		}

		if !isAuthed(authHash, username, password) {
			// TODO: update status?
			return "not authed", nil
		}

		return next(context)
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
