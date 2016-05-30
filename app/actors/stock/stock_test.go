package stock_test

import (
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/tests"
	"github.com/ottemo/foundation/utils"

	"fmt"
	"testing"
)

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
		"options": {
			"color": {
				"order": 1,
				"required": true,
				"options": {
					"black": {"sku": "-black", "qty": 1},
					"blue":  {"sku": "-blue",  "qty": 5},
					"green": {"sku": "-green", "price": "+1"}
				}
			},
			"size": {
				"order": 2,
				"required": true,
				"options": {
					"s":  {"sku": "-s",  "price": 1.0, "qty": 5},
					"l":  {"sku": "-l",  "price": 1.5, "qty": 1},
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
	// defer productModel.Delete()

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
		"sku": "test",
		"name": "Test Product",
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
	// defer productModel.Delete()

	productID := productModel.GetID()
	stock := product.GetRegisteredStock()

	// define options
	optionsWrapY := map[string]interface{}{"wrap": "Y"}
	optionsSizeS := map[string]interface{}{"size": "S"}
	optionsColorRedSizeS := map[string]interface{}{"color": "Red", "size": "S"}
	optionsColorGreenSizeXL := map[string]interface{}{"color": "Green", "size": "XL"}

	// set stock
	stock.SetProductQty(productID, optionsSizeS, 20)
	stock.SetProductQty(productID, optionsColorRedSizeS, 5)
	stock.SetProductQty(productID, optionsColorGreenSizeXL, 20)

	// if options not managed by stock
	// example: "wrap": "Y"
	stock.UpdateProductQty(productID, map[string]interface{}{"wrap": "Y"}, -5)
	// gets qty for specified options if it doesn't exist returns minimal from matched options
	qtyWrapY := stock.GetProductQty(productID, optionsWrapY)
	productTestModel, _ := product.LoadProductByID(productID)
	if qty := productTestModel.GetQty(); qty == 95 && qtyWrapY == 95 {
		fmt.Println("case 1 -5 \n\t baseQty: ", qty, "\n\t qty Wrap - Y: ", qtyWrapY)
	} else {
		t.Error("case 1 error")
		return
	}

	// if one option managed by stock
	// example: "size": "S"
	stock.UpdateProductQty(productID, map[string]interface{}{"size": "S"}, -5)
	qtySizeS := stock.GetProductQty(productID, optionsSizeS)
	productTestModel, _ = product.LoadProductByID(productID)
	if qty := productTestModel.GetQty(); qty == 90 && qtySizeS == 15 {
		fmt.Println("\ncase 2 -5 \n\t baseQty: ", qty, "\n\t qty Size - S: ", qtySizeS)
	} else {
		t.Error("case 2 error")
		return
	}

	// if multiple options managed by stock
	// example: "color": "Red", "size": "S"}
	stock.UpdateProductQty(productID, map[string]interface{}{"color": "Red", "size": "S"}, -1)
	// TODO: check is it possible to add more than we have
	qtySizeS = stock.GetProductQty(productID, optionsSizeS)
	qtyColorRedSizeS := stock.GetProductQty(productID, optionsColorRedSizeS)
	productTestModel, _ = product.LoadProductByID(productID)
	if qty := productTestModel.GetQty(); qty == 89 && qtySizeS == 14 && qtyColorRedSizeS == 4 {
		fmt.Println("\ncase 3 -1 \n\t baseQty: ", qty, "\n\t qty Size - S: ", qtySizeS, "\n\t qty Color - Red, Size - S: ", qtyColorRedSizeS)
	} else {
		t.Error("case 3 error")
		return
	}

	// example: "color": "Green", "size": "XL"
	stock.UpdateProductQty(productID, map[string]interface{}{"color": "Green", "size": "XL"}, -5)
	qtyColorGreenSizeXL := stock.GetProductQty(productID, optionsColorGreenSizeXL)
	productTestModel, _ = product.LoadProductByID(productID)
	if qty := productTestModel.GetQty(); qty == 84 && qtyColorGreenSizeXL == 15 {
		fmt.Println("\ncase 4 -5 \n\t baseQty: ", qty, "\n\t qty Color - Green, Size - XL: ", qtyColorGreenSizeXL)
	} else {
		t.Error("case 4 error")
		return
	}

	// if exist multiple options managed by stock and one not managed by stock
	// "color": "Red", "size": "S", "wrap": "Y"
	stock.UpdateProductQty(productID, map[string]interface{}{"color": "Red", "size": "S", "wrap": "Y"}, -1)
	qtyWrapY = stock.GetProductQty(productID, optionsWrapY)
	qtySizeS = stock.GetProductQty(productID, optionsSizeS)
	qtyColorRedSizeS = stock.GetProductQty(productID, optionsColorRedSizeS)
	productTestModel, _ = product.LoadProductByID(productID)
	if qty := productTestModel.GetQty(); qty == 83 && qtyWrapY == 83 && qtySizeS == 13 && qtyColorRedSizeS == 3 {
		fmt.Println("\ncase 5 -1 \n\t baseQty: ", qty, "\n\t qty Wrap - Y: ", qtyWrapY, "\n\t qty Size - S: ", qtySizeS, "\n\t qty Color - Red, Size - S: ", qtyColorRedSizeS)
	} else {
		t.Error("case 5 error")
		return
	}

}
