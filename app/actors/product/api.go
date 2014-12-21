package product

import (
	"math/rand"
	"mime"
	"strings"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	// 1. DefaultProduct API
	//----------------------
	//TODO: shorten endpoint to just 'product/:id' as GET verb is enough - jwv
	err = api.GetRestService().RegisterAPI("product", "GET", "get/:id", restGetProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	//TODO: shorten endpoint to just 'product' as POST verb assumes creation - jwv
	err = api.GetRestService().RegisterAPI("product", "POST", "create", restCreateProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	//TODO: shorten endpoint to just 'product/:id' as PUT verb describes it as an update - jwv
	err = api.GetRestService().RegisterAPI("product", "PUT", "update/:id", restUpdateProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	//TODO: shorten endpoint to just 'product/:id' as verb DELETE describes delete action - jwv
	err = api.GetRestService().RegisterAPI("product", "DELETE", "delete/:id", restDeleteProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	//TODO: shorten endpoint to just 'product/attribute/:attribute' as DELETE verb indicates purpose - jwv
	//TODO: shorten endpoint to just 'product/attribute' - jwv
	err = api.GetRestService().RegisterAPI("product", "GET", "attribute/list", restListProductAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "DELETE", "attribute/remove/:attribute", restRemoveProductAttribute)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "attribute/add", restAddProductAttribute)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "attribute/edit/:attribute", restEditProductAttribute)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	//TODO: shorten endpoint to just 'product/media/:productID/"mediaType/:mediaName' - jwv
	err = api.GetRestService().RegisterAPI("product", "GET", "media/get/:productID/:mediaType/:mediaName", restMediaGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "media/list/:productID/:mediaType", restMediaList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "media/path/:productID/:mediaType", restMediaPath)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	//TODO: remove 'add' from endpoint URL - jwv
	err = api.GetRestService().RegisterAPI("product", "POST", "media/add/:productID/:mediaType/:mediaName", restMediaAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	//TODO: remove 'remove' from endpoint URL - jwv
	err = api.GetRestService().RegisterAPI("product", "DELETE", "media/remove/:productID/:mediaType/:mediaName", restMediaRemove)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	//TODO: change to a GET as we are retrieving a list not creating one - jwv
	err = api.GetRestService().RegisterAPI("product", "POST", "related/:productID", restRelatedList)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// 2. DefaultProductCollection API
	//--------------------------------
	err = api.GetRestService().RegisterAPI("product", "GET", "list", restListProducts)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "POST", "list", restListProducts)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("product", "GET", "count", restCountProducts)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

//----------------------
// 1. DefaultProduct API
//----------------------

// WEB REST API function used to obtain product attributes information
func restListProductAttributes(params *api.StructAPIHandlerParams) (interface{}, error) {
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attrInfo := productModel.GetAttributesInfo()

	return attrInfo, nil
}

// WEB REST API function used to edit existing custom attribute fields (except id and name)
func restEditProductAttribute(params *api.StructAPIHandlerParams) (interface{}, error) {
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attributeName, isSpecified := params.RequestURLParams["attribute"]
	if !isSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cb8f7251e22b460597bbe239df6c7aac", "attribute name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, attribute := range productModel.GetAttributesInfo() {
		if attribute.Attribute == attributeName {
			if attribute.IsStatic == true {
				return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2893262fa61a42f89c75e763e0a5c8ca", "can't edit static attributes")
			}

			for key, value := range reqData {
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

	return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2893262fa61a42f89c75e763e0a5c8ca", "attribute not found")
}

// WEB REST API function used to add new one custom attribute
func restAddProductAttribute(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	attributeName, isSpecified := reqData["Attribute"]
	if !isSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2f7aec81dba84cadb68323c5d0a08cf5", "attribute name was not specified")
	}

	attributeLabel, isSpecified := reqData["Label"]
	if !isSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "934578478e4d4536898543f340a1abc4", "attribute label was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
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
		Type:       "text",
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

	for key, value := range reqData {
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

// WEB REST API function used to remove custom attribute of product
func restRemoveProductAttribute(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	attributeName, isSpecified := params.RequestURLParams["attribute"]
	if !isSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "cb8f7251e22b460597bbe239df6c7aac", "attribute name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
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

// WEB REST API function used to obtain all product attributes
//   - product id must be specified in request URI "http://[site:port]/product/get/:id"
func restGetProduct(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "feb3a463622b477ea22dc0a3fd1972dc", "product id was not specified")
	}

	// load product operation
	//-----------------------
	productModel, err := product.LoadProductByID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// not allowing to see disabled products if not admin
	if api.ValidateAdminRights(params) != nil && productModel.GetEnabled() == false {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "153673ac100840b5ada92286ad3f02b0", "product not available")
	}

	return productModel.ToHashMap(), nil
}

// WEB REST API used to create new one product
//   - product attributes must be included in POST form
//   - sku and name attributes required
func restCreateProduct(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(params.RequestURLParams, "sku", "name") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "2a0cf2b0215e4b53bf5598fbfe22cd27", "product name and/or sku were not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// create product operation
	//-------------------------
	productModel, err := product.GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
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

// WEB REST API used to delete product
//   - product attributes must be included in POST form
func restDeleteProduct(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//--------------------
	productID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f35af17081724ec0b30dab883231d222", "product id was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
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

// WEB REST API used to update existing product
//   - product id must be specified in request URI
//   - product attributes must be included in POST form
func restUpdateProduct(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isSpecifiedID := params.RequestURLParams["id"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c91e8fc7ca7740d1823ce50f90b8b4b5", "product id was not specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fffccbad455a4fff81d48919ae3a5c35", "unexpected request content")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// update operations
	//------------------
	productModel, err := product.LoadProductByID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attrName, attrVal := range reqData {
		err = productModel.Set(attrName, attrVal)
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

// WEB REST API used to add media for a product
//   - product id, media type must be specified in request URI
func restMediaPath(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isIDSpecified := params.RequestURLParams["productID"]
	if !isIDSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "6597ff92f2ee4233bcf9eb73b957fb05", "product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "75c0074158734be19fa0df9d2956d3de", "media type was not specified")
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

// WEB REST API used to add media for a product
//   - product id, media type must be specified in request URI
func restMediaList(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isIDSpecified := params.RequestURLParams["productID"]
	if !isIDSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "526774508a7f49c9a47251d0e80bc7ca", "product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "b8b31a9f6fac47b389e2c9b3e589a8f6", "media type was not specified")
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

// WEB REST API used to add media for a product
//   - product id, media type and media name must be specified in request URI
//   - media contents must be included as file in POST form
func restMediaAdd(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isIDSpecified := params.RequestURLParams["productID"]
	if !isIDSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a4696c5d327642728d868061e57743a5", "product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f3ea9a01412a4af29496cb58cdb8139d", "media type was not specified")
	}

	mediaName, isNameSpecified := params.RequestURLParams["mediaName"]
	if !isNameSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "23fb7617f19a4505b70610f7898fd980", "media name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// income file processing
	//-----------------------
	file, _, err := params.Request.FormFile("file")
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	fileSize, _ := file.Seek(0, 2)
	fileContents := make([]byte, fileSize)

	file.Seek(0, 0)
	file.Read(fileContents)

	// add media operation
	//--------------------
	productModel, err := product.GetProductModelAndSetID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.AddMedia(mediaType, mediaName, fileContents)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// WEB REST API used to add media for a product
//   - product id, media type and media name must be specified in request URI
func restMediaRemove(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isIDSpecified := params.RequestURLParams["productID"]
	if !isIDSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f5f77b7f66064bdda1130a3b26f5759c", "product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "e81b841f82534b66ac7d2cc9a484044c", "media type was not specified")
	}

	mediaName, isNameSpecified := params.RequestURLParams["mediaName"]
	if !isNameSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "63b37b083b2148b79058291bb7e635a1", "media name was not specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, env.ErrorDispatch(err)
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

// WEB REST API used to get media contents for a product
//   - product id, media type and media name must be specified in request URI
func restMediaGet(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isIDSpecified := params.RequestURLParams["productID"]
	if !isIDSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d33b8a67359f4a3eb626f58b6c70f09f", "product id was not specified")
	}

	mediaType, isTypeSpecified := params.RequestURLParams["mediaType"]
	if !isTypeSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "d081b726caf44694baaa7b1801ca9713", "media type was not specified")
	}

	mediaName, isNameSpecified := params.RequestURLParams["mediaName"]
	if !isNameSpecified {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "124c8b9d1a6b491c97baa03e8c828337", "media name was not specified")
	}

	params.ResponseWriter.Header().Set("Content-Type", mime.TypeByExtension(mediaName))

	// list media operation
	//---------------------
	productModel, err := product.GetProductModelAndSetID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return productModel.GetMedia(mediaType, mediaName)
}

//--------------------------------
// 2. DefaultProductCollection API
//--------------------------------

// WEB REST API function used to obtain product list we have in database
//   - only [_id, sku, name] attributes returns by default
func restListProducts(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// api function routine
	//---------------------
	productCollectionModel, err := product.GetProductCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// limit parameter handle
	productCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, productCollectionModel.GetDBCollection())

	// not allowing to see disabled products if not admin
	if err := api.ValidateAdminRights(params); err != nil {
		productCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := productCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	return productCollectionModel.List()
}

func restRelatedList(params *api.StructAPIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	productID, isSpecifiedID := params.RequestURLParams["productID"]
	if !isSpecifiedID {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "55aa2eee04074094a90a5d69d8c1efcc", "product id was not specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	count := 5
	if utils.InterfaceToInt(reqData["count"]) > 0 {
		count = utils.InterfaceToInt(reqData["count"])
	}

	// load product operation
	//-----------------------
	productModel, err := product.LoadProductByID(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// result := make([]models.StructListItem, 0)
	var result []models.StructListItem

	relatedPids := utils.InterfaceToArray(productModel.Get("related_pids"))

	if len(relatedPids) > count {
		// indexes := make([]int, 0)
		var indexes []int
		for len(indexes) < count {

			new := rand.Intn(len(relatedPids))

			inArray := false
			for _, b := range indexes {
				if utils.InterfaceToInt(b) == new {
					inArray = true
				}
			}
			if !inArray {
				indexes = append(indexes, new)
			}
		}
		for _, index := range indexes {
			if productID := utils.InterfaceToString(relatedPids[index]); productID != "" {
				if productModel, err := product.LoadProductByID(productID); err == nil {

					// not allowing to see disabled products if not admin
					if err := api.ValidateAdminRights(params); err != nil {
						if productModel.GetEnabled() == false {
							continue
						}
					}

					if err == nil {
						resultItem := new(models.StructListItem)

						mediaPath, err := productModel.GetMediaPath("image")
						if err != nil {
							return result, env.ErrorDispatch(err)
						}

						resultItem.ID = productModel.GetID()
						resultItem.Name = "[" + productModel.GetSku() + "] " + productModel.GetName()
						resultItem.Image = ""
						resultItem.Desc = productModel.GetShortDescription()

						if productModel.GetDefaultImage() != "" {
							resultItem.Image = mediaPath + productModel.GetDefaultImage()
						}

						if extra, isExtra := reqData["extra"]; isExtra {
							resultItem.Extra = make(map[string]interface{})
							extra := utils.Explode(utils.InterfaceToString(extra), ",")
							for _, value := range extra {
								resultItem.Extra[value] = productModel.Get(value)
							}
						}

						result = append(result, *resultItem)
					}
				}
			}
		}
	} else {
		for _, productID := range relatedPids {
			if productID == "" {
				continue
			}

			productModel, err := product.LoadProductByID(utils.InterfaceToString(productID))
			if err == nil {
				resultItem := new(models.StructListItem)

				mediaPath, err := productModel.GetMediaPath("image")
				if err != nil {
					return result, env.ErrorDispatch(err)
				}

				resultItem.ID = productModel.GetID()
				resultItem.Name = "[" + productModel.GetSku() + "] " + productModel.GetName()
				resultItem.Image = ""
				resultItem.Desc = productModel.GetShortDescription()

				if productModel.GetDefaultImage() != "" {
					resultItem.Image = mediaPath + productModel.GetDefaultImage()
				}

				if extra, isExtra := reqData["extra"]; isExtra {
					resultItem.Extra = make(map[string]interface{})
					extra := utils.Explode(utils.InterfaceToString(extra), ",")
					for _, value := range extra {
						resultItem.Extra[value] = productModel.Get(value)
					}
				}

				result = append(result, *resultItem)
			}
		}
	}

	return result, nil
}

// WEB REST API function used to obtain visitors count in model collection
func restCountProducts(params *api.StructAPIHandlerParams) (interface{}, error) {

	productCollectionModel, err := product.GetProductCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := productCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	// not allowing to see disabled products if not admin
	if err := api.ValidateAdminRights(params); err != nil {
		productCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	return dbCollection.Count()
}
