package product

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// init makes package self-initialization routine
func init() {
	productInstance := new(DefaultProduct)
	var _ product.InterfaceProduct = productInstance
	models.RegisterModel(product.ConstModelNameProduct, productInstance)

	collectionInstance := new(DefaultProductCollection)
	var _ product.InterfaceProductCollection = collectionInstance
	models.RegisterModel(product.ConstModelNameProductCollection, collectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	var shouldFillVisibleField = !collection.HasColumn("visible")

	collection.AddColumn("enabled", db.ConstTypeBoolean, true)
	collection.AddColumn("sku", db.ConstTypeVarchar, true)
	collection.AddColumn("name", db.ConstTypeVarchar, true)
	collection.AddColumn("short_description", db.ConstTypeVarchar, false)
	collection.AddColumn("description", db.ConstTypeText, false)
	collection.AddColumn("default_image", db.ConstTypeVarchar, false)
	collection.AddColumn("price", db.ConstTypeMoney, false)
	collection.AddColumn("weight", db.ConstTypeFloat, false)
	collection.AddColumn("options", db.ConstTypeJSON, false)
	collection.AddColumn("visible", db.ConstTypeBoolean, false)

	collection.AddColumn("related_pids", db.TypeArrayOf(db.ConstTypeID), false)

	if shouldFillVisibleField {
		env.Log(ConstErrorModule, env.ConstLogPrefixInfo, "Field 'visible' have been added. Make all products visible.")
		app.OnAppStart(fillVisibleField)
	} else {
		env.Log(ConstErrorModule, env.ConstLogPrefixInfo, "'visible' value need no update.")
	}

	return nil
}

// fillVisibleField makes all product visible on adding "visible" column
func fillVisibleField() error {
	// get product collection
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// update products
	for _, currentProduct := range productCollection.ListProducts() {
		err = currentProduct.Set("visible", true)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err := currentProduct.Save()
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}
