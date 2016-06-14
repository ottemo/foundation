package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/utils"

	// using standard set of packages

	"github.com/ottemo/foundation/app/models/product"
	_ "github.com/ottemo/foundation/basebuild"
	"github.com/ottemo/foundation/env"
)

func init() {
	// time.Unix() should be in UTC (as it could be not by default)
	time.Local = time.UTC
}

// executable file start point
func main() {
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
		newOptions := ConvertProductOptionsToSnakeCase(currentProduct)
		err = currentProduct.Set("options", newOptions)
		if err != nil {
			fmt.Println(env.ErrorDispatch(err))
		}

		err := currentProduct.Save()
		if err != nil {
			fmt.Println(env.ErrorDispatch(err))
		}
	}
}

// ConvertProductOptionsToSnakeCase updates option keys for product to case_snake
func ConvertProductOptionsToSnakeCase(product product.InterfaceProduct) map[string]interface{} {

	newOptions := make(map[string]interface{})

	// product options
	for optionsName, currentOption := range product.GetOptions() {
		currentOption := utils.InterfaceToMap(currentOption)

		if option, present := currentOption["options"]; present {
			newOptionValues := make(map[string]interface{})

			// option values
			for key, value := range utils.InterfaceToMap(option) {
				newOptionValues[utils.StrToSnakeCase(key)] = value
			}

			currentOption["options"] = newOptionValues
		}
		newOptions[utils.StrToSnakeCase(optionsName)] = currentOption
	}

	return newOptions
}
