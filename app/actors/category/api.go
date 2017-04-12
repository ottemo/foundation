package category

import (
	"io/ioutil"
	"mime"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	// Admin rights mix the response
	service.GET("categories", APIListCategories)
	service.GET("categories/tree", APIGetCategoriesTree)
	service.GET("categories/attributes", APIGetCategoryAttributes)

	service.GET("category/:categoryID", APIGetCategory)
	service.GET("category/:categoryID/layers", APIGetCategoryLayers)

	service.GET("category/:categoryID/products", APIGetCategoryProducts)

	service.GET("category/:categoryID/media/:mediaType/:mediaName", APIGetMedia)
	service.GET("category/:categoryID/media/:mediaType", APIListMedia)
	service.GET("category/:categoryID/mediapath/:mediaType", APIGetMediaPath)

	// Admin Only
	service.POST("category", api.IsAdminHandler(APICreateCategory))
	service.PUT("category/:categoryID", api.IsAdminHandler(APIUpdateCategory))
	service.DELETE("category/:categoryID", api.IsAdminHandler(APIDeleteCategory))

	service.POST("category/:categoryID/product/:productID", api.IsAdminHandler(APIAddProductToCategory))
	service.DELETE("category/:categoryID/product/:productID", api.IsAdminHandler(APIRemoveProductFromCategory))

	service.POST("category/:categoryID/media/:mediaType/:mediaName", api.IsAdminHandler(APIAddMediaForCategory))
	service.DELETE("category/:categoryID/media/:mediaType/:mediaName", api.IsAdminHandler(APIRemoveMediaForCategory))

	return nil
}

// APIListCategories returns a list of available categories
//   - if "action" parameter is set to "count" result value will be just a number of list items
//   - for a not admins available categories are limited to enabled ones
func APIListCategories(context api.InterfaceApplicationContext) (interface{}, error) {

	categoryCollectionModel, err := category.GetCategoryCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// applying requested filters
	if err := models.ApplyFilters(context, categoryCollectionModel.GetDBCollection()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7d867be7-731c-4af1-952f-23dd297e56c3", err.Error())
	}

	// excluding disabled categories for a regular visitor
	if !api.IsAdminSession(context) {
		if err := categoryCollectionModel.GetDBCollection().AddFilter("enabled", "=", true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4c40fe1d-34f3-4b21-8c53-e0a6d074eab0", err.Error())
		}
	}

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return categoryCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	if err := categoryCollectionModel.ListLimit(models.GetListLimit(context)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "eb982f7a-97b4-45f6-b5ad-84960d60050e", err.Error())
	}

	// extra parameter handle
	if err := models.ApplyExtraAttributes(context, categoryCollectionModel); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "deeb537f-2a4a-4808-9e59-e827160581f4", err.Error())
	}

	listItems, err := categoryCollectionModel.List()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var result []map[string]interface{}

	for _, listItem := range listItems {

		itemImages, err := mediaStorage.GetAllSizes(category.ConstModelNameCategory, listItem.ID, ConstCategoryMediaTypeImage)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// move default image to first position in array
		if listItem.Image != "" && len(itemImages) > 1 {
			found := false
			for index, images := range itemImages {
				for sizeName, sizeValue := range images {
					basicName := strings.Replace(sizeValue, "_"+sizeName, "", -1)
					if strings.Contains(basicName, listItem.Image) {
						found = true
						itemImages = append(itemImages[:index], itemImages[index+1:]...)
						itemImages = append([]map[string]string{images}, itemImages...)
					}
					break
				}
				if found {
					break
				}
			}
		}

		item := map[string]interface{}{
			"ID":     listItem.ID,
			"Name":   listItem.Name,
			"Desc":   listItem.Desc,
			"Extra":  listItem.Extra,
			"Image":  listItem.Image,
			"Images": itemImages,
		}

		result = append(result, item)
	}

	return result, nil
}

// APICreateCategory creates a new category
//   - category attributes must be provided in request content
//   - "name" attribute required
func APICreateCategory(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if _, present := requestData["name"]; !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "521b50d6-0d98-491a-8e3a-37678fbbccfe", "category name was not specified")
	}

	// create category operation
	//-------------------------
	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
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

