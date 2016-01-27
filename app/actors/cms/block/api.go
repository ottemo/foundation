package block

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	if service := api.GetRestService(); service != nil {
		// print a stack trace to the logs if an error occurs
		handleError := func(err error) {
			if err != nil {
				env.ErrorDispatch(err)
			}
		}

		// package block endpoints
		handleError(service.RegisterAPI("cms/blocks", api.GET, APIListCMSBlocks))
		handleError(service.RegisterAPI("cms/blocks/attributes", api.GET, APIListCMSBlockAttributes))
		handleError(service.RegisterAPI("cms/block/:blockID", api.GET, APIGetCMSBlock))
		handleError(service.RegisterAPI("cms/block", api.POST, APICreateCMSBlock))
		handleError(service.RegisterAPI("cms/block/:blockID", api.PUT, APIUpdateCMSBlock))
		handleError(service.RegisterAPI("cms/block/:blockID", api.DELETE, APIDeleteCMSBlock))
	}

	return nil
}

// APIListCMSBlockAttributes returns a list of CMS block attributes
func APIListCMSBlockAttributes(context api.InterfaceApplicationContext) (interface{}, error) {

	cmsBlock, err := cms.GetCMSBlockModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cmsBlock.GetAttributesInfo(), nil
}

// APIListCMSBlocks returns a list of existing CMS blocks
//   - if "action" parameter is set to "count" result value will be just a number of list items
func APIListCMSBlocks(context api.InterfaceApplicationContext) (interface{}, error) {

	// taking CMS block collection model
	cmsBlockCollectionModel, err := cms.GetCMSBlockCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// applying request filters
	models.ApplyFilters(context, cmsBlockCollectionModel.GetDBCollection())

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return cmsBlockCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	cmsBlockCollectionModel.ListLimit(models.GetListLimit(context))

	// extra parameter handle
	models.ApplyExtraAttributes(context, cmsBlockCollectionModel)

	return cmsBlockCollectionModel.List()
}

// APIGetCMSBlock return specified CMS block information
//   - CMS block id should be specified in "blockID" argument
//   - CMS block content can be a text template, so "evaluated" field in response is that template evaluation result
func APIGetCMSBlock(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	reqBlockID := context.GetRequestArgument("blockID")
	if reqBlockID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a6dd2812-5070-4869-8ae2-90c4bd28bf69", "cms block id should be specified")
	}
	blockID := utils.InterfaceToString(reqBlockID)

	// operation
	//----------
	cmsBlock, err := cms.LoadCMSBlockByID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := cmsBlock.ToHashMap()
	result["evaluated"] = cmsBlock.EvaluateContent()

	return result, nil
}

// APICreateCMSBlock creates a new CMS block
//   - CMS block attributes should be specified in request content
func APICreateCMSBlock(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsBlockModel, err := cms.GetCMSBlockModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		cmsBlockModel.Set(attribute, value)
	}

	cmsBlockModel.SetID("")
	cmsBlockModel.Save()

	return cmsBlockModel.ToHashMap(), nil
}

// APIUpdateCMSBlock updates existing CMS block
//   - CMS block id should be specified in "blockID" argument
func APIUpdateCMSBlock(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	blockID := context.GetRequestArgument("blockID")
	if blockID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a7f8db95-7495-49ba-9307-baa7d5f7ecef", "cms block id should be specified")
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsBlockModel, err := cms.LoadCMSBlockByID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		cmsBlockModel.Set(attribute, value)
	}

	cmsBlockModel.SetID(blockID)
	cmsBlockModel.Save()

	return cmsBlockModel.ToHashMap(), nil
}

// APIDeleteCMSBlock deletes existing CMS block
//   - CMS block id should be specified in "blockID" argument
func APIDeleteCMSBlock(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	blockID := context.GetRequestArgument("blockID")
	if blockID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8dd275d4-efaf-4e67-b24d-67b28acd74e5", "cms block id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsBlockModel, err := cms.GetCMSBlockModelAndSetID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cmsBlockModel.Delete()

	return "ok", nil
}
