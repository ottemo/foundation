package cart_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/api/rest"
	cartActor "github.com/ottemo/foundation/app/actors/cart"
	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/product"
)

const (
	constTestSku         = "sku"
	constTestSkuModifier = "-mod"
	constExpectedPrice   = 11.0
)

func TestMain(m *testing.M) {
	err := test.StartAppInTestingMode()
	if err != nil {
		fmt.Println("Unable to start app in testing mode:", err)
	}

	os.Exit(m.Run())
}

func TestCartItem(t *testing.T) {
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

	var productModel = createProductFromJson(t, `{
		"_id": "123456789012345678904444",
		"sku": "`+constTestSku+`",
		"name": "Test",
		"price": 1,
		"options": {
			"color": {
				"controls_inventory": true, "key": "color", "label": "color",
				"options": {
					"red_2": {
						"key": "red_2", "label": "red-2", "order": 1,
						"sku": "`+constTestSkuModifier+`",
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

	// test
	testCartAddItem(t, cartItem)

	// save cart
	applicationContext := new(rest.DefaultRestApplicationContext)
	applicationContext.SetSession(currentCheckout.GetSession())

	currentCheckout.GetSession().Set(cart.ConstSessionKeyCurrentCart, currentCart.GetID())

	err = currentCart.Save()
	if err != nil {
		t.Error(err)
	}

	// test
	testAPICartInfo(t, applicationContext)
}

func testCartAddItem(t *testing.T, cartItem cart.InterfaceCartItem) {
	cartItemProduct := cartItem.GetProduct()

	var expectedSku = constTestSku + constTestSkuModifier
	var gotSku = cartItemProduct.GetSku()
	if gotSku != expectedSku {
		t.Errorf("Incorrect Sku. Expected: '%v'. Got: '%v'.", expectedSku, gotSku)
	}

	var gotPrice = cartItemProduct.GetPrice()
	if gotPrice != constExpectedPrice {
		t.Errorf("Incorrect Price. Expected: '%v'. Got: '%v'.", constExpectedPrice, gotPrice)
	}

	fmt.Println("testCartAddItem DONE")
}

func testAPICartInfo(t *testing.T, applicationContext api.InterfaceApplicationContext) {
	loadedCart, err := cartActor.APICartInfo(applicationContext)
	if err != nil {
		t.Error(err)
	}

	var loadedCartMap = loadedCart.(map[string]interface{})
	var items = loadedCartMap["items"].([]map[string]interface{})
	var itemProduct = items[0]["product"].(map[string]interface{})

	var expectedSku = constTestSku + constTestSkuModifier
	var gotSku = itemProduct["sku"].(string)
	if gotSku != expectedSku {
		t.Errorf("Incorrect Sku. Expected: '%v'. Got: '%v'.", expectedSku, gotSku)
	}

	var gotPrice = itemProduct["price"].(float64)
	if gotPrice != constExpectedPrice {
		t.Errorf("Incorrect Price. Expected: '%v'. Got: '%v'.", constExpectedPrice, gotPrice)
	}

	fmt.Println("testAPICartInfo DONE")
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
