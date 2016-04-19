package stock_test

import (
	"github.com/ottemo/foundation/tests"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"

	"testing"
)

func TestStock(t *testing.T) {
	err := tests.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
		return;
	}

	if config := env.GetConfig(); config != nil {
		if config.GetValue("general.stock.enabled") != true {
			err := env.GetConfig().SetValue("general.stock.enabled", true);
			if err != nil {
				t.Error(err)
				return;
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
		return;
	}

	productModel, err := product.GetProductModel()
	if err != nil {
		t.Error(err)
		return;
	}

	err = productModel.FromHashMap(productData)
	if err != nil {
		t.Error(err)
		return;
	}

	err = productModel.Save()
	if err != nil {
		t.Error(err)
		return;
	}
	// defer productModel.Delete()

	productID := productModel.GetID()

	productTestModel, _ := product.LoadProductByID(productID)
	productTestModel.ApplyOptions(map[string]interface{} {"color": "black", "size": "s"})
	if qty := productTestModel.GetQty(); qty != 1 {
		t.Error("The black,s color qty should be 1 and not", qty)
		return;
	}

	// TODO: find out why second ApplyOptions call to existing model have no effect
	productTestModel, _ = product.LoadProductByID(productID)
	productTestModel.ApplyOptions(map[string]interface{} {"color": "blue", "size": "s"})
	if qty := productTestModel.GetQty(); qty != 5 {
		t.Error("The blue,s color qty should be 5 and not", qty)
		return;
	}

	productTestModel, _ = product.LoadProductByID(productID)
	productTestModel.ApplyOptions(map[string]interface{} {"color": "green", "size": "xl"})
	if qty := productTestModel.GetQty(); qty != 10 {
		t.Error("The green,xl color qty should be 10 and not", qty)
		return;
	}

}