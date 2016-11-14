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

	"github.com/stretchr/testify/assert"

	"fmt"

	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/product"
)

const (
	PRESENT = "$present"
	ABSENT  = "$absent"
)

type testDataType struct {
	productJson    string
	optionsToApply map[string]interface{}
	testValues     map[string]interface{}

	additionalProductJson string
}

func TestProductApplyOptions(t *testing.T) {

	start(t)

	var product = populateProductModel(t, `{
		"_id": "123456789012345678901234",
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
				"label": "FieldOption", "order": 2, "price": "+13", "required": false,
				"sku": "-fo", "type": "field"
			},
			"another_option":{
				"code": "another_option", "controls_inventory": false, "key": "another_option",
				"label": "AnotherOption", "order": 3, "price": "14", "required": false,
				"sku": "-ao", "type": "field"
			},
			"color" : {
				"code": "color", "controls_inventory": true, "key": "color", "label": "Color",
				"order": 1, "required": true, "type": "select",
				"options" : {
					"black": {"order": "3", "key": "black", "label": "Black", "price": 1.3, "sku": "-black"},
					"blue":  {"order": "1", "key": "blue",  "label": "Blue",  "price": 2.0, "sku": "-blue"},
					"red":   {
						"order": "2", "key": "red",   "label": "Red",   "price": 100, "sku": "-red"
					}
				}
			}
		}
	}`)

	appliedOptions := map[string]interface{}{
		"color":        "red",
		"field_option": "field_option value",
	}

	checkJson := `{
		"sku": "test-red-fo",
		"price": 113,
		"options": {
			"field_option": "` + PRESENT + `",
			"another_option": "` + ABSENT + `",
			"color": {
				"options": {
					"red": "` + PRESENT + `",
					"black": "` + ABSENT + `",
					"blue": "` + ABSENT + `"
				}
			}
		}
	}`

	product = applyOptions(t, product, appliedOptions)

	check, err := utils.DecodeJSONToInterface(checkJson)
	if err != nil {
		fmt.Println("checkJson: " + checkJson)
		t.Error(err)
	}

	checkResults(t, product.ToHashMap(), check.(map[string]interface{}))
}

