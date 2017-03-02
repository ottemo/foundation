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

	service := api.GetRestService()

	service.GET("cms/blocks", APIListCMSBlocks)
	service.GET("cms/blocks/attributes", APIListCMSBlockAttributes)
	service.GET("cms/block/:blockID", APIGetCMSBlock)

	// Admin Only
	service.POST("cms/block", api.IsAdminHandler(APICreateCMSBlock))
	service.PUT("cms/block/:blockID", api.IsAdminHandler(APIUpdateCMSBlock))
	service.DELETE("cms/block/:blockID", api.IsAdminHandler(APIDeleteCMSBlock))

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
	if err := models.ApplyFilters(context, cmsBlockCollectionModel.GetDBCollection()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "64d4f10c-680e-43e6-9d52-2d8710607dc2", err.Error())
	}

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return cmsBlockCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	if err := cmsBlockCollectionModel.ListLimit(models.GetListLimit(context)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0cc02545-ee3f-4816-a67f-adf2a64c267c", err.Error())
	}

	// extra parameter handle
	if err := models.ApplyExtraAttributes(context, cmsBlockCollectionModel); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "49379308-09cc-4df9-87c4-757ba9a2484a", err.Error())
	}

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

	// operation
	//----------
	cmsBlockModel, err := cms.GetCMSBlockModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		if err := cmsBlockModel.Set(attribute, value); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c1a265dd-453f-40a7-b8ff-a82a0f1ab066", err.Error())
		}
	}

	if err := cmsBlockModel.SetID(""); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a06ff37b-9572-4d14-a6ad-b8123dffc996", err.Error())
	}
	if err := cmsBlockModel.Save(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d3b8ea22-b012-4f75-85c1-d5905edfd625", err.Error())
	}

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

	// operation
	//----------
	cmsBlockModel, err := cms.LoadCMSBlockByID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		if err := cmsBlockModel.Set(attribute, value); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "01656722-c2e5-41d0-b4b3-78ec33f1235b", err.Error())
		}
	}

	if err := cmsBlockModel.SetID(blockID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7f360300-9762-4b05-9448-2a95771668e4", err.Error())
	}
	if err := cmsBlockModel.Save(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "22785238-740c-49f7-852d-7d9bc498ae31", err.Error())
	}

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

	// operation
	//----------
	cmsBlockModel, err := cms.GetCMSBlockModelAndSetID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if err := cmsBlockModel.Delete(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fe2262d8-d386-4c4a-8edc-5e8388dffd7a", err.Error())
	}

	return "ok", nil
}
