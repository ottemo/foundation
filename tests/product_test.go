package tests

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/utils"
)

// function used to test most product model operations
func TestProductsOperations(tst *testing.T) {

	// starting application and getting product model
	err := StartAppInTestingMode()
	if err != nil {
		tst.Error(err)
	}

	productModel, err := product.GetProductModel()
	if err != nil {
		tst.Error(err)
	}

	// looking for "test" custom attribute
	found := false
	for _, attributeInfo := range productModel.GetAttributesInfo() {
		if attributeInfo.Attribute == "test" {
			found = true
			break
		}
	}
	if found {
		err = productModel.RemoveAttribute("test")
		if err != nil {
			tst.Error(err)
		}
	}

	// adding new attribute to system
	err = productModel.AddNewAttribute(models.StructAttributeInfo{
		Model:      product.ConstModelNameProduct,
		Collection: "product", // TODO: Custom attribute helper should set this by self
		Attribute:  "test",
		Type:       utils.ConstDataTypeText,
		IsRequired: true,
		IsStatic:   false,
		Label:      "Test Attribute",
		Group:      "General",
		Editors:    "text",
		Options:    "",
		Default:    "",
		Validators: "",
		IsLayered:  true,
	})
	if err != nil {
		tst.Error(err)
	}

	// making data for test object
	productData, err := utils.DecodeJSONToStringKeyMap(`{
		"sku": "test",
		"name": "Test Product",
		"short_description": "something short",
		"description": "something long",
		"default_image": "",
		"price": 1.1,
		"weight": 0.5,
		"test": "ok",
		"options" : {
			"Color" : {
				"type" : "select", "required" : true, "price_type" : "fixed", "label" : "Color",
				"options" : {
					"black": {"order": "3", "label": "black", "price": 1.3, "price_type": "percent", "sku": "black"},
					"blue":  {"order": "1", "label": "blue",  "price": 2.0, "price_type": "percent", "sku": "blue"},
					"red":   {"order": "2", "label": "red",   "price": 100, "price_type": "percent", "sku": "red"}
				}
			}
		}}`)
	if err != nil {
		tst.Error(err)
	}

	// setting values for product
	err = productModel.FromHashMap(productData)
	if err != nil {
		tst.Error(err)
	}

	// saving product
	err = productModel.Save()
	if err != nil {
		tst.Error(err)
	}

	// loading just saved product
	productModel, err = product.LoadProductByID(productModel.GetID())
	if err != nil {
		tst.Error(err)
	}

	// deleting product
	err = productModel.Delete()
	if err != nil {
		tst.Error(err)
	}

	// removing added attribute
	err = productModel.RemoveAttribute("test")
	if err != nil {
		tst.Error(err)
	}
}

// idle benchmark test, to check go bench not lying
func BenchmarkSleep1sec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(1000000000)
	}
}

// benchmarks product list obtain operation
func BenchmarkList50Products(b *testing.B) {
	err := StartAppInTestingMode()
	if err != nil {
		b.Error(err)
	}

	err = MakeSureProductsCount(100)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		productCollection, err := product.GetProductCollectionModel()
		if err != nil {
			b.Error(err)
		}

		productCollection.ListLimit(0, 50)
		if err != nil {
			b.Error(err)
		}

		_, err = productCollection.List()
		if err != nil {
			b.Error(err)
		}
	}
}

// benchmarks product load operation throught collection
func BenchmarkProductCntLoad(b *testing.B) {
	err := StartAppInTestingMode()
	if err != nil {
		b.Error(err)
	}

	err = MakeSureProductsCount(100)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		productCollection, err := product.GetProductCollectionModel()
		if err != nil {
			b.Error(err)
		}

		_, err = productCollection.GetDBCollection().Count()
		if err != nil {
			b.Error(err)
		}
		//productsCount := 100

		//		err = productCollection.ListLimit(rand.Intn(productsCount), 1)
		//		if err != nil {
		//			b.Error(err)
		//		}
		//
		//		productCollection.ListProducts()
	}
}

// benchmarks product load operation
func BenchmarkProductLoad(b *testing.B) {
	err := StartAppInTestingMode()
	if err != nil {
		b.Error(err)
	}

	err = MakeSureProductsCount(100)
	if err != nil {
		b.Error(err)
	}

	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		b.Error(err)
	}

	productDBCollection := productCollection.GetDBCollection()
	productDBCollection.SetResultColumns("_id")
	productIds, err := productDBCollection.Load()
	count := len(productIds)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		randomID := utils.InterfaceToString(productIds[rand.Intn(count)]["_id"])
		product.LoadProductByID(randomID)
	}
}

// func BenchmarkGetAllProductsParallel(b *testing.B) {
// 	app.Start()
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			productCollection, err := product.GetProductCollectionModel()
// 			if err != nil {
// 				b.Error(err)
// 			}
// 			_, err = productCollection.List()
// 			if err != nil {
// 				b.Error(err)
// 			}
// 		}
// 	})
// }