func TestConfigurableProductApplyOption(t *testing.T) {

	start(t)

	var simpleProduct = createProductFromJson(t, `{
		"sku": "test-simple",
		"enabled": "true",
		"name": "Test Simple Product",
		"short_description": "something short",
		"description": "something long",
		"default_image": "",
		"price": 1.0,
		"weight": 0.4,
		"inventory": [
			{
        			"options": { },
        			"qty": 13
			}
		]
	}`)

	var configurable = populateProductModel(t, `{
		"_id": "123456789012345678901234",
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
						"`+product.ConstOptionSimpleIDsName+`": ["`+simpleProduct.GetID()+`"]
					}
				}
			}
		}
	}`)

	appliedOptions := map[string]interface{}{
		"color":        "red",
		"field_option": "field_option value",
	}

	checkJson := `{
		"_id": "` + configurable.GetID() + `",
		"sku": "` + simpleProduct.GetSku() + `",
		"price": 1.0,
		"options": {
			"field_option": {
				"value": "` + utils.InterfaceToString(appliedOptions["field_option"]) + `"
			},
			"color": {
				"options": {
					"red": "` + PRESENT + `",
					"black": "` + ABSENT + `",
					"blue": "` + ABSENT + `"
				}
			}
		}
	}`

	configurable = applyOptions(t, configurable, appliedOptions)

	check, err := utils.DecodeJSONToInterface(checkJson)
	if err != nil {
		fmt.Println("checkJson: " + checkJson)
		t.Error(err)
	}

	checkResults(t, configurable.ToHashMap(), check.(map[string]interface{}))

	deleteProduct(t, simpleProduct)
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
		"weight": 0.4,
		"inventory": [
			{
        			"options": { },
        			"qty": 13
			}
		]
	}`)

	var configurable = populateProductModel(t, `{
		"_id": "123456789012345678901234",
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
						"`+product.ConstOptionSimpleIDsName+`": ["`+simpleProduct.GetID()+`"]
					}
				}
			},
			"size" : {
				"code": "size", "controls_inventory": true, "key": "size", "label": "Size",
				"order": 2, "required": true, "type": "select",
				"options" : {
					"xl": {"order": "1", "key": "xl", "label": "xxl", "price": 1.3, "sku": "-xl"},
					"xxl":   {
						"order": "2", "key": "xxl",   "label": "xxl",   "price": 100, "sku": "-xxl",
						"`+product.ConstOptionSimpleIDsName+`": ["`+simpleProduct.GetID()+`"]
					}
				}
			}
		}
	}`)

	appliedOptions := map[string]interface{}{
		"color":        "red",
		"size":         "xxl",
		"field_option": "field_option value",
	}

	checkJson := `{
		"_id": "` + configurable.GetID() + `",
		"sku": "` + simpleProduct.GetSku() + `",
		"price": 1.0,
		"options": {
			"field_option": {
				"value": "` + utils.InterfaceToString(appliedOptions["field_option"]) + `"
			},
			"color": {
				"options": {
					"red": "` + PRESENT + `",
					"black": "` + ABSENT + `",
					"blue": "` + ABSENT + `"
				}
			},
			"size": {
				"options": {
					"xl": "` + ABSENT + `",
					"xxl": "` + PRESENT + `"
				}
			}
		}
	}`

	configurable = applyOptions(t, configurable, appliedOptions)

	check, err := utils.DecodeJSONToInterface(checkJson)
	if err != nil {
		fmt.Println("checkJson: " + checkJson)
		t.Error(err)
	}

	checkResults(t, configurable.ToHashMap(), check.(map[string]interface{}))

	deleteProduct(t, simpleProduct)
}

func TestFailedConfigurableProductApplyOptions(t *testing.T) {

	start(t)

	var simpleProduct = createProductFromJson(t, `{
		"sku": "test-simple",
		"enabled": "true",
		"name": "Test Simple Product",
		"short_description": "something short",
		"description": "something long",
		"default_image": "",
		"price": 1.0,
		"weight": 0.4,
		"inventory": [
			{
        			"options": { },
        			"qty": 13
			}
		]
	}`)

	var configurable = populateProductModel(t, `{
		"_id": "123456789012345678901234",
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
						"`+product.ConstOptionSimpleIDsName+`": ["`+simpleProduct.GetID()+`"]
					}
				}
			},
			"size" : {
				"code": "size", "controls_inventory": true, "key": "size", "label": "Size",
				"order": 2, "required": true, "type": "select",
				"options" : {
					"xl": {"order": "1", "key": "xl", "label": "xxl", "price": 1.3, "sku": "-xl"},
					"xxl":   {
						"order": "2", "key": "xxl",   "label": "xxl",   "price": 100, "sku": "-xxl",
						"`+product.ConstOptionSimpleIDsName+`": ["`+simpleProduct.GetID()+`"]
					}
				}
			}
		}
	}`)

	appliedOptions := map[string]interface{}{
		"color":        "red",
		"size":         "xxl",
		"field_option": "field_option value",
	}

	checkJson := `{
		"_id": "` + configurable.GetID() + `",
		"sku": "` + simpleProduct.GetSku() + `",
		"price": 1.0,
		"options": {
			"field_option": {
				"value": "` + utils.InterfaceToString(appliedOptions["field_option"]) + `"
			},
			"color": {
				"options": {
					"black": "` + PRESENT + `"
				}
			},
			"size": {
				"options": {
					"xl": "` + PRESENT + `"
				}
			}
		}
	}`

	configurable = applyOptions(t, configurable, appliedOptions)

	check, err := utils.DecodeJSONToInterface(checkJson)
	if err != nil {
		fmt.Println("checkJson: " + checkJson)
		t.Error(err)
	}
	checkFailedResults(t, configurable.ToHashMap(), check.(map[string]interface{}))

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
		fmt.Println("json: " + json)
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

func checkResults(
	t *testing.T,
	valueMap map[string]interface{},
	checkMap map[string]interface{}) {
	for key, checkValue := range checkMap {
		value := valueMap[key]
		switch typedCheckValue := checkValue.(type) {
		case map[string]interface{}:
			checkResults(t, value.(map[string]interface{}), typedCheckValue)
		default:
			valueStr := utils.InterfaceToString(value)
			checkValueStr := utils.InterfaceToString(checkValue)
			if checkValueStr == PRESENT {
				if _, present := valueMap[key]; !present {
					t.Error("Key [" + key + "] not present.")
				}
			} else if checkValueStr == ABSENT {
				if _, present := valueMap[key]; present {
					t.Error("Key [" + key + "] present.")
				}
			} else {
				assert.Equal(t, valueStr, checkValueStr)
			}
		}
	}
}

func checkFailedResults(
	t *testing.T,
	valueMap map[string]interface{},
	checkMap map[string]interface{}) {
        // @todo need refactoring if one option not failed function not work
	for key, checkValue := range checkMap {
		value := valueMap[key]
		switch typedCheckValue := checkValue.(type) {
		case map[string]interface{}:
			checkFailedResults(t, value.(map[string]interface{}), typedCheckValue)
		default:
			valueStr := utils.InterfaceToString(value)
			checkValueStr := utils.InterfaceToString(checkValue)
			if checkValueStr == PRESENT {
				if _, present := valueMap[key]; present {
					t.Error("Key [" + key + "] present.")
				}
			} else if checkValueStr == ABSENT {
				if _, present := valueMap[key]; !present {
					t.Error("Key [" + key + "] not present.")
				}
			} else {
				assert.Equal(t, valueStr, checkValueStr)
			}
		}
	}
}
