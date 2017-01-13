package product

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"mime"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/subscription"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	// Public
	service.GET("products", APIListProducts)
	service.GET("product/:productID", APIGetProduct)

	service.GET("products/attributes", APIListProductAttributes)

	service.GET("product/:productID/media/:mediaType/:mediaName", APIGetMedia) // @DEPRECATED
	service.GET("product/:productID/media/:mediaType", APIListMedia)           // @DEPRECATED
	service.GET("product/:productID/mediapath/:mediaType", APIGetMediaPath)    // @DEPRECATED

	// Related
	service.GET("product/:productID/related", APIListRelatedProducts)

	// Admin Only
	service.POST("product", api.IsAdmin(APICreateProduct))
	service.PUT("product/:productID", api.IsAdmin(APIUpdateProduct))
	service.DELETE("product/:productID", api.IsAdmin(APIDeleteProduct))

	service.POST("products/attribute", api.IsAdmin(APICreateProductAttribute))
	service.PUT("products/attribute/:attribute", api.IsAdmin(APIUpdateProductAttribute))
	service.DELETE("products/attribute/:attribute", api.IsAdmin(APIDeleteProductsAttribute))

	service.POST("product/:productID/media/:mediaType/:mediaName", api.IsAdmin(APIAddMediaForProduct))
	service.DELETE("product/:productID/media/:mediaType/:mediaName", api.IsAdmin(APIRemoveMediaForProduct))
	service.PUT("product/:productID/media/:mediaType/:mediaName", api.IsAdmin(APIRenameMediaForProduct))

	// TODO: remove after patching
	service.GET("patch/options", api.IsAdmin(APIPatchOptions))

	return nil
}

// APIPatchOptions converts product options to snake case in products and subscriptions
// TODO: remove after patching
func APIPatchOptions(context api.InterfaceApplicationContext) (interface{}, error) {
	warnings := make(map[string]string)

	// get product collection
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return warnings, env.ErrorDispatch(err)
	}

	// update products option
	for _, currentProduct := range productCollection.ListProducts() {
		newOptions := ConvertProductOptionsToSnakeCase(currentProduct)
		err = currentProduct.Set("options", newOptions)
		if err != nil {
			warnings["product ID: "+currentProduct.GetID()] = utils.InterfaceToString(err)
		}

		err := currentProduct.Save()
		if err != nil {
			return warnings, env.ErrorDispatch(err)
		}
	}

	// get subscriptions collection
	subscriptionCollection, err := subscription.GetSubscriptionCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	currentCart, err := cart.GetCartModel()
	if err != nil {
		return warnings, env.ErrorDispatch(err)
	}

	for _, currentSubscription := range subscriptionCollection.ListSubscriptions() {
		var updatedItems []subscription.StructSubscriptionItem
		for _, subscriptionItem := range currentSubscription.GetItems() {
			updatedOptions := make(map[string]interface{})
			// Labels where used as a key for options key: value, so we will convert both of them
			for optionKey, optionValue := range subscriptionItem.Options {
				updatedOptions[utils.StrToSnakeCase(optionKey)] = utils.StrToSnakeCase(utils.InterfaceToString(optionValue))
			}
			subscriptionItem.Options = updatedOptions
			if _, err = currentCart.AddItem(subscriptionItem.ProductID, subscriptionItem.Qty, subscriptionItem.Options); err != nil {
				warnings["subscription ID: "+currentSubscription.GetID()] = utils.InterfaceToString(err)
			}

			updatedItems = append(updatedItems, subscriptionItem)
		}

		currentSubscription.Set("items", updatedItems)

		err = currentSubscription.Save()
		if err != nil {
			return warnings, env.ErrorDispatch(err)
		}
	}

	return warnings, nil
}

// ConvertProductOptionsToSnakeCase updates option keys for product to case_snake
// TODO: remove after patching
func ConvertProductOptionsToSnakeCase(product product.InterfaceProduct) map[string]interface{} {

	newOptions := make(map[string]interface{})

	// product options
	for optionsName, currentOption := range product.GetOptions() {
		currentOption := utils.InterfaceToMap(currentOption)

		if option, present := currentOption["options"]; present {
			newOptionValues := make(map[string]interface{})

			// option values
			for key, value := range utils.InterfaceToMap(option) {
				newOptionValues[utils.StrToSnakeCase(key)] = value

			}

			currentOption["options"] = newOptionValues

		}
		newOptions[utils.StrToSnakeCase(optionsName)] = currentOption

	}

	return newOptions
}

// APIListProductAttributes returns a list of product attributes
func APIListProductAttributes(context api.InterfaceApplicationContext) (interface{}, error) {
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := productModel.GetAttributesInfo()

	return attrInfo, nil
}

