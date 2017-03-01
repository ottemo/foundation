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
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/product"
)

const (
	ConstPresent = "$present"
	ConstAbsent  = "$absent"
)

type testDataType struct {
	productJson    string
	optionsToApply map[string]interface{}
	testValues     map[string]interface{}

	additionalProductJson string
}

func TestMain(m *testing.M) {
	err := test.StartAppInTestingMode()
	if err != nil {
		fmt.Println("Unable to start app in testing mode:", err)
	}

	os.Exit(m.Run())
}

func TestProductApplyOptions(t *testing.T) {
	fmt.Println("TestProductApplyOptions")

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
			"field_option": "` + ConstPresent + `",
			"another_option": "` + ConstAbsent + `",
			"color": {
				"options": {
					"red": "` + ConstPresent + `",
					"black": "` + ConstAbsent + `",
					"blue": "` + ConstAbsent + `"
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

func TestProductApplyOptionsQty(t *testing.T) {
	fmt.Println("TestProductApplyOptionsQty")

	if config := env.GetConfig(); config != nil {
		if config.GetValue("general.stock.enabled") != true {
			err := env.GetConfig().SetValue("general.stock.enabled", true)
			if err != nil {
				t.Error(err)
			}
		}
	}

	var product = createProductFromJson(t, `{
		"_id": "123456789012345678904444",
		"sku": "test",
		"name": "Test",
		"price": 1,
		"options": {
			"color": {
				"controls_inventory": true, "key": "color", "label": "color",
				"options": {
					"red_2": {"key": "red_2", "label": "red-2", "order": 1, "price": ""},
					"red_3": {"key": "red_3", "label": "red-3", "order": 2}
				},
				"order": 1, "required": true, "type": "select"
			}
		},
		"inventory": [
		    {
			"options": { },
			"qty": 5
		    },
		    {
			"options": {
			    "color": "red_2"
			},
			"qty": 2
		    },
		    {
			"options": {
			    "color": "red_3"
			},
			"qty": 3
		    }
		],
		"qty": 5
	}`)

	appliedOptions := map[string]interface{}{
		"color": "red_2",
	}

	checkJson := `{
		"options": {
			"color": {
				"options": {
					"red_2": "` + ConstPresent + `"
				}
			}
		},
		"qty": "2"
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
	fmt.Println("TestConfigurableProductApplyOption")

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
						"`+product.ConstOptionProductIDs+`": ["`+simpleProduct.GetID()+`"]
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
					"red": "` + ConstPresent + `",
					"black": "` + ConstAbsent + `",
					"blue": "` + ConstAbsent + `"
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
	fmt.Println("TestConfigurableProductApplyOptions")

	var simpleProduct1 = createProductFromJson(t, `{
		"sku": "test-simple-1",
		"enabled": "true",
		"name": "Test Simple Product 1",
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

	var simpleProduct2 = createProductFromJson(t, `{
		"sku": "test-simple-2",
		"enabled": "true",
		"name": "Test Simple Product 2",
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

	var simpleProduct3 = createProductFromJson(t, `{
		"sku": "test-simple-3",
		"enabled": "true",
		"name": "Test Simple Product 3",
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

	var simpleProduct4 = createProductFromJson(t, `{
		"sku": "test-simple-4",
		"enabled": "true",
		"name": "Test Simple Product 4",
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
					"blue":  {"order": "1", "key": "blue",  "label": "Blue",  "price": 2.0, "sku": "-blue",
						"`+product.ConstOptionProductIDs+`": ["`+simpleProduct3.GetID()+`", "`+simpleProduct4.GetID()+`"]
					},
					"red":   {
						"order": "2", "key": "red",   "label": "Red",   "price": 100, "sku": "-red",
						"`+product.ConstOptionProductIDs+`": ["`+simpleProduct1.GetID()+`", "`+simpleProduct2.GetID()+`"]
					}
				}
			},
			"size" : {
				"code": "size", "controls_inventory": true, "key": "size", "label": "Size",
				"order": 2, "required": true, "type": "select",
				"options" : {
					"xl": {"order": "1", "key": "xl", "label": "xxl", "price": 1.3, "sku": "-xl",
						"`+product.ConstOptionProductIDs+`": ["`+simpleProduct1.GetID()+`", "`+simpleProduct3.GetID()+`"]
					},
					"xxl":   {
						"order": "2", "key": "xxl",   "label": "xxl",   "price": 100, "sku": "-xxl",
						"`+product.ConstOptionProductIDs+`": ["`+simpleProduct2.GetID()+`", "`+simpleProduct4.GetID()+`"]
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
		"sku": "` + simpleProduct2.GetSku() + `",
		"price": 1.0,
		"options": {
			"field_option": {
				"value": "` + utils.InterfaceToString(appliedOptions["field_option"]) + `"
			},
			"color": {
				"options": {
					"red": "` + ConstPresent + `",
					"black": "` + ConstAbsent + `",
					"blue": "` + ConstAbsent + `"
				}
			},
			"size": {
				"options": {
					"xl": "` + ConstAbsent + `",
					"xxl": "` + ConstPresent + `"
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

	deleteProduct(t, simpleProduct1)
	deleteProduct(t, simpleProduct2)
	deleteProduct(t, simpleProduct3)
	deleteProduct(t, simpleProduct4)
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
			if checkValueStr == ConstPresent {
				if _, present := valueMap[key]; !present {
					t.Error("Key [" + key + "] not present.")
				}
			} else if checkValueStr == ConstAbsent {
				if _, present := valueMap[key]; present {
					t.Error("Key [" + key + "] present.")
				}
			} else {
				assert.Equal(t, checkValueStr, valueStr, key)
			}
		}
	}
}
