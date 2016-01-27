package config

import (
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("config/groups", api.GET, restConfigGroups)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config/item/:path", api.GET, restConfigInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config/values", api.GET, restConfigList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config/values/refresh", api.GET, restConfigReload)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config/value/:path", api.GET, restConfigGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config/value/:path", api.POST, restConfigRegister)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config/value/:path", api.PUT, restConfigSet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("config/value/:path", api.DELETE, restConfigUnRegister)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API to get value information about config items with type [ConstConfigTypeGroup]
func restConfigGroups(context api.InterfaceApplicationContext) (interface{}, error) {

	config := env.GetConfig()
	return config.GetGroupItems(), nil
}

// WEB REST API to get value information about config items with type [ConstConfigTypeGroup]
func restConfigList(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	config := env.GetConfig()
	return config.ListPathes(), nil
}

// WEB REST API to get value information about item(s) matching path
func restConfigInfo(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	config := env.GetConfig()
	return config.GetItemsInfo(context.GetRequestArgument("path")), nil
}

// WEB REST API used to get value of particular item in config
//   - path should be without any wildcard
func restConfigGet(context api.InterfaceApplicationContext) (interface{}, error) {
	config := env.GetConfig()

	configItemPath := context.GetRequestArgument("path")

	info := config.GetItemsInfo(configItemPath)
	if len(info) == 1 {
		itemInfo := info[0]

		if itemInfo.Type == env.ConstConfigTypeSecret ||
			strings.Contains(itemInfo.Editor, "password") ||
			strings.Contains(itemInfo.Type, "password") ||
			strings.Contains(itemInfo.Path, "password") ||
			strings.Contains(itemInfo.Path, "login") ||
			strings.Contains(itemInfo.Path, "admin") {

			// check rights
			if err := api.ValidateAdminRights(context); err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	return config.GetValue(context.GetRequestArgument("path")), nil
}

// WEB REST API used to set value of particular item in config
//   - path should be without any wildcard
func restConfigSet(context api.InterfaceApplicationContext) (interface{}, error) {
	config := env.GetConfig()

	var setValue interface{}

	setValue = context.GetRequestContent()
	configPath := context.GetRequestArgument("path")

	content, err := api.GetRequestContentAsMap(context)
	if err == nil {
		if contentValue, present := content["value"]; present {
			setValue = contentValue
		}
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = config.SetValue(configPath, setValue)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return config.GetValue(configPath), env.ErrorDispatch(err)
}

// WEB REST API used to add new config Item to a config system
func restConfigRegister(context api.InterfaceApplicationContext) (interface{}, error) {
	inputData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
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
func restConfigUnRegister(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	config := env.GetConfig()

	err := config.UnregisterItem(context.GetRequestArgument("path"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API used to re-load config from DB
func restConfigReload(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	config := env.GetConfig()

	err := config.Reload()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}