// APIUpdateProductAttribute updates existing custom attribute of product model
//   - attribute name/code should be provided in "attribute" argument
//   - attribute parameters should be provided in request content
//   - attribute parameters "id" and "name" will be ignored
//   - static attributes can not be changed
func APIUpdateProductAttribute(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attributeName := context.GetRequestArgument("attribute")
	if attributeName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cb8f7251-e22b-4605-97bb-e239df6c7aac", "attribute name was not specified")
	}

	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, attribute := range productModel.GetAttributesInfo() {
		if attribute.Attribute == attributeName {
			if attribute.IsStatic == true {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2893262f-a61a-42f8-9c75-e763e0a5c8ca", "can't edit static attributes")
			}

			for key, value := range requestData {
				switch strings.ToLower(key) {
				case "label":
					attribute.Label = utils.InterfaceToString(value)
				case "group":
					attribute.Group = utils.InterfaceToString(value)
				case "editors":
					attribute.Editors = utils.InterfaceToString(value)
				case "options":
					attribute.Options = utils.InterfaceToString(value)
				case "default":
					attribute.Default = utils.InterfaceToString(value)
				case "validators":
					attribute.Validators = utils.InterfaceToString(value)
				case "isrequired", "required":
					attribute.IsRequired = utils.InterfaceToBool(value)
				case "islayered", "layered":
					attribute.IsLayered = utils.InterfaceToBool(value)
				case "ispublic", "public":
					attribute.IsPublic = utils.InterfaceToBool(value)
				}
			}
			err := productModel.EditAttribute(attributeName, attribute)
			if err != nil {
				return nil, err
			}
			return attribute, nil
		}
	}

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2893262f-a61a-42f8-9c75-e763e0a5c8ca", "attribute not found")
}

// APICreateProductAttribute creates a new custom attribute for a product model
//   - attribute parameters "Attribute" and "Label" are required
func APICreateProductAttribute(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attributeName, isSpecified := requestData["Attribute"]
	if !isSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2f7aec81-dba8-4cad-b683-23c5d0a08cf5", "attribute name was not specified")
	}

	attributeLabel, isSpecified := requestData["Label"]
	if !isSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "93457847-8e4d-4536-8985-43f340a1abc4", "attribute label was not specified")
	}

	// make product attribute operation
	//---------------------------------
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attribute := models.StructAttributeInfo{
		Model:      product.ConstModelNameProduct,
		Collection: ConstCollectionNameProduct,
		Attribute:  utils.InterfaceToString(attributeName),
		Type:       utils.ConstDataTypeText,
		IsRequired: false,
		IsStatic:   false,
		Label:      utils.InterfaceToString(attributeLabel),
		Group:      "General",
		Editors:    "text",
		Options:    "",
		Default:    "",
		Validators: "",
		IsLayered:  false,
		IsPublic:   false,
	}

	for key, value := range requestData {
		switch strings.ToLower(key) {
		case "type":
			attribute.Type = utils.InterfaceToString(value)
		case "group":
			attribute.Group = utils.InterfaceToString(value)
		case "editors":
			attribute.Editors = utils.InterfaceToString(value)
		case "options":
			attribute.Options = utils.InterfaceToString(value)
		case "default":
			attribute.Default = utils.InterfaceToString(value)
		case "validators":
			attribute.Validators = utils.InterfaceToString(value)
		case "isrequired", "required":
			attribute.IsRequired = utils.InterfaceToBool(value)
		case "islayered", "layered":
			attribute.IsLayered = utils.InterfaceToBool(value)
		case "ispublic", "public":
			attribute.IsPublic = utils.InterfaceToBool(value)
		}
	}

	err = productModel.AddNewAttribute(attribute)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return attribute, nil
}

