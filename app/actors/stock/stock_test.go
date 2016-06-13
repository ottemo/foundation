package stock_test

import (
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/tests"
	"github.com/ottemo/foundation/utils"

	"fmt"
	"testing"
)

// TestStock validates product inventory model to works properly
func TestStock(t *testing.T) {
	err := tests.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
		return
	}

	if config := env.GetConfig(); config != nil {
		if config.GetValue("general.stock.enabled") != true {
			err := env.GetConfig().SetValue("general.stock.enabled", true)
			if err != nil {
				t.Error(err)
				return
			}
		}
	}

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
			{"options": {"size":  "s"},     "qty": 5 },
			{"options": {"size":  "l"},     "qty": 1 }
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

	productTestModel, _ := product.LoadProductByID(productID)
	productTestModel.ApplyOptions(map[string]interface{}{"color": "black", "size": "s"})
	if qty := productTestModel.GetQty(); qty != 1 {
		t.Error("The black,s color qty should be 1 and not", qty)
		return
	}

	// TODO: find out why second ApplyOptions call to existing model have no effect
	productTestModel, _ = product.LoadProductByID(productID)
	productTestModel.ApplyOptions(map[string]interface{}{"color": "blue", "size": "s"})
	if qty := productTestModel.GetQty(); qty != 5 {
		t.Error("The blue,s color qty should be 5 and not", qty)
		return
	}

	productTestModel, _ = product.LoadProductByID(productID)
	productTestModel.ApplyOptions(map[string]interface{}{"color": "green", "size": "xl"})
	if qty := productTestModel.GetQty(); qty != 10 {
		t.Error("The green,xl color qty should be 10 and not", qty)
		return
	}

}

// TestStock validates product inventory model calculations
func TestDecrementingStock(t *testing.T) {
	err := tests.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
		return
	}

	if config := env.GetConfig(); config != nil {
		if config.GetValue("general.stock.enabled") != true {
			err := env.GetConfig().SetValue("general.stock.enabled", true)
			if err != nil {
				t.Error(err)
				return
			}
		}
	}

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
	stock := product.GetRegisteredStock()

	// test options
	optionsWrapY := map[string]interface{}{"wrap": "Y"}
	optionsSizeS := map[string]interface{}{"size": "S"}
	optionsColorRedSizeS := map[string]interface{}{"color": "Red", "size": "S"}
	optionsColorGreenSizeXL := map[string]interface{}{"color": "Green", "size": "XL"}

	// setting stock for test options
	stock.SetProductQty(productID, optionsSizeS, 20)
	stock.SetProductQty(productID, optionsColorRedSizeS, 5)
	stock.SetProductQty(productID, optionsColorGreenSizeXL, 20)

	// Test Case 1
	stock.UpdateProductQty(productID, map[string]interface{}{"wrap": "Y"}, -5)

	productTestModel, _ := product.LoadProductByID(productID)
	qty := productTestModel.GetQty()
	qtyWrapY := stock.GetProductQty(productID, optionsWrapY)
	if qty != 95 || qtyWrapY != 95 {
		msg := fmt.Sprintln("Test case 1 error")
		msg += fmt.Sprintln("\t qty:", qty, "(95 expected)")
		msg += fmt.Sprintln("\t qty(Wrap=Y):", qtyWrapY, "(95 expected)")

		t.Error(msg)
		return
	}

	// Test Case 2
	stock.UpdateProductQty(productID, map[string]interface{}{"size": "S"}, -5)

	productTestModel, _ = product.LoadProductByID(productID)
	qty = productTestModel.GetQty()
	qtySizeS := stock.GetProductQty(productID, optionsSizeS)
	if qty != 90 || qtySizeS != 15 {
		msg := fmt.Sprintln("Test case 2 error")
		msg += fmt.Sprintln("\t qty:", qty, "(90 expected)")
		msg += fmt.Sprintln("\t qty(Size=S):", qtySizeS, "(15 expected)")

		t.Error(msg)
		return
	}

	// Test Case 3
	stock.UpdateProductQty(productID, map[string]interface{}{"color": "Red", "size": "S"}, -1)

	// TODO: check is it possible to add more than we have
	productTestModel, _ = product.LoadProductByID(productID)
	qty = productTestModel.GetQty()
	qtySizeS = stock.GetProductQty(productID, optionsSizeS)
	qtyColorRedSizeS := stock.GetProductQty(productID, optionsColorRedSizeS)
	if qty != 89 || qtySizeS != 14 || qtyColorRedSizeS != 4 {
		msg := fmt.Sprintln("Test case 3 error")
		msg += fmt.Sprintln("\t qty:", qty, "(89 expected)")
		msg += fmt.Sprintln("\t qty(Size=S):", qtySizeS, "(14 expected)")
		msg += fmt.Sprintln("\t qty(Color=Red, Size=S):", qtyColorRedSizeS, "(4 expected)")

		t.Error(msg)
		return
	}

	// Test Case 4
	stock.UpdateProductQty(productID, map[string]interface{}{"color": "Green", "size": "XL"}, -5)

	productTestModel, _ = product.LoadProductByID(productID)
	qty = productTestModel.GetQty()
	qtyColorGreenSizeXL := stock.GetProductQty(productID, optionsColorGreenSizeXL)
	if qty != 84 || qtyColorGreenSizeXL != 15 {
		msg := fmt.Sprintln("Test case 4 error")
		msg += fmt.Sprintln("\t qty:", qty, "(84 expected)")
		msg += fmt.Sprintln("\t qty(Color=Green, Size=XL):", qtyColorGreenSizeXL, "(15 expected)")

		t.Error(msg)
		return
	}

	// Test Case 5
	stock.UpdateProductQty(productID, map[string]interface{}{"color": "Red", "size": "S", "wrap": "Y"}, -1)

	productTestModel, _ = product.LoadProductByID(productID)
	qty = productTestModel.GetQty()
	qtyWrapY = stock.GetProductQty(productID, optionsWrapY)
	qtySizeS = stock.GetProductQty(productID, optionsSizeS)
	qtyColorRedSizeS = stock.GetProductQty(productID, optionsColorRedSizeS)
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
