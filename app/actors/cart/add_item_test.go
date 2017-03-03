package cart_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/product"
)

func TestMain(m *testing.M) {
	err := test.StartAppInTestingMode()
	if err != nil {
		fmt.Println("Unable to start app in testing mode:", err)
	}

	os.Exit(m.Run())
}

func TestCartAddItem(t *testing.T) {

	currentVisitor, err := test.GetRandomVisitor()
	if err != nil {
		t.Error(err)
	}

	currentCheckout, err := test.GetNewCheckout(currentVisitor)
	if err != nil {
		t.Error(err)
	}

	currentCart := currentCheckout.GetCart()
	if err != nil {
		t.Error(err)
	}

	var testSku = "sku"
	var testSkuModifier = "-mod"

	var productModel = createProductFromJson(t, `{
		"_id": "123456789012345678904444",
		"sku": "`+testSku+`",
		"name": "Test",
		"price": 1,
		"options": {
			"color": {
				"controls_inventory": true, "key": "color", "label": "color",
				"options": {
					"red_2": {
						"key": "red_2", "label": "red-2", "order": 1,
						"sku": "`+testSkuModifier+`",
						"price": "+10"
					},
					"red_3": {"key": "red_3", "label": "red-3", "order": 2, "sku": "-3"}
				},
				"order": 1, "required": true, "type": "select"
			}
		},
		"inventory": [
		    {"options": { }, "qty": 5},
		    {"options": {"color": "red_2"}, "qty": 2},
		    {"options": {"color": "red_3"}, "qty": 3}
		],
		"qty": 5,
		"enabled": true,
		"visible": true
	}`)

	appliedOptions := map[string]interface{}{
		"color": "red_2",
	}

	cartItem, err := currentCart.AddItem(productModel.GetID(), 1, appliedOptions)
	if err != nil {
		t.Error(err)
	}

	cartItemProduct := cartItem.GetProduct()

	// test
	var expectedSku = testSku + testSkuModifier
	var gotSku = cartItemProduct.GetSku()
	if gotSku != expectedSku {
		t.Errorf("Incorrect Sku. Expected: '%v'. Got: '%v'.", expectedSku, gotSku)
	}

	var expectedPrice = 11.0
	var gotPrice = cartItemProduct.GetPrice()
	if gotPrice != expectedPrice {
		t.Errorf("Incorrect Price. Expected: '%v'. Got: '%v'.", expectedPrice, gotPrice)
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