// APIDeleteProductsAttribute removes existing custom attribute of a product model
//   - attribute name/code should be provided in "attribute" argument
func APIDeleteProductsAttribute(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//--------------------
	attributeName := context.GetRequestArgument("attribute")
	if attributeName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cb8f7251-e22b-4605-97bb-e239df6c7aac", "attribute name was not specified")
	}

	// remove attribute actions
	//-------------------------
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.RemoveAttribute(attributeName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIGetProduct return specified product information
//   - product id should be specified in "productID" argument
func APIGetProduct(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "feb3a463-622b-477e-a22d-c0a3fd1972dc", "product id was not specified")
	}

	// load product operation
	//-----------------------
	productModel, err := product.LoadProductByID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// not allowing to see disabled products if not admin
	if api.ValidateAdminRights(context) != nil && (!productModel.GetEnabled() || !utils.InterfaceToBool(productModel.Get("visible"))) {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "153673ac-1008-40b5-ada9-2286ad3f02b0", "product not available")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// get product
	//-------------

	result := productModel.ToHashMap()

	itemImages, err := mediaStorage.GetAllSizes(product.ConstModelNameProduct, productModel.GetID(), ConstProductMediaTypeImage)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	defaultImage := productModel.GetDefaultImage()

	// move default image to first position in array
	if defaultImage != "" && len(itemImages) > 1 {
		defaultImageName := defaultImage[strings.LastIndex(defaultImage, "/")+1 : strings.Index(defaultImage, ".")]
		found := false
		for index, images := range itemImages {
			for _, sizeValue := range images {
				if strings.Contains(sizeValue, defaultImageName) {
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

	result["images"] = itemImages

	return result, nil
}

// APICreateProduct creates a new product
//   - product attributes must be provided in request content
//   - "sku" and "name" attributes are required
func APICreateProduct(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(context.GetRequestArguments(), "sku", "name") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2a0cf2b0-215e-4b53-bf55-98fbfe22cd27", "product name and/or sku were not specified")
	}

	// create product operation
	//-------------------------
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range requestData {
		err := productModel.Set(attribute, value)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = productModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return productModel.ToHashMap(), nil
}

// APIDeleteProduct deletes existing product
//   - product id must be specified in "productID" argument
func APIDeleteProduct(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//--------------------
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f35af170-8172-4ec0-b30d-ab883231d222", "product id was not specified")
	}

	// delete operation
	//-----------------
	productModel, err := product.GetProductModelAndSetID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.Delete()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIUpdateProduct updates existing product
//   - product id should be specified in "productID" argument
//   - product attributes should be specified in content
func APIUpdateProduct(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c91e8fc7-ca77-40d1-823c-e50f90b8b4b5", "product id was not specified")
	}

	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fffccbad-455a-4fff-81d4-8919ae3a5c35", "unexpected request content")
	}

	// update operations
	//------------------
	productModel, err := product.LoadProductByID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attrName, attrVal := range requestData {
		err = productModel.Set(attrName, attrVal)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	err = productModel.Save()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := productModel.ToHashMap()

	return result, nil
}

// APIGetMediaPath returns relative path to product media files within media library
//   - product id, media type must be specified in "productID" and "mediaType" arguments
func APIGetMediaPath(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6597ff92-f2ee-4233-bcf9-eb73b957fb05", "product id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if mediaType == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "75c00741-5873-4be1-9fa0-df9d2956d3de", "media type was not specified")
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	mediaList, err := productModel.GetMediaPath(mediaType)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaList, nil
}

// APIListMedia returns lost of media files assigned to specified product
//   - product id, media type must be specified in "productID" and "mediaType" arguments
func APIListMedia(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "52677450-8a7f-49c9-a472-51d0e80bc7ca", "product id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b8b31a9f-6fac-47b3-89e2-c9b3e589a8f6", "media type was not specified")
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	mediaList, err := productModel.ListMedia(mediaType)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaList, nil
}

// APIAddMediaForProduct uploads and assigns media file send in request for a specified product
//   - product id, media type and media name should be specified in "productID", "mediaType" and "mediaName" arguments
//   - media file should be provided in "file" field
func APIAddMediaForProduct(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a4696c5d-3276-4272-8d86-8061e57743a5", "product id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if mediaType == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f3ea9a01-412a-4af2-9496-cb58cdb8139d", "media type was not specified")
	}

	mediaName := context.GetRequestArgument("mediaName")
	if mediaName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "23fb7617-f19a-4505-b706-10f7898fd980", "media name was not specified")
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
	productModel, err := product.GetProductModelAndSetID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// Adding timestamp to image name to prevent overwriting
	mediaNameParts := strings.SplitN(mediaName, ".", 2)
	mediaName = mediaNameParts[0] + "_" + utils.InterfaceToString(time.Now().Unix()) + "." + mediaNameParts[1]

	err = productModel.AddMedia(mediaType, mediaName, fileContents)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIRemoveMediaForProduct removes media content from specified product
