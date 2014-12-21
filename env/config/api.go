package config

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/utils"
)

// setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("config", "GET", "groups", restConfigGroups)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config", "GET", "info/:path", restConfigInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config", "GET", "list", restConfigList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config", "GET", "get/:path", restConfigGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config", "POST", "set/:path", restConfigSet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config", "POST", "register", restConfigRegister)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config", "DELETE", "unregister/:path", restConfigUnRegister)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config", "GET", "reload", restConfigReload)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API to get value information about config items with type [ConstConfigItemGroupType]
func restConfigGroups(params *api.StructAPIHandlerParams) (interface{}, error) {

	config := env.GetConfig()
	return config.GetGroupItems(), nil
}

// WEB REST API to get value information about config items with type [ConstConfigItemGroupType]
func restConfigList(params *api.StructAPIHandlerParams) (interface{}, error) {

	config := env.GetConfig()
	return config.ListPathes(), nil
}

// WEB REST API to get value information about item(s) matching path
func restConfigInfo(params *api.StructAPIHandlerParams) (interface{}, error) {

	config := env.GetConfig()
	return config.GetItemsInfo(params.RequestURLParams["path"]), nil
}

// WEB REST API used to get value of particular item in config
//   - path should be without any wildcard
func restConfigGet(params *api.StructAPIHandlerParams) (interface{}, error) {

	config := env.GetConfig()
	return config.GetValue(params.RequestURLParams["path"]), nil
}

// WEB REST API used to set value of particular item in config
//   - path should be without any wildcard
func restConfigSet(params *api.StructAPIHandlerParams) (interface{}, error) {
	config := env.GetConfig()

	var setValue interface{}

	setValue = params.RequestContent
	configPath := params.RequestURLParams["path"]

	content, err := api.GetRequestContentAsMap(params)
	if err == nil {
		if contentValue, present := content["value"]; present {
			setValue = contentValue
		}
	}

	err = config.SetValue(configPath, setValue)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return config.GetValue(configPath), env.ErrorDispatch(err)
}

// WEB REST API used to add new config Item to a config system
func restConfigRegister(params *api.StructAPIHandlerParams) (interface{}, error) {
	inputData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	config := env.GetConfig()

	configItem := env.StructConfigItem{
		Path:  utils.InterfaceToString(utils.GetFirstMapValue(inputData, "path", "Path")),
		Value: utils.GetFirstMapValue(inputData, "value"),

		Type: utils.InterfaceToString(utils.GetFirstMapValue(inputData, "type", "Type")),

		Editor:  utils.InterfaceToString(utils.GetFirstMapValue(inputData, "editor", "Editor")),
		Options: utils.InterfaceToString(utils.GetFirstMapValue(inputData, "options", "Options")),

		Label:       utils.InterfaceToString(utils.GetFirstMapValue(inputData, "label", "Label")),
		Description: utils.InterfaceToString(utils.GetFirstMapValue(inputData, "description", "Description")),

		Image: utils.InterfaceToString(utils.GetFirstMapValue(inputData, "image", "Image")),
	}

	config.RegisterItem(configItem, nil)

	return configItem, nil
}

// WEB REST API used to remove config item from system
func restConfigUnRegister(params *api.StructAPIHandlerParams) (interface{}, error) {
	config := env.GetConfig()

	err := config.UnregisterItem(params.RequestURLParams["path"])
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API used to re-load config from DB
func restConfigReload(params *api.StructAPIHandlerParams) (interface{}, error) {
	config := env.GetConfig()

	err := config.Reload()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}
