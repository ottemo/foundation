package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ottemo/foundation/app"

	// using standard set of packages
	_ "github.com/ottemo/foundation/basebuild"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/product"
	//productActor "github.com/ottemo/foundation/app/actors/product"
)

func init() {
	// time.Unix() should be in UTC (as it could be not by default)
	time.Local = time.UTC
}

// executable file start point
func mainProductOptonsUpdate() {
	defer app.End() // application close event

	// application start event
	if err := app.Start(); err != nil {
		env.LogError(err)
		fmt.Println(err.Error())
		os.Exit(0)
	}

	// get product collection
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		fmt.Println(env.ErrorDispatch(err))
	}

	// update products option
	for _, currentProduct := range productCollection.ListProducts() {
//		newOptions := productActor.UpdateProductOptions(currentProduct)
//		currentProduct.Set("options", newOptions)
		err := currentProduct.Save()
		if err != nil {
			fmt.Println(env.ErrorDispatch(err))
		}
	}
}