//   - product id, media type and media name should be specified in "productID", "mediaType" and "mediaName" arguments
func APIRemoveMediaForProduct(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f5f77b7f-6606-4bdd-a113-0a3b26f5759c", "product id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if mediaType == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e81b841f-8253-4b66-ac7d-2cc9a484044c", "media type was not specified")
	}

	mediaName := context.GetRequestArgument("mediaName")
	if mediaName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "63b37b08-3b21-48b7-9058-291bb7e635a1", "media name was not specified")
	}

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.RemoveMedia(mediaType, mediaName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIRenameMediaForProduct renames media file for a specified product
//   - product id, media type and media name should be specified in "productID", "mediaType" and "mediaName" arguments
//   - new media name should be provided in "newMediaName" field
func APIRenameMediaForProduct(context api.InterfaceApplicationContext) (interface{}, error) {
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check request context
	//---------------------
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2dc543fb-0fa9-4900-93f1-56031fa68dc1", "product id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if mediaType == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "91ea7d50-8577-4b9b-ae04-3b5022d189c9", "media type was not specified")
	}

	mediaName := context.GetRequestArgument("mediaName")
	if mediaName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "379426c8-b1c1-44ee-b2e0-e0704e6f6e0f", "media name was not specified")
	}

	newMediaName, present := requestData["newMediaName"]
	if !present {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8205e4f3-0805-42cd-a6c3-b8c6d139d52a", "new media name was not specified")
	}

	// add media operation
	//--------------------
	productModel, err := product.GetProductModelAndSetID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// load media
	fileContents, err := productModel.GetMedia(mediaType, mediaName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// remove media
	err = productModel.RemoveMedia(mediaType, mediaName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check image format
	decodedImage, imageFormat, err := image.Decode(bytes.NewReader(fileContents))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// convert image format if not default
	if imageFormat != ConstSwatchImageDefaultFormat {
		buffer := bytes.NewBuffer(nil)
		err = jpeg.Encode(buffer, decodedImage, nil)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		fileContents = buffer.Bytes()
	}

	mediaName = utils.InterfaceToString(newMediaName) + "_" + utils.InterfaceToString(time.Now().Unix()) + "." + ConstSwatchImageDefaultExtention

	err = productModel.AddMedia(mediaType, mediaName, fileContents)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIGetMedia returns media contents for a product (file assigned to a product)
//   - product id, media type and media name must be specified in "productID", "mediaType" and "mediaName" arguments
//   - on success case not a JSON data returns, but media file
func APIGetMedia(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d33b8a67-359f-4a3e-b626-f58b6c70f09f", "product id was not specified")
	}

	mediaType := context.GetRequestArgument("mediaType")
	if mediaType == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d081b726-caf4-4694-baaa-7b1801ca9713", "media type was not specified")
	}

	mediaName := context.GetRequestArgument("mediaName")
	if mediaName == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "124c8b9d-1a6b-491c-97ba-a03e8c828337", "media name was not specified")
	}

	context.SetResponseContentType(mime.TypeByExtension(mediaName))

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return productModel.GetMedia(mediaType, mediaName)
}

// APIListProducts returns a list of available products
//   - if "action" parameter is set to "count" result value will be just a number of list items
//   - visitors can not see disabled products, but administrators can
func APIListProducts(context api.InterfaceApplicationContext) (interface{}, error) {

	var productCollectionModel product.InterfaceProductCollection
	var err error

	if productCollectionModel, err = product.GetProductCollectionModel(); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// filters handle
	models.ApplyFilters(context, productCollectionModel.GetDBCollection())

	// exclude disabled and hidden products for visitors, but not Admins
	if err := api.ValidateAdminRights(context); err != nil {
		productCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
		productCollectionModel.GetDBCollection().AddFilter("visible", "=", true)
	}

	// check "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return productCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	productCollectionModel.ListLimit(models.GetListLimit(context))

	// extra parameter handle
	models.ApplyExtraAttributes(context, productCollectionModel)

	listItems, err := productCollectionModel.List()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	var result []map[string]interface{}

	for _, listItem := range listItems {

		itemImages, err := mediaStorage.GetAllSizes(product.ConstModelNameProduct, listItem.ID, ConstProductMediaTypeImage)
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// move default image to first position in array
		if listItem.Image != "" && len(itemImages) > 1 {
			defaultImageName := listItem.Image[strings.LastIndex(listItem.Image, "/")+1 : strings.Index(listItem.Image, ".")]
			found := false
			for index, images := range itemImages {
				for _, sizeValue := range images {
					if strings.Contains(sizeValue, defaultImageName) {
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

// APIListRelatedProducts returns related products list for a given product
func APIListRelatedProducts(context api.InterfaceApplicationContext) (interface{}, error) {

	var result []map[string]interface{}

	productID := context.GetRequestArgument("productID")
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "55aa2eee-0407-4094-a90a-5d69d8c1efcc", "product id was not specified")
	}

	productModel, err := product.LoadProductByID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	relatedPids := utils.InterfaceToArray(productModel.Get("related_pids"))

	productsCollection, _ := product.GetProductCollectionModel()
	productsCollection.GetDBCollection().AddFilter("_id", "in", relatedPids)

	// if you aren't an admin the product must be enabled
	if err := api.ValidateAdminRights(context); err != nil {
		productsCollection.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// add a limit
	productsCollection.ListLimit(models.GetListLimit(context))

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, relatedProduct := range productsCollection.ListProducts() {
		productInfo := relatedProduct.ToHashMap()

		defaultImage := utils.InterfaceToString(productInfo["default_image"])
		productInfo["image"], err = mediaStorage.GetSizes(product.ConstModelNameProduct, relatedProduct.GetID(), ConstProductMediaTypeImage, defaultImage)
		if err != nil {
			env.ErrorDispatch(err)
		}

		result = append(result, productInfo)
	}

	return result, nil
}
