package stock_test

import (
	"fmt"
	"testing"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/product"
)

// TestStock validates product inventory model to works properly
func TestStock(t *testing.T) {
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
		return
	}

	initConfig(t)

	productData, err := utils.DecodeJSONToStringKeyMap(`{
		"sku": "test",
		"name": "Test Product",
		"short_description": "something short",
		"description": "something long",
		"default_image": "",
		"price": 1,
		"weight": 1,
		"qty": 10,
		"inventory": [
			{"options": {"color": "black"}, "qty": 1 },
			{"options": {"color": "blue"},  "qty": 5 },
			{"options": {"color": "green"}, "qty": 2 },
			{"options": {"size":  "s"},     "qty": 5 },
			{"options": {"size":  "l"},     "qty": 1 },
			{"options": {"size":  "xl"},    "qty": 5 }
		],
		"options": {
			"color": {
				"order": 1,
				"required": true,
				"options": {
					"black": {"sku": "-black"},
					"blue":  {"sku": "-blue"},
					"green": {"sku": "-green", "price": "+1"}
				}
			},
			"size": {
				"order": 2,
				"required": true,
				"options": {
					"s":  {"sku": "-s",  "price": 1.0 },
					"l":  {"sku": "-l",  "price": 1.5 },
					"xl": {"sku": "-xl", "price": 2.0 }
				}
			}
		}
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	productModel, err := product.GetProductModel()
	if err != nil {
		t.Error(err)
		return
	}

	err = productModel.FromHashMap(productData)
	if err != nil {
		t.Error(err)
		return
	}

	err = productModel.Save()
	if err != nil {
		t.Error(err)
		return
	}
	defer productModel.Delete()

	productID := productModel.GetID()
	registeredStock := product.GetRegisteredStock()

	qty := registeredStock.GetProductQty(productID, map[string]interface{}{"color": "black", "size": "s"})
	if qty != 1 {
		t.Error("The black,s color qty should be 1 and not", qty)
		return
	}

	qty = registeredStock.GetProductQty(productID, map[string]interface{}{"color": "blue", "size": "s"})
	if qty != 5 {
		t.Error("The blue,s color qty should be 5 and not", qty)
		return
	}

	qty = registeredStock.GetProductQty(productID, map[string]interface{}{"color": "green", "size": "xl"})
	if qty != 2 {
		t.Error("The green,xl color qty should be 2 and not", qty)
		return
	}

}

// TestStock validates product inventory model calculations
func TestDecrementingStock(t *testing.T) {
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
		return
	}

	initConfig(t)

	productData, err := utils.DecodeJSONToStringKeyMap(`{
		"sku": "test2",
		"name": "Test Product 2",
		"short_description": "something short",
		"description": "something long",
		"default_image": "",
		"price": 1,
		"weight": 1,
		"qty": 100,
		"options": {
			"color": {
				"order": 1,
				"required": true,
				"options": {
					"black": {"sku": "-black"},
					"blue":  {"sku": "-blue"},
					"green": {"sku": "-green", "price": "+1"}
				}
			},
			"size": {
				"order": 2,
				"required": true,
				"options": {
					"s":  {"sku": "-s",  "price": 1.0},
					"l":  {"sku": "-l",  "price": 1.5},
					"xl": {"sku": "-xl", "price": 2.0 }
				}
			},
			"wrap": {
				"order": 3,
				"required": true,
				"options": {
					"Y":  {"sku": "-s",  "price": 1.0},
					"N":  {"sku": "-l",  "price": 1.5}
				}
			}
		}
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	productModel, err := product.GetProductModel()
	if err != nil {
		t.Error(err)
		return
	}

	err = productModel.FromHashMap(productData)
	if err != nil {
		t.Error(err)
		return
	}

	err = productModel.Save()
	if err != nil {
		t.Error(err)
		return
	}
	defer productModel.Delete()

	productID := productModel.GetID()
	registeredStock := product.GetRegisteredStock()

	// test options
	optionsWrapY := map[string]interface{}{"wrap": "Y"}
	optionsSizeS := map[string]interface{}{"size": "S"}
	optionsColorRedSizeS := map[string]interface{}{"color": "Red", "size": "S"}
	optionsColorGreenSizeXL := map[string]interface{}{"color": "Green", "size": "XL"}

	// setting stock for test options
	registeredStock.SetProductQty(productID, optionsSizeS, 20)
	registeredStock.SetProductQty(productID, optionsColorRedSizeS, 5)
	registeredStock.SetProductQty(productID, optionsColorGreenSizeXL, 20)

	// Test Case 1
	registeredStock.UpdateProductQty(productID, map[string]interface{}{"wrap": "Y"}, -5)

	qty := registeredStock.GetProductQty(productID, map[string]interface{}{})
	qtyWrapY := registeredStock.GetProductQty(productID, optionsWrapY)
	if qty != 95 || qtyWrapY != 95 {
		msg := fmt.Sprintln("Test case 1 error")
		msg += fmt.Sprintln("\t qty:", qty, "(95 expected)")
		msg += fmt.Sprintln("\t qty(Wrap=Y):", qtyWrapY, "(95 expected)")

		t.Error(msg)
		return
	}

	// Test Case 2
	registeredStock.UpdateProductQty(productID, map[string]interface{}{"size": "S"}, -5)

	qty = registeredStock.GetProductQty(productID, map[string]interface{}{})
	qtySizeS := registeredStock.GetProductQty(productID, optionsSizeS)
	if qty != 90 || qtySizeS != 15 {
		msg := fmt.Sprintln("Test case 2 error")
		msg += fmt.Sprintln("\t qty:", qty, "(90 expected)")
		msg += fmt.Sprintln("\t qty(Size=S):", qtySizeS, "(15 expected)")

		t.Error(msg)
		return
	}

	// Test Case 3
	registeredStock.UpdateProductQty(productID, map[string]interface{}{"color": "Red", "size": "S"}, -1)

	// TODO: check is it possible to add more than we have
	qty = registeredStock.GetProductQty(productID, map[string]interface{}{})
	qtySizeS = registeredStock.GetProductQty(productID, optionsSizeS)
	qtyColorRedSizeS := registeredStock.GetProductQty(productID, optionsColorRedSizeS)
	if qty != 89 || qtySizeS != 14 || qtyColorRedSizeS != 4 {
		msg := fmt.Sprintln("Test case 3 error")
		msg += fmt.Sprintln("\t qty:", qty, "(89 expected)")
		msg += fmt.Sprintln("\t qty(Size=S):", qtySizeS, "(14 expected)")
		msg += fmt.Sprintln("\t qty(Color=Red, Size=S):", qtyColorRedSizeS, "(4 expected)")

		t.Error(msg)
		return
	}

	// Test Case 4
	registeredStock.UpdateProductQty(productID, map[string]interface{}{"color": "Green", "size": "XL"}, -5)

	qty = registeredStock.GetProductQty(productID, map[string]interface{}{})
	qtyColorGreenSizeXL := registeredStock.GetProductQty(productID, optionsColorGreenSizeXL)
	if qty != 84 || qtyColorGreenSizeXL != 15 {
		msg := fmt.Sprintln("Test case 4 error")
		msg += fmt.Sprintln("\t qty:", qty, "(84 expected)")
		msg += fmt.Sprintln("\t qty(Color=Green, Size=XL):", qtyColorGreenSizeXL, "(15 expected)")

		t.Error(msg)
		return
	}

	// Test Case 5
	registeredStock.UpdateProductQty(productID, map[string]interface{}{"color": "Red", "size": "S", "wrap": "Y"}, -1)

	qty = registeredStock.GetProductQty(productID, map[string]interface{}{})
	qtyWrapY = registeredStock.GetProductQty(productID, optionsWrapY)
	qtySizeS = registeredStock.GetProductQty(productID, optionsSizeS)
	qtyColorRedSizeS = registeredStock.GetProductQty(productID, optionsColorRedSizeS)
	if qty != 83 && qtyWrapY != 83 && qtySizeS != 13 && qtyColorRedSizeS != 3 {
		msg := fmt.Sprintln("Test case 5 error")
		msg += fmt.Sprintln("\t qty=", qty, "(83 expected)")
		msg += fmt.Sprintln("\t qty(Wrap: Y)=", qtyWrapY, "(83 expected)")
		msg += fmt.Sprintln("\t qty(Size: S)=", qtySizeS, "(13 expected)")
		msg += fmt.Sprintln("\t qty(Color: Red, Size: S)=", qtyColorRedSizeS, "(3 expected)")

		t.Error(msg)
		return
	}

}

// TestCountAfterSetInventory checks if duplicates have not been generated
func TestCountAfterSetInventory(t *testing.T) {
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
		return
	}

	initConfig(t)

	productData, err := utils.DecodeJSONToStringKeyMap(`{
		"sku": "test 3",
		"name": "Test Product 3",
		"short_description": "something short 3",
		"description": "something long 3",
		"default_image": "",
		"price": 3,
		"weight": 3,
		"qty": 30,
		"inventory": [
			{"options": {"color": "black"}, "qty": 1 },
			{"options": {"color": "blue"},  "qty": 5 },
			{"options": {"color": "green"}, "qty": 2 },
			{"options": {"size":  "s"},     "qty": 5 },
			{"options": {"size":  "l"},     "qty": 1 },
			{"options": {"size":  "xl"},    "qty": 5 },
			{"options": {"size":  12},      "qty": 12 },
			{"options": {"color": "black", "size":  "xl"},    "qty": 10 }
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	var testTable = []map[string]interface{}{
		{
			"data": `[
				{"options": {"color": "black"}, "qty": 3 },
				{"options": {"color": "black", "size":  "xl"},    "qty": 13 }
			]`,
			"testCount": 2,
		},
		{
			"data": `[
				{"options": {"size":  "12"},      "qty": 12 }
			]`,
			"testCount": 1,
		},
	}

	for testIdx, testItem := range testTable {
		productModel, err := product.GetProductModel()
		if err != nil {
			t.Error(err)
			return
		}

		err = productModel.FromHashMap(productData)
		if err != nil {
			t.Error(err)
			return
		}

		newInventory, err := utils.DecodeJSONToArray(testItem["data"])
		if err != nil {
			t.Error("Test error:", err)
		}

		err = productModel.Set("inventory", newInventory)
		if err != nil {
			t.Error("Test error:", err)
		}

		var inventory = productModel.Get("inventory")

		if len(utils.InterfaceToArray(inventory)) != testItem["testCount"] {
			t.Error("Test:", testIdx, ". Incorrect number of options:", len(utils.InterfaceToArray(inventory)), ", should be ", testItem["testCount"])
		}
	}
}

// TestDuplicatesForSetInventory checks if new inventory contains no duplicates
func TestDuplicatesForSetInventory(t *testing.T) {
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
		return
	}

	initConfig(t)

	productData, err := utils.DecodeJSONToStringKeyMap(`{
		"sku": "test 4",
		"name": "Test Product 4",
		"short_description": "something short 4",
		"description": "something long 4",
		"default_image": "",
		"price": 4,
		"weight": 4,
		"qty": 40,
		"inventory": [
		]
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	var testData = `[
		{"options": {"color": "black"}, "qty": 3 },
		{"options": {"color": "black"}, "qty": 4 }
	]`

	productModel, err := product.GetProductModel()
	if err != nil {
		t.Error(err)
		return
	}

	err = productModel.FromHashMap(productData)
	if err != nil {
		t.Error(err)
		return
	}

	newInventory, err := utils.DecodeJSONToArray(testData)
	if err != nil {
		t.Error("Test error:", err)
		return
	}

	err = productModel.Set("inventory", newInventory)
	if err == nil {
		t.Error("Should be error, because new data contains duplicates.")
		return
	}
}

// initConfig initializes configuration for tests
func initConfig(t *testing.T) {
	if config := env.GetConfig(); config != nil {
		if config.GetValue("general.stock.enabled") != true {
			err := env.GetConfig().SetValue("general.stock.enabled", true)
			if err != nil {
				t.Error(err)
			}
		}
	}
}
