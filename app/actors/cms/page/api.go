package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/cms"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("cms", "GET", "page/attributes", restCMSPageAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "GET", "page/list", restCMSPageList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "POST", "page/list", restCMSPageList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "GET", "page/count", restCMSPageCount)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "GET", "page/get/:id", restCMSPageGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "POST", "page/add", restCMSPageAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "PUT", "page/update/:id", restCMSPageUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms", "DELETE", "page/delete/:id", restCMSPageDelete)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to get CMS page available attributes information
func restCMSPageAttributes(params *api.StructAPIHandlerParams) (interface{}, error) {

	cmsPage, err := cms.GetCMSPageModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cmsPage.GetAttributesInfo(), nil
}

// WEB REST API function used to obtain CMS pages list
func restCMSPageList(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4813d3ab-67e4-4abf-b1f8-906e75cdd907", "unexpected request content")
		}
		reqData = make(map[string]interface{})
	}

	// operation start
	//----------------
	cmsPageCollectionModel, err := cms.GetCMSPageCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// limit parameter handle
	cmsPageCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, cmsPageCollectionModel.GetDBCollection())

	// not allowing to see disabled if not admin
	if err := api.ValidateAdminRights(params); err != nil {
		cmsPageCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := cmsPageCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	return cmsPageCollectionModel.List()
}

// WEB REST API function used to obtain CMS pages count in model collection
func restCMSPageCount(params *api.StructAPIHandlerParams) (interface{}, error) {

	cmsPageCollectionModel, err := cms.GetCMSPageCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := cmsPageCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	// not allowing to see disabled if not admin
	if err := api.ValidateAdminRights(params); err != nil {
		dbCollection.AddFilter("enabled", "=", true)
	}

	return dbCollection.Count()
}

// WEB REST API function to get CMS page information
func restCMSPageGet(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqPageID, present := params.RequestURLParams["id"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4c0a288a-03a4-4375-9fb3-1abd7981aecd", "cms page id should be specified")
	}
	pageID := utils.InterfaceToString(reqPageID)

	// operation
	//----------
	cmsPage, err := cms.LoadCMSPageByID(pageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// not allowing to see disabled if not admin
	if api.ValidateAdminRights(params) != nil && !cmsPage.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fa76f5ac-0cce-4670-9e62-197a600ec0b9", "cms page is not available")
	}

	result := cmsPage.ToHashMap()
	result["evaluated"] = cmsPage.EvaluateContent()

	return result, nil
}

// WEB REST API for adding new CMS page in system
func restCMSPageAdd(params *api.StructAPIHandlerParams) (interface{}, error) {

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
	cmsPageModel, err := cms.GetCMSPageModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		cmsPageModel.Set(attribute, value)
	}

	cmsPageModel.SetID("")
	cmsPageModel.Save()

	return cmsPageModel.ToHashMap(), nil
}

// WEB REST API for update existing CMS page in system
func restCMSPageUpdate(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	pageID, present := params.RequestURLParams["id"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f128b02f-4ca5-494b-920d-5f320d112636", "cms page id should be specified")
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
	cmsPageModel, err := cms.LoadCMSPageByID(pageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		cmsPageModel.Set(attribute, value)
	}

	cmsPageModel.SetID(pageID)
	cmsPageModel.Save()

	return cmsPageModel.ToHashMap(), nil
}

// WEB REST API used to delete CMS page from system
func restCMSPageDelete(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	pageID, present := params.RequestURLParams["id"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "20545153-b171-4638-b47e-15317e85262c", "cms page id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsPageModel, err := cms.GetCMSPageModelAndSetID(pageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cmsPageModel.Delete()

	return "ok", nil
}
