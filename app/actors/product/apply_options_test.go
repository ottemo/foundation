package product_test

// This package provides additional product package tests
// To run it use command line
//
// $ go test -tags sqlite github.com/ottemo/foundation/app/actors/product/
//
// or, if fmt.Println output required, use it with "-v" flag
//
// $ go test -v -tags sqlite github.com/ottemo/foundation/app/actors/product/

import (
	"testing"

	"fmt"

	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/product"
)

type testDataType struct {
	productJson    string
	optionsToApply map[string]interface{}
	testValues     map[string]interface{}

	additionalProductJson string
}

func TestConfigurableProductApplyOptions(t *testing.T) {

	start(t)

	var simpleProduct = createProductFromJson(t, `{
		"sku": "test-simple",
		"enabled": "true",
		"name": "Test Simple Product",
		"short_description": "something short",
		"description": "something long",
		"default_image": "",
		"price": 1.0,
		"weight": 0.4
	}`)

	var configurable = populateProductModel(t, `{
		"id": "123456789012345678901234",
		"sku": "test",
		"name": "Test Product",
		"short_description": "something short",
		"description": "something long",
		"default_image": "",
		"price": 1.1,
		"weight": 0.5,
		"test": "ok",
		"options" : {
			"field_option":{
				"code": "field_option", "controls_inventory": false, "key": "field_option",
				"label": "FieldOption", "order": 2, "price": "13", "required": false,
				"sku": "-fo", "type": "field"},
			"color" : {
				"code": "color", "controls_inventory": true, "key": "color", "label": "Color",
				"order": 1, "required": true, "type": "select",
				"options" : {
					"black": {"order": "3", "key": "black", "label": "Black", "price": 1.3, "sku": "-black"},
					"blue":  {"order": "1", "key": "blue",  "label": "Blue",  "price": 2.0, "sku": "-blue"},
					"red":   {
						"order": "2", "key": "red",   "label": "Red",   "price": 100, "sku": "-red",
						"` + product.ConstOptionSimpleIDsName + `": ["` + simpleProduct.GetID() + `"]
					}
				}
			}
		}
	}`)

	configurable = applyOptions(t, configurable, map[string]interface{}{
		"color":        "red",
		"field_option": "field_option value",
	})

	deleteProduct(t, simpleProduct)
}

func start(t *testing.T) {
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}
}

func createProductFromJson(t *testing.T, json string) product.InterfaceProduct {
	productData, err := utils.DecodeJSONToStringKeyMap(json)
	if err != nil {
		t.Error(err)
	}

	productModel, err := product.GetProductModel()
	if err != nil || productModel == nil {
		t.Error(err)
	}

	err = productModel.FromHashMap(productData)
	if err != nil {
		t.Error(err)
	}

	err = productModel.Save()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(utils.InterfaceToString("\n= Additional product original ID: " + productModel.GetID()))

	return productModel
}

func deleteProduct(t *testing.T, productModel product.InterfaceProduct) {
	err := productModel.Delete()
	if err != nil {
		t.Error(err)
	}
}

func populateProductModel(t *testing.T, json string) product.InterfaceProduct {
	productData, err := utils.DecodeJSONToStringKeyMap([]byte(json))
	if err != nil {
		t.Error(err)
	}

	productModel, err := product.GetProductModel()
	if err != nil || productModel == nil {
		t.Error(err)
	}

	err = productModel.FromHashMap(productData)
	if err != nil {
		t.Error(err)
	}

	return productModel
}

func applyOptions(
	t *testing.T,
	productModel product.InterfaceProduct,
	options map[string]interface{}) product.InterfaceProduct {

	err := productModel.ApplyOptions(options)
	if err != nil {
		t.Error("Error applying options")
	}

	return productModel
}

