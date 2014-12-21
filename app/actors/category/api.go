package category

import (
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("category", "GET", "list", restListCategories)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "list", restListCategories)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "count", restCountCategories)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "create", restCreateCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "PUT", "update/:id", restUpdateCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "DELETE", "delete/:id", restDeleteCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "get/:id", restGetCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "layers/:id", restListCategoryLayers)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "layers/:id", restListCategoryLayers)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("category", "GET", "product/list/:categoryID", restListCategoryProducts)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "POST", "product/list/:categoryID", restListCategoryProducts)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "product/count/:categoryID", restCategoryProductsCount)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "product/add/:categoryID/:productID", restAddCategoryProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("category", "GET", "product/remove/:categoryID/:productID", restRemoveCategoryProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("category", "GET", "attribute/list", restListCategoryAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("category", "GET", "tree", restGetCategoriesTree)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function used to obtain category list we have in database
//   - parent categories and categorys will not be present in list
func restListCategories(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "760bdd87468c4e6f913103e9df4fc6ce", "unexpected request content")
		}
		reqData = make(map[string]interface{})
	}

	// operation start
	//----------------
	categoryCollectionModel, err := category.GetCategoryCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// limit parameter handler
	categoryCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, categoryCollectionModel.GetDBCollection())

	// not allowing to see disabled if not admin
	if err := api.ValidateAdminRights(params); err != nil {
		categoryCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// extra parameter handler
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := categoryCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	return categoryCollectionModel.List()
}

// WEB REST API function used to obtain categories count in model collection
func restCountCategories(params *api.StructAPIHandlerParams) (interface{}, error) {
	categoryCollectionModel, err := category.GetCategoryCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := categoryCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	// not allowing to see disabled if not admin
	if err := api.ValidateAdminRights(params); err != nil {
		categoryCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	return dbCollection.Count()
}

// WEB REST API used to create new category
//   - category attributes must be included in POST form
//   - name attribute required
func restCreateCategory(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if _, present := reqData["name"]; !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "521b50d60d98491a8e3a37678fbbccfe", "category name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// create category operation
	//-------------------------
	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		err := categoryModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = categoryModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return categoryModel.ToHashMap(), nil
}

// WEB REST API used to delete category
func restDeleteCategory(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	categoryID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f1507b00337e490382445cf959dde3a4", "category id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// delete operation
	//-----------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.Delete()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API used to update existing category
//   - category id must be specified in request URI
//   - category attributes must be included in POST form
func restUpdateCategory(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	categoryID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "389975e7611c4d6c8b4dbca450f5f7e7", "category id was not specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// update operations
	//------------------
	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attrName, attrVal := range reqData {
		err = categoryModel.Set(attrName, attrVal)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = categoryModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return categoryModel.ToHashMap(), nil
}

// WEB REST API function used to obtain category attributes information
func restListCategoryAttributes(params *api.StructAPIHandlerParams) (interface{}, error) {
	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := categoryModel.GetAttributesInfo()

	return attrInfo, nil
}

// WEB REST API function used to obtain layered navigation options for products in category
//   - category id must be specified in request URI
func restListCategoryLayers(params *api.StructAPIHandlerParams) (interface{}, error) {
	categoryID := params.RequestURLParams["id"]

	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if api.ValidateAdminRights(params) != nil && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d46dadf8373a4247a81efbbe39a7fe74", "category is not available")
	}

	productsCollection := categoryModel.GetProductsCollection()
	productsDBCollection := productsCollection.GetDBCollection()

	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	productAttributesInfo := productModel.GetAttributesInfo()

	result := make(map[string]interface{})

	api.ApplyFilters(params, productsDBCollection)

	// not allowing to see disabled products if not admin
	if err := api.ValidateAdminRights(params); err != nil {
		productsDBCollection.AddFilter("enabled", "=", true)
	}

	for _, productAttribute := range productAttributesInfo {
		if productAttribute.IsLayered {
			distinctValues, _ := productsDBCollection.Distinct(productAttribute.Attribute)
			result[productAttribute.Attribute] = distinctValues
		}
	}

	return result, nil
}

// WEB REST API function used to list product in category
//   - category id must be specified in request URI
func restListCategoryProducts(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	categoryID, isSpecifiedID := params.RequestURLParams["categoryID"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "12400cff34fe4cf5ac6e41625f9e3d5a", "category id was not specified")
	}

	// product list operation
	//-----------------------
	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if api.ValidateAdminRights(params) != nil && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9a6f080ddfa44f8c8a0cec31cbe1cd87", "category is not available")
	}

	productsCollection := categoryModel.GetProductsCollection()

	api.ApplyFilters(params, productsCollection.GetDBCollection())

	// not allowing to see disabled products if not admin
	if err := api.ValidateAdminRights(params); err != nil {
		productsCollection.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// preparing product information
	var result []map[string]interface{}

	for _, productModel := range productsCollection.ListProducts() {
		productInfo := productModel.ToHashMap()
		if defaultImage, present := productInfo["default_image"]; present {
			mediaPath, err := productModel.GetMediaPath("image")
			if defaultImage, ok := defaultImage.(string); ok && defaultImage != "" && err == nil {
				productInfo["default_image"] = mediaPath + defaultImage
			}
		}
		result = append(result, productInfo)
	}

	return result, nil
}

// WEB REST API function used to add product in category
//   - category and product ids must be specified in request URI
func restAddCategoryProduct(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	categoryID, isSpecifiedID := params.RequestURLParams["categoryID"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2fc40bdc3c8e4a9c910bcca62cda1b77", "category id was not specified")
	}
	productID, isSpecifiedID := params.RequestURLParams["productID"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "190a3d6077694908b38380bc143128da", "product id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// category product add operation
	//-------------------------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.AddProduct(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API function used to remove product from category
//   - category and product ids must be specified in request URI
func restRemoveCategoryProduct(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	categoryID, isSpecifiedID := params.RequestURLParams["categoryID"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "dc083799bccb48c8bd565dcd0a0f6031", "category id was not specified")
	}
	productID, isSpecifiedID := params.RequestURLParams["productID"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9ffd262641924726849ec7a2416fab3a", "product id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// category product add operation
	//-------------------------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.RemoveProduct(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API function used to obtain all product attributes
//   - product id must be specified in request URI "http://[site:port]/product/get/:id"
func restGetCategory(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	categoryID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3c336fd71a184aea9eb0460d746f8dfa", "category id was not specified")
	}

	// load product operation
	//-----------------------
	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if api.ValidateAdminRights(params) != nil && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "80615e04f43d42a4948239a5e7f8ccb7", "category is not available")
	}

	return categoryModel.ToHashMap(), nil
}

// WEB REST API function used to list product in category
//   - category id must be specified in request URI
func restCategoryProductsCount(params *api.StructAPIHandlerParams) (interface{}, error) {

	categoryID, isSpecifiedID := params.RequestURLParams["categoryID"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e1003839d91c444596c830e481b8347e", "category id was not specified")
	}

	// product list operation
	//-----------------------
	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if api.ValidateAdminRights(params) != nil && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9c2efb6acd284dd2ba7f5c1fe7c0df30", "category is not available")
	}

	productsDBCollection := categoryModel.GetProductsCollection().GetDBCollection()
	api.ApplyFilters(params, productsDBCollection)

	// not allowing to see disabled products if not admin
	if err := api.ValidateAdminRights(params); err != nil {
		productsDBCollection.AddFilter("enabled", "=", true)
	}

	return productsDBCollection.Count()
}

