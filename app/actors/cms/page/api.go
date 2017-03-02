package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/env"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()
	service.GET("cms/pages", APIListCMSPages)
	service.GET("cms/pages/attributes", APIListCMSPageAttributes)
	service.GET("cms/page/:pageID", APIGetCMSPage)

	// Admin Only
	service.POST("cms/page", api.IsAdminHandler(APICreateCMSPage))
	service.PUT("cms/page/:pageID", api.IsAdminHandler(APIUpdateCMSPage))
	service.DELETE("cms/page/:pageID", api.IsAdminHandler(APIDeleteCMSPage))

	return nil
}

// APIListCMSPageAttributes returns a list of CMS block attributes
func APIListCMSPageAttributes(context api.InterfaceApplicationContext) (interface{}, error) {

	cmsPage, err := cms.GetCMSPageModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cmsPage.GetAttributesInfo(), nil
}

// APIListCMSPages returns a list of existing CMS pages
//   - if "action" parameter is set to "count" result value will be just a number of list items
func APIListCMSPages(context api.InterfaceApplicationContext) (interface{}, error) {

	cmsPageCollectionModel, err := cms.GetCMSPageCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// applying requested filters
	if err := models.ApplyFilters(context, cmsPageCollectionModel.GetDBCollection()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "91f07843-374a-43ba-89c6-9c1dd2c84b28", err.Error())
	}

	// excluding disabled pages for a regular visitor
	if !api.IsAdminSession(context) {
		if err := cmsPageCollectionModel.GetDBCollection().AddFilter("enabled", "=", true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "348d2715-84a3-43a6-a3d8-dc65a3fc3b88", err.Error())
		}
	}

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return cmsPageCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	if err := cmsPageCollectionModel.ListLimit(models.GetListLimit(context)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3d2c6ffd-5702-40f6-917d-90d302c0cd4d", err.Error())
	}

	// extra parameter handle
	if err := models.ApplyExtraAttributes(context, cmsPageCollectionModel); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0bfd8539-0c2b-4914-a379-ee5e92ec94ed", err.Error())
	}

	return cmsPageCollectionModel.List()
}

// APIGetCMSPage return specified CMS page information
//   - CMS page id should be specified in "pageID" argument
//   - CMS page content can be a text template, so "evaluated" field in response is that template evaluation result
func APIGetCMSPage(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	pageID := context.GetRequestArgument("pageID")
	if pageID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4c0a288a-03a4-4375-9fb3-1abd7981aecd", "cms page id should be specified")
	}

	// operation
	//----------
	cmsPage, err := cms.LoadCMSPageByID(pageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// not allowing to see disabled if not admin
	if !api.IsAdminSession(context) && !cmsPage.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fa76f5ac-0cce-4670-9e62-197a600ec0b9", "cms page is not available")
	}

	result := cmsPage.ToHashMap()
	result["evaluated"] = cmsPage.EvaluateContent()

	return result, nil
}

// APICreateCMSPage creates a new CMS block
//   - CMS page attributes should be specified in request content
func APICreateCMSPage(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsPageModel, err := cms.GetCMSPageModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		if err := cmsPageModel.Set(attribute, value); err != nil {
			_ = env.ErrorDispatch(err)
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d3339e7a-f651-48e8-8959-91a96471e788", "internal error")
		}
	}

	if err := cmsPageModel.SetID(""); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dc87c362-bfd2-4865-8cfa-4ec11ccc76f5", err.Error())
	}
	if err := cmsPageModel.Save(); err != nil {
		_ = env.ErrorDispatch(err)
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bb702fce-8a97-4913-aa8d-d2ceea3d6f56", "unable to save page")
	}

	return cmsPageModel.ToHashMap(), nil
}

// APIUpdateCMSPage updates existing CMS page
//   - CMS page id should be specified in "pageID" argument
func APIUpdateCMSPage(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	pageID := context.GetRequestArgument("pageID")
	if pageID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f128b02f-4ca5-494b-920d-5f320d112636", "cms page id should be specified")
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsPageModel, err := cms.LoadCMSPageByID(pageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		if err := cmsPageModel.Set(attribute, value); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f82fdb15-c2f6-407c-b5da-e81fae20068d", err.Error())
		}
	}

	if err := cmsPageModel.SetID(pageID); err != nil {
		_ = env.ErrorDispatch(err)
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8e660cca-8748-4759-808f-6531e11e6776", "internal error")
	}
	if err := cmsPageModel.Save(); err != nil {
		_ = env.ErrorDispatch(err)
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1fa828f2-ee96-4533-aab1-cd6c5ba6f570", "unable to save page")
	}

	return cmsPageModel.ToHashMap(), nil
}

// APIDeleteCMSPage deletes specified CMS page
//   - CMS page id should be specified in "pageID" argument
func APIDeleteCMSPage(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	pageID := context.GetRequestArgument("pageID")
	if pageID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "20545153-b171-4638-b47e-15317e85262c", "cms page id should be specified")
	}

	// operation
	//----------
	cmsPageModel, err := cms.GetCMSPageModelAndSetID(pageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if err := cmsPageModel.Delete(); err != nil {
		_ = env.ErrorDispatch(err)
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1d11f5ef-0fca-4d87-9c70-14cdd56a677e", "unable to delete page")
	}

	return "ok", nil
}
