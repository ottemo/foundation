package block

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("cms", "GET", "block/attributes", restCMSBlockAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "GET", "block/list", restCMSBlockList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "POST", "block/list", restCMSBlockList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "GET", "block/count", restCMSBlockCount)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "GET", "block/get/:id", restCMSBlockGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "POST", "block/add", restCMSBlockAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "PUT", "block/update/:id", restCMSBlockUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "DELETE", "block/delete/:id", restCMSBlockDelete)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to get CMS block available attributes information
func restCMSBlockAttributes(params *api.StructAPIHandlerParams) (interface{}, error) {

	cmsBlock, err := cms.GetCMSBlockModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cmsBlock.GetAttributesInfo(), nil
}

// WEB REST API function used to obtain CMS blocks list
func restCMSBlockList(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b61f75c4-2be1-4547-8764-d7040b1edb2f", "unexpected request content")
		}
		reqData = make(map[string]interface{})
	}

	// operation start
	//----------------
	cmsBlockCollectionModel, err := cms.GetCMSBlockCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// limit parameter handle
	cmsBlockCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, cmsBlockCollectionModel.GetDBCollection())

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := cmsBlockCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	return cmsBlockCollectionModel.List()
}

// WEB REST API function used to obtain CMS blocks count in model collection
func restCMSBlockCount(params *api.StructAPIHandlerParams) (interface{}, error) {

	cmsBlockCollectionModel, err := cms.GetCMSBlockCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := cmsBlockCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	return dbCollection.Count()
}

// WEB REST API function to get CMS block information
func restCMSBlockGet(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqBlockID, present := params.RequestURLParams["id"]
	if !present {
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

// WEB REST API for adding new CMS block in system
func restCMSBlockAdd(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
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
	cmsBlockModel, err := cms.GetCMSBlockModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		cmsBlockModel.Set(attribute, value)
	}

	cmsBlockModel.SetID("")
	cmsBlockModel.Save()

	return cmsBlockModel.ToHashMap(), nil
}

// WEB REST API for update existing CMS block in system
func restCMSBlockUpdate(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	blockID, present := params.RequestURLParams["id"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a7f8db95-7495-49ba-9307-baa7d5f7ecef", "cms block id should be specified")
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
	cmsBlockModel, err := cms.LoadCMSBlockByID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		cmsBlockModel.Set(attribute, value)
	}

	cmsBlockModel.SetID(blockID)
	cmsBlockModel.Save()

	return cmsBlockModel.ToHashMap(), nil
}

// WEB REST API used to delete CMS block from system
func restCMSBlockDelete(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	blockID, present := params.RequestURLParams["id"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8dd275d4-efaf-4e67-b24d-67b28acd74e5", "cms block id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
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
