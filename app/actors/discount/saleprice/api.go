package saleprice

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {
	service := api.GetRestService()

	// Admin Only
	//-----------

	service.GET("saleprices", api.IsAdmin(listAllScheduled))

	service.POST("saleprice", api.IsAdmin(createSalePrice))
	service.GET("saleprice/:id", api.IsAdmin(priceByID))
	service.PUT("saleprice/:id", api.IsAdmin(updateByID))
	service.DELETE("saleprice/:id", api.IsAdmin(deleteByID))

	return nil
}

// listAllScheduled returns list of all registered sale price promotions scheduled.
func listAllScheduled(context api.InterfaceApplicationContext) (interface{}, error) {
	salePriceCollectionModel, err := saleprice.GetSalePriceCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// applying requested filters
	if err := models.ApplyFilters(context, salePriceCollectionModel.GetDBCollection()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c19d40b6-dc72-4ee1-b851-40cdf3636991", err.Error())
	}

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return salePriceCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	if err := salePriceCollectionModel.ListLimit(models.GetListLimit(context)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "34a048a7-e6f6-4041-b9af-6c884fd74f09", err.Error())
	}

	// extra parameter handle
	if err := models.ApplyExtraAttributes(context, salePriceCollectionModel); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0d1cc9dd-d7c4-4190-9fb7-dac7aa69d3fb", err.Error())
	}

	listItems, err := salePriceCollectionModel.List()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return listItems, nil
}

// createSalePrice checks input parameters and store new Sale Price
func createSalePrice(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(requestData, "amount", "start_datetime", "product_id") {
		context.SetResponseStatusBadRequest()
		return nil, newErrorHelper(
			"Required fields 'amount', 'start_datetime', 'product_id', cannot be blank.",
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

// priceByID returns a sale price with the specified ID
// * sale price id should be specified in the "id" argument
func priceByID(context api.InterfaceApplicationContext) (interface{}, error) {

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

// updateByID updates sale price
func updateByID(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	salePriceID := context.GetRequestArgument("id")
	if salePriceID == "" {
		return nil, newErrorHelper("Required field 'id' is blank or absend.", "fc0e40f9-c51e-46a7-b53b-c480be5bd556")
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

// deleteByID deletes specified sale price
func deleteByID(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	salePriceID := context.GetRequestArgument("id")
	if salePriceID == "" {
		return nil, newErrorHelper("Required field 'id' is blank or absend.", "1cb07d89-5447-4de5-b15a-3a8e76f8a818")
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
