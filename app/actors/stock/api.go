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

	service := api.GetRestService()

	service.GET("stock/:productID", APIGetProductStock)
	service.POST("stock/:productID/:qty", APISetStockQty)
	service.PUT("stock/:productID/:qty", APIUpdateStockQty)
	service.DELETE("stock/:productID", APIDeleteStockQty)

	service.POST("product/:productID/stock", APIGetProductQty)

	return nil
}

// APIGetProductStock returns stock information for particular product
//   - returns qty for all specified product-option pairs
//   - product id should be specified in "productID" argument
func APIGetProductStock(context api.InterfaceApplicationContext) (interface{}, error) {

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = dbCollection.AddFilter("product_id", "=", context.GetRequestArgument("productID"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return dbRecords, nil
}

// APIGetProductQty returns available stock qty for particular product-options pair
//   - product id should be specified in "productID" argument
//   - product options should be specified in "options" field of content
func APIGetProductQty(context api.InterfaceApplicationContext) (interface{}, error) {

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a4d6f5ae-35be-44c4-a1ff-6b7c17e05a73", "unexpected request content")
	}

	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "0d7805ef-e13e-47c3-8873-5ccf29749863", "no registered stock manager")
	}

	productID := context.GetRequestArgument("productID")
	options := make(map[string]interface{})
	if requestedOptions, present := requestData["options"]; present {
		options = utils.InterfaceToMap(requestedOptions)
	}

	return stockManager.GetProductQty(productID, options), nil
}

// APISetStockQty sets amount qty for a particular product-options pair
//   - product id and qty should be specified in "productID" and "qty" arguments
//   - product options should be specified in "options" field of content
func APISetStockQty(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b72e9050-cd30-4d9e-a3da-20b843fce518", "no registered stock manager")
	}

	productID := context.GetRequestArgument("productID")
	qty := utils.InterfaceToInt(context.GetRequestArgument("qty"))

	options := make(map[string]interface{})
	if requestedOptions, present := requestData["options"]; present {
		options = utils.InterfaceToMap(requestedOptions)
	}

	return stockManager.SetProductQty(productID, options, qty), nil
}

// APIUpdateStockQty increases qty on particular product-options pair for a delta value
//   - product id and delta should be specified in "productID" and "qty" arguments
//   - negative delta will decrease amount
//   - product options should be specified in "options" field of content
func APIUpdateStockQty(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c03d0b95-400e-415f-8c4a-26863993adbc", "no registered stock manager")
	}

	productID := context.GetRequestArgument("productID")
	qty := utils.InterfaceToInt(context.GetRequestArgument("qty"))

	options := make(map[string]interface{})
	if requestedOptions, present := requestData["options"]; present {
		options = utils.InterfaceToMap(requestedOptions)
	}

	return stockManager.UpdateProductQty(productID, options, qty), nil
}

// APIDeleteStockQty deletes stock records for a product-options pair
//   - product id should be specified in "productID" argument
//   - product options should be specified in "options" field of content
func APIDeleteStockQty(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9a4ca9aa-a12f-4913-9505-e05802a94c32", "no registered stock manager")
	}

	productID := context.GetRequestArgument("productID")

	options := make(map[string]interface{})
	if requestedOptions, present := requestData["options"]; present {
		options = utils.InterfaceToMap(requestedOptions)
	}

	return stockManager.RemoveProductQty(productID, options), nil
}