// WEB REST API function used to categories menu
func restGetCategoriesTree(params *api.StructAPIHandlerParams) (interface{}, error) {

	var result = make([]map[string]interface{}, 0)

	collection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = collection.AddFilter("enabled", "=", true)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = collection.AddSort("path", false)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	rowData, err := collection.Load()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var categoryStack []map[string]interface{}
	var pathStack []string

	for _, row := range rowData {

		currentItem := make(map[string]interface{})
		currentItem["id"] = row["_id"]
		currentItem["name"] = row["name"]
		currentItem["child"] = make([]map[string]interface{}, 0)

		// calculating current path
		currentPath, ok := row["path"].(string)
		if !ok {
			continue
		}

		for idx := len(pathStack) - 1; idx >= 0; idx-- {
			parentPath := pathStack[idx]

			// if we found parent
			if strings.Contains(currentPath, parentPath) {
				parent := categoryStack[idx]

				parentChild, ok := parent["child"].([]map[string]interface{})
				if !ok {
					return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c94e54a453dc4ed2bac47e9d93958765", "category tree builder internal error")
				}

				parent["child"] = append(parentChild, currentItem)

				break
			} else {
				pathStack = pathStack[0:idx]
				categoryStack = categoryStack[0:idx]
			}
		}

		if len(categoryStack) == 0 {
			result = append(result, currentItem)
		}

		categoryStack = append(categoryStack, currentItem)
		pathStack = append(pathStack, currentPath)

	}

	return result, nil
}
