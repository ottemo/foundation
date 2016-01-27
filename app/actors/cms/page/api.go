package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/env"
)

// setupAPI configures package related API endpoint routines
func setupAPI() error {

	if service := api.GetRestService(); service != nil {
		// print a stack trace to the logs if an error occurs
		handleError := func(err error) {
			if err != nil {
				env.ErrorDispatch(err)
			}
		}

		// package page endpoints
		handleError(service.RegisterAPI("cms/pages", api.GET, APIListCMSPages))
		handleError(service.RegisterAPI("cms/pages/attributes", api.GET, APIListCMSPageAttributes))
		handleError(service.RegisterAPI("cms/page/:pageID", api.GET, APIGetCMSPage))
		handleError(service.RegisterAPI("cms/page", api.POST, APICreateCMSPage))
		handleError(service.RegisterAPI("cms/page/:pageID", api.PUT, APIUpdateCMSPage))
		handleError(service.RegisterAPI("cms/page/:pageID", api.DELETE, APIDeleteCMSPage))
	}

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
	models.ApplyFilters(context, cmsPageCollectionModel.GetDBCollection())

	// excluding disabled pages for a regular visitor
	if err := api.ValidateAdminRights(context); err != nil {
		cmsPageCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return cmsPageCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	cmsPageCollectionModel.ListLimit(models.GetListLimit(context))

	// extra parameter handle
	models.ApplyExtraAttributes(context, cmsPageCollectionModel)

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
	if api.ValidateAdminRights(context) != nil && !cmsPage.GetEnabled() {
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

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsPageModel, err := cms.GetCMSPageModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		cmsPageModel.Set(attribute, value)
	}

	cmsPageModel.SetID("")
	cmsPageModel.Save()

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

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsPageModel, err := cms.LoadCMSPageByID(pageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		cmsPageModel.Set(attribute, value)
	}

	cmsPageModel.SetID(pageID)
	cmsPageModel.Save()

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

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
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
