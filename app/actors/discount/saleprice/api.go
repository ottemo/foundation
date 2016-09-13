package saleprice

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	service := api.GetRestService()

	// Admin Only
	//-----------

	service.GET("saleprices", api.IsAdmin(AdminAPIReadSalePriceList))

	service.POST("saleprice", api.IsAdmin(AdminAPICreateSalePrice))
	service.GET("saleprice/:id", api.IsAdmin(AdminAPIReadSalePrice))
	service.PUT("saleprice/:id", api.IsAdmin(AdminAPIUpdateSalePrice))
	service.DELETE("saleprice/:id", api.IsAdmin(AdminAPIDeleteSalePrice))

	return nil
}

// AdminAPIReadSalePriceList returns list of all registered sale prices.
func AdminAPIReadSalePriceList(context api.InterfaceApplicationContext) (interface{}, error) {
	salePriceCollectionModel, err := saleprice.GetSalePriceCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// applying requested filters
	models.ApplyFilters(context, salePriceCollectionModel.GetDBCollection())

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return salePriceCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	salePriceCollectionModel.ListLimit(models.GetListLimit(context))

	// extra parameter handle
	models.ApplyExtraAttributes(context, salePriceCollectionModel)

	listItems, err := salePriceCollectionModel.List()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return listItems, nil
}

// AdminAPICreateSalePrice checks input parameters and store new Sale Price
func AdminAPICreateSalePrice(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "amount", "start_datetime", "end_datetime", "product_id") {
		context.SetResponseStatusBadRequest()
		return nil, newErrorHelper(
			"Required fields 'amount', 'start_datetime', 'end_datetime', 'product_id', cannot be blank.",
			"a54d2879-d080-42fb-a733-1411911bd4d1")
	}

	// operation
	//----------
	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		err := salePriceModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = salePriceModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return salePriceModel.ToHashMap(), nil
}

// AdminAPIReadSalePrice returns a sale price with the specified ID
// * sale price id should be specified in the "id" argument
func AdminAPIReadSalePrice(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	salePriceID := context.GetRequestArgument("id")
	if salePriceID == "" {
		return nil, newErrorHelper("Required field 'id' is blank or absend.", "beb06bd0-db31-4daa-9fdd-d9872da7fdd6")
	}

	// operation
	//-------------------------
	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = salePriceModel.Load(salePriceID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return salePriceModel.ToHashMap(), nil
}

// AdminAPIUpdateSalePrice updates sale price
func AdminAPIUpdateSalePrice(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	salePriceID := context.GetRequestArgument("id")
	if salePriceID == "" {
		return nil, newErrorHelper("Required field 'id' is blank or absend.", "beb06bd0-db31-4daa-9fdd-d9872da7fdd6")
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = salePriceModel.Load(salePriceID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attrName, attrVal := range requestData {
		err = salePriceModel.Set(attrName, attrVal)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = salePriceModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return salePriceModel.ToHashMap(), nil
}

// AdminAPIDeleteSalePrice deletes specified sale price
func AdminAPIDeleteSalePrice(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	salePriceID := context.GetRequestArgument("id")
	if salePriceID == "" {
		return nil, newErrorHelper("Required field 'id' is blank or absend.", "beb06bd0-db31-4daa-9fdd-d9872da7fdd6")
	}

	// operation
	//-------------------------
	salePriceModel, err := saleprice.GetSalePriceModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = salePriceModel.Load(salePriceID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = salePriceModel.Delete()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "Delete Successful", nil
}