// APIDeleteCategory deletes existing category
//   - category id should be specified in "categoryID" argument
func APIDeleteCategory(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//--------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f1507b00-337e-4903-8244-5cf959dde3a4", "category id was not specified")
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

// APIUpdateCategory modifies existing category
//   - category id must be specified as "id" argument
//   - category attributes must be specified in content
func APIUpdateCategory(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "389975e7-611c-4d6c-8b4d-bca450f5f7e7", "category id was not specified")
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// update operations
	//------------------
	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attrName, attrVal := range requestData {
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

// APIGetCategoryAttributes returns a list of category attributes
func APIGetCategoryAttributes(context api.InterfaceApplicationContext) (interface{}, error) {
	categoryModel, err := category.GetCategoryModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := categoryModel.GetAttributesInfo()

	return attrInfo, nil
}

// APIGetCategoryLayers enumerates category attributes and their possible values which is used for layered navigation
//   - category id should be specified in "id" argument
func APIGetCategoryLayers(context api.InterfaceApplicationContext) (interface{}, error) {
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3edd349a-269f-4400-8762-d685bb8e3a5b", "category id was not specified")
	}

	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// retrieve data only from attributes specified in "attributeCodes"
	attributeCodes := strings.Split(context.GetRequestArgument("attributeCodes"), ",")
	if len(attributeCodes) == 1 && attributeCodes[0] == "" {
		attributeCodes = []string{}
	}

	if !api.IsAdminSession(context) && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d46dadf8-373a-4247-a81e-fbbe39a7fe74", "category is not available")
	}

	productsCollection := categoryModel.GetProductsCollection()
	productsDBCollection := productsCollection.GetDBCollection()

	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	productAttributesInfo := productModel.GetAttributesInfo()

	result := make(map[string]interface{})

	if err := models.ApplyFilters(context, productsDBCollection); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "876d62b2-4d80-455c-8dcd-ac7971197df6", err.Error())
	}

	// not allowing to see disabled products if not admin
	if !api.IsAdminSession(context) {
		if err := productsDBCollection.AddFilter("enabled", "=", true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ea8e2ba1-c9df-484a-ac53-1b9fa43fcab1", err.Error())
		}
	}

	for _, productAttribute := range productAttributesInfo {
		if (productAttribute.IsLayered && (len(attributeCodes) == 0 || utils.IsInListStr(productAttribute.Attribute, attributeCodes))) {
			distinctValues, _ := productsDBCollection.Distinct(productAttribute.Attribute)
			result[productAttribute.Attribute] = distinctValues
		}
	}

	return result, nil
}

// APIGetCategoryProducts returns category related products
//   - category id should be specified in "categoryID" argument
func APIGetCategoryProducts(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "12400cff-34fe-4cf5-ac6e-41625f9e3d5a", "category id was not specified")
	}

	// product list operation
	categoryModel, err := category.LoadCategoryByID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !api.IsAdminSession(context) && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9a6f080d-dfa4-4f8c-8a0c-ec31cbe1cd87", "category is not available")
	}

	productsCollection := categoryModel.GetProductsCollection()

	if err := models.ApplyFilters(context, productsCollection.GetDBCollection()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ee3b2e74-81d4-460d-80ef-e9f973744ecd", err.Error())
	}

	// not allowing to see disabled and hidden products if not admin
	if !api.IsAdminSession(context) {
		if err := productsCollection.GetDBCollection().AddGroupFilter("visitor", "enabled", "=", true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "00051bb2-83a7-484f-8ad8-51697385afa1", err.Error())
		}
		if err := productsCollection.GetDBCollection().AddGroupFilter("visitor", "visible", "=", true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9b02437f-fceb-492e-a70e-f5cfb10d048d", err.Error())
		}
	}

	// checking for a "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return productsCollection.GetDBCollection().Count()
	}

	// limit parameter handle
	if err := productsCollection.ListLimit(models.GetListLimit(context)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d36c2098-5a68-4073-ae4b-a49ca56e9f27", err.Error())
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// preparing product information
	var result []map[string]interface{}

	for _, productModel := range productsCollection.ListProducts() {
		productInfo := productModel.ToHashMap()

		productInfo["image"], err = mediaStorage.GetAllSizes(product.ConstModelNameProduct, productModel.GetID(), ConstCategoryMediaTypeImage)
		if err != nil {
			_ = env.ErrorDispatch(err)
		}
		result = append(result, productInfo)
	}

	return result, nil
}

// APIAddProductToCategory adds product to category
//   - category id and product id  should be specified in "categoryID" and "productID" arguments
func APIAddProductToCategory(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2fc40bdc-3c8e-4a9c-910b-cca62cda1b77", "category id was not specified")
	}
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "190a3d60-7769-4908-b383-80bc143128da", "product id was not specified")
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

// APIRemoveProductFromCategory removes product from category
//   - category id and product id  should be specified in "categoryID" and "productID" arguments
func APIRemoveProductFromCategory(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "dc083799-bccb-48c8-bd56-5dcd0a0f6031", "category id was not specified")
	}
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9ffd2626-4192-4726-849e-c7a2416fab3a", "product id was not specified")
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

