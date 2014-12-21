package app

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	var err error

	err = api.GetRestService().RegisterAPI("app", "GET", "login", restLogin)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app", "POST", "login", restLogin)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app", "GET", "logout", restLogout)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("app", "GET", "rights", restRightsInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function login application with root rights
func restLogin(params *api.StructAPIHandlerParams) (interface{}, error) {

	var requestLogin string
	var requestPassword string

	if params.Request.Method == "GET" && utils.KeysInMapAndNotBlank(params.RequestGETParams, "login", "password") {
		requestLogin = params.RequestGETParams["login"]
		requestPassword = params.RequestGETParams["password"]

	} else {

		reqData, err := api.GetRequestContentAsMap(params)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		if !utils.KeysInMapAndNotBlank(reqData, "login", "password") {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fee28a56adb144b9a0e21c9be6bd6fdb", "login and password should be specified")
		}

		requestLogin = utils.InterfaceToString(reqData["login"])
		requestPassword = utils.InterfaceToString(reqData["password"])
	}

	rootLogin := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathStoreRootLogin))
	rootPassword := utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathStoreRootPassword))

	if requestLogin == rootLogin && requestPassword == rootPassword {
		params.Session.Set(api.ConstSessionKeyAdminRights, true)

		return "ok", nil
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "68546aa8a6be4c31ac44ea4278dfbdb0", "wrong login or password")
}

// WEB REST API function logout application - session data clear
func restLogout(params *api.StructAPIHandlerParams) (interface{}, error) {
	err := params.Session.Close()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	return "ok", nil
}

// WEB REST API function to get info about current rights
func restRightsInfo(params *api.StructAPIHandlerParams) (interface{}, error) {
	result := make(map[string]interface{})

	result["is_admin"] = utils.InterfaceToBool(params.Session.Get(api.ConstSessionKeyAdminRights))

	return result, nil
}
