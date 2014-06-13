package default_product

import (
	"errors"
	"net/http"

	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/product"
)

func jsonError(err error) map[string]interface{} {
	return map[string]interface{} { "error": err.Error() }
}

// http://127.0.0.1:9000/AddProductAttribute
// http://127.0.0.1:9000/AddProductAttribute?name=x&type=int&group=Others
func AddProductAttributeRestAPI(req *http.Request) map[string]interface{} {
	queryParams := req.URL.Query()

	model, err := models.GetModel("Product")
	if err != nil { return jsonError(err) }

	attribute := models.T_AttributeInfo {
		Model:      "product",
		Collection: "product",
		Attribute:  "test",
		Type:       "text",
		Label:      "Test Attribute",
		Group:      "General",
		Editors:    "text",
		Options:    "",
		Default:    "",
	}


	for param, value := range queryParams {
		switch param {
		case "type":
			attribute.Type = value[0]
		case "attribute":
			attribute.Attribute = value[0]
		case "label":
			attribute.Label = value[0]
		case "group":
			attribute.Group = value[0]
		case "editors":
			attribute.Editors = value[0]
		case "options":
			attribute.Options = value[0]
		case "default":
			attribute.Default = value[0]
		}
	}


	if prod, ok := model.(models.I_CustomAttributes); ok {
		if err := prod.AddNewAttribute(attribute); err != nil {
			return jsonError( errors.New("Product new attribute error: " + err.Error()) )
		}
	} else {
		return jsonError( errors.New("product model is not I_CustomAttributes") )
	}


	return map[string]interface{} {"ok": true, "attribute": attribute}
}


// http://127.0.0.1:9000/LoadProduct?id=5
func LoadProductRestAPI(req *http.Request) map[string]interface{} {
	queryParams := req.URL.Query()

	productId := queryParams.Get("id")
	if productId == "" {
		return jsonError( errors.New("product 'id' was not specified") )
	}

	if model, err := models.GetModel("Product"); err == nil {
		if model, ok := model.(product.I_Product); ok {

			err = model.Load( productId )
			if err != nil { return jsonError(err) }

			return model.ToHashMap()
		}
	}

	return jsonError( errors.New("Something went wrong...") )
}


// http://127.0.0.1:9000/CreateProduct/xx-25/some
func CreateProductRestAPI(req *http.Request) map[string]interface{} {
	queryParams := req.URL.Query()

	if queryParams.Get("sku") == "" || queryParams.Get("name") == "" {
		return jsonError( errors.New("product 'name' and/or 'sku' was not specified") )
	}

	if model, err := models.GetModel("Product"); err == nil {
		if model, ok := model.(product.I_Product); ok {

			for attribute, value := range queryParams {
				err := model.Set(attribute, value[0])
				if err != nil { return jsonError(err) }
			}

			err := model.Save()
			if err != nil { return jsonError(err) }

			return model.ToHashMap()
		}
	}

	return jsonError( errors.New("Something went wrong...") )
}
