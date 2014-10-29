package main

import (
	"testing"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/product"
)


func TestGetAllProducts(t *testing.T) {
	app.Start()

	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		t.Error(err)
	}
	_, err = productCollection.List()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGetAllProducts(b *testing.B) {
	app.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		productCollection, err := product.GetProductCollectionModel()
		if err != nil {
			b.Error(err)
		}

		_, err = productCollection.List()
		if err != nil {
			b.Error(err)
		}
	}

}

// func BenchmarkGetAllProductsParallel(b *testing.B) {
// 	app.Start()
//
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			productCollection, err := product.GetProductCollectionModel()
// 			if err != nil {
// 				b.Error(err)
// 			}
//
// 			_, err = productCollection.List()
// 			if err != nil {
// 				b.Error(err)
// 			}
// 		}
// 	})
// }
