package stock

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("stock", "GET", "info/:productID", restStockInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("stock", "POST", "get/:productID", restStockGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("stock", "POST", "set/:productID/:qty", restStockSet)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("stock", "PUT", "update/:productID/:delta", restStockUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("stock", "DELETE", "remove/:productID", restStockRemove)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API used obtain stock information on particular product-options pair
func restStockInfo(params *api.StructAPIHandlerParams) (interface{}, error) {

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		env.ErrorDispatch(err)
		return nil, err
	}

	err = dbCollection.AddFilter("product_id", "=", params.RequestURLParams["productID"])
	if err != nil {
		env.ErrorDispatch(err)
		return nil, err
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		env.ErrorDispatch(err)
		return nil, err
	}

	return dbRecords, nil
}

// WEB REST API used get available qty on particular product-options pair
func restStockGet(params *api.StructAPIHandlerParams) (interface{}, error) {

	requestData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a4d6f5ae35be44c4a1ff6b7c17e05a73", "unexpected request content")
	}

	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0d7805efe13e47c388735ccf29749863", "no registered stock manager")
	}

	productID := params.RequestURLParams["productID"]
	options := make(map[string]interface{})
	if requestedOptions, present := requestData["options"]; present {
		options = utils.InterfaceToMap(requestedOptions)
	}

	return stockManager.GetProductQty(productID, options), nil
}

// WEB REST API used set available qty on particular product-options pair
func restStockSet(params *api.StructAPIHandlerParams) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b72e9050cd304d9ea3da20b843fce518", "no registered stock manager")
	}

	productID := params.RequestURLParams["productID"]
	qty := utils.InterfaceToInt(params.RequestURLParams["qty"])

	options := make(map[string]interface{})
	if requestedOptions, present := requestData["options"]; present {
		options = utils.InterfaceToMap(requestedOptions)
	}

	return stockManager.SetProductQty(productID, options, qty), nil
}

// WEB REST API used to increase available qty on particular product-options pair for a delta value
func restStockUpdate(params *api.StructAPIHandlerParams) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c03d0b95400e415f8c4a26863993adbc", "no registered stock manager")
	}

	productID := params.RequestURLParams["productID"]
	qty := utils.InterfaceToInt(params.RequestURLParams["delta"])

	options := make(map[string]interface{})
	if requestedOptions, present := requestData["options"]; present {
		options = utils.InterfaceToMap(requestedOptions)
	}

	return stockManager.UpdateProductQty(productID, options, qty), nil
}

// WEB REST API used to remove product-options qty data from stock
func restStockRemove(params *api.StructAPIHandlerParams) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9a4ca9aaa12f49139505e05802a94c32", "no registered stock manager")
	}

	productID := params.RequestURLParams["productID"]

	options := make(map[string]interface{})
	if requestedOptions, present := requestData["options"]; present {
		options = utils.InterfaceToMap(requestedOptions)
	}

	return stockManager.RemoveProductQty(productID, options), nil
}
