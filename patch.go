package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/utils"

	// using standard set of packages

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/subscription"
	_ "github.com/ottemo/foundation/basebuild"
	"github.com/ottemo/foundation/env"
)

func init() {
	// time.Unix() should be in UTC (as it could be not by default)
	time.Local = time.UTC

}

// executable file start point
func main1() {
	defer func (){
		if err := app.End(); err != nil { // application close event
			fmt.Println(err.Error())
		}
	}()

	fmt.Println("starting application " + time.Now().String())
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

	fmt.Println("update products options" + time.Now().String())
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

	// get product collection
	subscriptionCollection, err := subscription.GetSubscriptionCollectionModel()
	if err != nil {
		fmt.Println(env.ErrorDispatch(err))
	}
	currentCart, err := cart.GetCartModel()
	if err != nil {
		fmt.Println(env.ErrorDispatch(err))
	}

	for _, currentSubscription := range subscriptionCollection.ListSubscriptions() {
		for _, subscriptionItem := range currentSubscription.GetItems() {
			updatedOptions := make(map[string]interface{})
			// Labels where used as a key for options key: value, so we will convert both of them
			for optionKey, optionValue := range subscriptionItem.Options {
				updatedOptions[utils.StrToSnakeCase(optionKey)] = utils.StrToSnakeCase(utils.InterfaceToString(optionValue))
			}
			subscriptionItem.Options = updatedOptions
			if _, err = currentCart.AddItem(subscriptionItem.ProductID, subscriptionItem.Qty, subscriptionItem.Options); err != nil {
				fmt.Println(env.ErrorDispatch(err))
				fmt.Println(subscriptionItem.Options)
			}
		}

		err = currentSubscription.Save()
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
