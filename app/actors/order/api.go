package order

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("order", "GET", "attributes", restOrderAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("order", "GET", "list", restOrderList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("order", "POST", "list", restOrderList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("order", "GET", "count", restOrderCount)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("order", "GET", "get/:id", restOrderGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	// err = api.GetRestService().RegisterAPI("order", "POST", "add", restOrderAdd)
	// if err != nil {
	// 	return env.ErrorDispatch(err)
	// }
	err = api.GetRestService().RegisterAPI("order", "PUT", "update/:id", restOrderUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("order", "DELETE", "delete/:id", restOrderDelete)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to get order available attributes information
func restOrderAttributes(params *api.StructAPIHandlerParams) (interface{}, error) {

	orderModel, err := order.GetOrderModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return orderModel.GetAttributesInfo(), nil
}

// WEB REST API function used to obtain orders list
func restOrderList(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d4a758dd-3f03-44c5-b95e-c28a10aadf3c", "unexpected request content")
		}
		reqData = make(map[string]interface{})
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation start
	//----------------
	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// limit parameter handle
	orderCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, orderCollectionModel.GetDBCollection())

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := orderCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	return orderCollectionModel.List()
}

// WEB REST API function used to obtain orders count in model collection
func restOrderCount(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := orderCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	return dbCollection.Count()
}

// WEB REST API function to get order information
func restOrderGet(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqBlockID, present := params.RequestURLParams["id"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "723ef443-f974-4455-9be0-a8af13916554", "order id should be specified")
	}
	blockID := utils.InterfaceToString(reqBlockID)

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	orderModel, err := order.LoadOrderByID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := orderModel.ToHashMap()
	result["items"] = orderModel.GetItems()
	return result, nil
}

// WEB REST API for update existing order in system
func restOrderUpdate(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	blockID, present := params.RequestURLParams["id"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "20a08638-e9e6-428b-b70c-a418d7821e4b", "order id should be specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	orderModel, err := order.LoadOrderByID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		orderModel.Set(attribute, value)
	}

	orderModel.SetID(blockID)
	orderModel.Save()

	return orderModel.ToHashMap(), nil
}

// WEB REST API used to delete order from system
func restOrderDelete(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	blockID, present := params.RequestURLParams["id"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fc3011c7-e58c-4433-b9b0-881a7ba005cf", "order id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	orderModel, err := order.GetOrderModelAndSetID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderModel.Delete()

	return "ok", nil
}
