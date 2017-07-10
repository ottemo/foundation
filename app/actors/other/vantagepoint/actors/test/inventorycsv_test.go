package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ottemo/foundation/app/actors/other/vantagepoint/actors"
	"github.com/ottemo/foundation/app/actors/stock"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/product"
)

// --------------------------------------------------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := test.StartAppInTestingMode()
	if err != nil {
		fmt.Println("Unable to start app in testing mode:", err)
	}

	os.Exit(m.Run())
}

func TestInventoryCSVProcess(t *testing.T) {
	if err := prepareEnvironment(t); err != nil {
		t.Fatal(err)
	}

	processor, err := actors.NewInventoryCSV(&testEnv{})
	if err != nil {
		t.Fatal(err)
	}

	csvStr := `UPC Number,Stock
no-stock,10
have-stock,20`

	if err := processor.Process(strings.NewReader(csvStr)); err != nil {
		t.Fatal(err)
	}

	_ = testProductQuantity(t, "no-stock", 10)
	_ = testProductQuantity(t, "have-stock", 20)
	_ = testProductQuantity(t, "configurable", 23)
}

func testProductQuantity(t *testing.T, sku string, expectedQty int) error {
	products, err := getProductsBySku(t, sku)
	if err != nil {
		t.Fatal(err)
	}

	stockMgr := product.GetRegisteredStock()
	if stockMgr == nil {
		t.Fatal("stock is undefined")
	}

	for _, productModel := range products {
		qty := stockMgr.GetProductQty(productModel.GetID(), map[string]interface{}{})
		if qty != expectedQty {
			t.Errorf("sku [%s] qty [%d] expected [%d]", sku, qty, expectedQty)
		}
	}

	return nil
}

func prepareEnvironment(t *testing.T) error {
	prepareConfig(t)

	if err := prepareData(t); err != nil {
		return err
	}

	_ = createProductFromJson(t, `{"sku": "no-stock", "name": "NAME no-stock"}`)
	product02 := createProductFromJson(t, `{"sku": "have-stock", "name": "NAME have-stock",
		"inventory": [
		    {"options": {}, "qty": 5}
		]
	}`)
	_ = createProductFromJson(t, `{"sku": "configurable", "name": "NAME configurable",
		"options": {
			"color": {"key": "color",
				"options": {
					"red_2": {
						"key": "red_2",
						"_ids": ["`+product02.GetID()+`"]
					},
					"red_3": {"key": "red_3"}
				}
			}
		},
		"inventory": [
		    {"options": { }, "qty": 5},
		    {"options": {"color": "red_2"}, "qty": 2},
		    {"options": {"color": "red_3"}, "qty": 3}
		]
	}`)

	return nil
}

func prepareData(t *testing.T) error {
	// clear products
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return err
	}

	_, err = productCollection.GetDBCollection().Delete()
	if err != nil {
		return err
	}

	// clear inventory
	stockCollection, err := db.GetCollection(stock.ConstCollectionNameStock)
	if err != nil {
		return err
	}

	_, err = stockCollection.Delete()

	return err
}

func createProductFromJson(t *testing.T, json string) product.InterfaceProduct {
	productData, err := utils.DecodeJSONToStringKeyMap(json)
	if err != nil {
		fmt.Println("json issue: " + json)
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

func prepareConfig(t *testing.T) {
	if config := env.GetConfig(); config != nil {
		if config.GetValue("general.stock.enabled") != true {
			err := env.GetConfig().SetValue("general.stock.enabled", true)
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func getProductsBySku(t *testing.T, sku string) ([]product.InterfaceProduct, error) {
	collection, err := product.GetProductCollectionModel()
	if err != nil {
		return nil, err
	}

	if err := collection.ListFilterAdd("sku", "=", sku); err != nil {
		return nil, err
	}

	return collection.ListProducts(), nil
}