// APIGetCategory return specified category information
//   - category id should be specified in "categoryID" argument
func APIGetCategory(context api.InterfaceApplicationContext) (interface{}, error) {

	var categoryID string
	var categoryModel category.InterfaceCategory
	var err error

	// check request context
	if categoryID = context.GetRequestArgument("categoryID"); categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "3c336fd7-1a18-4aea-9eb0-460d746f8dfa", "category id was not specified")
	}

	// load category
	if categoryModel, err = category.LoadCategoryByID(categoryID); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !api.IsAdminSession(context) && !categoryModel.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "80615e04-f43d-42a4-9482-39a5e7f8ccb7", "category is not available")
	}

	result := categoryModel.ToHashMap()

	return result, nil
}

// APIGetCategoriesTree returns categories parent/child relation map
func APIGetCategoriesTree(context api.InterfaceApplicationContext) (interface{}, error) {

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
					return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c94e54a4-53dc-4ed2-bac4-7e9d93958765", "category tree builder internal error")
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

// APIGetMediaPath returns relative path to category media files within media library
//   - category id, media type must be specified in "categoryID" and "mediaType" arguments
func APIGetMediaPath(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "091d758d-9c3d-4e4f-8515-0f33a04b0354", "category id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if mediaType == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f5571ec1-a56d-4897-aaa8-6b14b466d82e", "media type was not specified")
	}

	// list media operation
	//---------------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	mediaList, err := categoryModel.GetMediaPath(mediaType)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaList, nil
}

// APIListMedia returns lost of media files assigned to specified category
//   - category id, media type must be specified in "categoryID" and "mediaType" arguments
func APIListMedia(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "363bf428-ba51-432e-9b80-7d4a0a11fe18", "category id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d05160cb-4df5-4917-bdb3-c06767d2483b", "media type was not specified")
	}

	// list media operation
	//---------------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	mediaList, err := categoryModel.ListMedia(mediaType)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaList, nil
}

// APIAddMediaForCategory uploads and assigns media file send in request for a specified category
//   - category id, media type and media name should be specified in "categoryID", "mediaType" and "mediaName" arguments
//   - media file should be provided in "file" field
func APIAddMediaForCategory(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "98e86d33-83c3-47c6-986b-11925c4f21ac", "category id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if mediaType == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "9db218b8-b37c-49a7-86d8-04a58d0f4fbc", "media type was not specified")
	}

	mediaName := context.GetRequestArgument("mediaName")
	if mediaName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "74cec364-c64d-493a-8848-e8af474eebc0", "media name was not specified")
	}

	// income file processing
	//-----------------------
	files := context.GetRequestFiles()
	if len(files) == 0 {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "75a2ddaf-b63d-4eed-b16d-4b32778f5fc1", "media file was not specified")
	}

	var fileContents []byte
	for _, fileReader := range files {
		contents, err := ioutil.ReadAll(fileReader)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		fileContents = contents
		break
	}

	// add media operation
	//--------------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Adding timestamp to image name to prevent overwriting
	mediaNameParts := strings.SplitN(mediaName, ".", 2)
	mediaName = mediaNameParts[0] + "_" + utils.InterfaceToString(time.Now().Unix()) + "." + mediaNameParts[1]

	err = categoryModel.AddMedia(mediaType, mediaName, fileContents)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIRemoveMediaForCategory removes media content from specified category
//   - category id, media type and media name should be specified in "categoryID", "mediaType" and "mediaName" arguments
func APIRemoveMediaForCategory(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "dfde9c02-97bb-4bc6-8874-669392ad32a1", "category id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if mediaType == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "224dd952-7c84-4ae4-882a-b15a002c37e6", "media type was not specified")
	}

	mediaName := context.GetRequestArgument("mediaName")
	if mediaName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "63b37b08-3b21-48b7-9058-291bb7e635a1", "media name was not specified")
	}

	// list media operation
	//---------------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.RemoveMedia(mediaType, mediaName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIGetMedia returns media contents for a category (file assigned to a category)
//   - category id, media type and media name must be specified in "categoryID", "mediaType" and "mediaName" arguments
//   - on success case not a JSON data returns, but media file
func APIGetMedia(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	categoryID := context.GetRequestArgument("categoryID")
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "41d5d10d-9c15-4144-b1e2-ebf47d1911d8", "category id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if mediaType == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "453e4717-8101-41b4-91c1-172d1b489662", "media type was not specified")
	}

	mediaName := context.GetRequestArgument("mediaName")
	if mediaName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "124c8b9d-1a6b-491c-97ba-a03e8c828337", "media name was not specified")
	}

	if err := context.SetResponseContentType(mime.TypeByExtension(mediaName)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3d8c7a09-3ee7-489f-899c-d6c32c3f5285", err.Error())
	}

	// list media operation
	//---------------------
	categoryModel, err := category.GetCategoryModelAndSetID(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return categoryModel.GetMedia(mediaType, mediaName)
}
