package saleprice

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	service := api.GetRestService()

	// Admin Only
	//-----------

	//service.GET("saleprices", api.IsAdmin(AdminAPIReadSalePriceList))
	//service.GET("saleprices/product/:id", api.IsAdmin(AdminAPIGetSalePriceListByProduct))

	service.POST("saleprice", api.IsAdmin(AdminAPICreateSalePrice))
	service.GET("saleprice/:id", api.IsAdmin(AdminAPIReadSalePrice))
	service.PUT("saleprice/:id", api.IsAdmin(AdminAPIUpdateSalePrice))
	service.DELETE("saleprice/:id", api.IsAdmin(AdminAPIDeleteSalePrice))

	return nil
}

// Returns list of all registered sale prices.
//func AdminAPIGetSalePriceList(context api.InterfaceApplicationContext) ([]map[string]interface{}, error) {
//
//	collection, err := db.GetCollection(ConstCollectionNameSalePrices)
//	if err != nil {
//		context.SetResponseStatusInternalServerError()
//		return nil, env.ErrorDispatch(err)
//	}
//
//	records, err := collection.Load()
//
//	return records, nil
//}

// Returns a list of registered sale prices for product
// * product id should be specified in the "product_id" argument
//func AdminAPIGetSalePriceListByProduct(context api.InterfaceApplicationContext) ([]map[string]interface{}, error) {
//
//	var postValues map[string]interface{}
//	var err error
//
//	if postValues, err = api.GetRequestContentAsMap(context); err != nil {
//		context.SetResponseStatusInternalServerError()
//		return nil, env.ErrorDispatch(err)
//	}
//
//	if !utils.KeysInMapAndNotBlank(postValues, "product_id") {
//		context.SetResponseStatusBadRequest()
//		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "5aea5a01-1601-40a3-9b4a-a6dcf7e4dff5", "Required field 'product_id' is not specified.")
//	}
//
//	var collection db.InterfaceDBCollection
//	if collection, err = db.GetCollection(ConstCollectionNameSalePrices); err != nil {
//		context.SetResponseStatusInternalServerError()
//		return nil, env.ErrorDispatch(err)
//	}
//
//	valueProductId := utils.InterfaceToString(postValues["product_id"])
//	if err = collection.AddFilter("product_id", "=", valueProductId); err != nil {
//		context.SetResponseStatusInternalServerError()
//		return nil, env.ErrorDispatch(err)
//	}
//
//	var records []map[string]interface{}
//	if records, err = collection.Load(); err != nil {
//		context.SetResponseStatusInternalServerError()
//		return nil, env.ErrorDispatch(err)
//	}
//
//	return records, nil
//}

// Check input parameters and store new Sale Price
func AdminAPICreateSalePrice(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(postValues, "amount", "start_datetime", "end_datetime", "product_id") {
		context.SetResponseStatusBadRequest()
		return nil, newErrorHelper(
			"Required fields 'amount', 'start_datetime', 'end_datetime', 'product_id', cannot be blank.",
			"a54d2879-d080-42fb-a733-1411911bd4d1")
	}

	// operation
	//----------
	return CreateSalePriceFromHashMapHelper(postValues)
}

// API returns a sale price with the specified ID
// * sale price id should be specified in the "id" argument
func AdminAPIReadSalePrice(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(postValues, "id") {
		context.SetResponseStatusBadRequest()
		return nil, newErrorHelper("Required field 'id' is blank or absend.", "5c00350a-1f06-4d7a-86a8-592d3c799f71")
	}

	// operation
	//-------------------------
	return ReadSalePriceHashMapByIDHelper(context.GetRequestArgument("id"))
}

// Update sale price
func AdminAPIUpdateSalePrice(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(postValues, "id") {
		context.SetResponseStatusBadRequest()
		return nil, newErrorHelper("Required field 'id' is blank or absend.", "dfd27d74-1d49-4baa-b64e-213324630765")
	}

	// operation
	//----------
	return UpdateSalePriceByHashMapHelper(context.GetRequestArgument("id"), postValues)
}

// Deletes specified sale price
func AdminAPIDeleteSalePrice(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(postValues, "id") {
		context.SetResponseStatusBadRequest()
		return nil, newErrorHelper("Required field 'id' is blank or absend.", "c91464a9-b920-4229-af4b-8a5945a862ca")
	}

	// operation
	//-------------------------
	err = DeleteSalePriceByIDHelper(context.GetRequestArgument("id"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "Delete Successful", nil
}
