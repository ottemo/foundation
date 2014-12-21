package product

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
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

	collection.AddColumn("enabled", db.ConstDBBasetypeBoolean, true)
	collection.AddColumn("sku", db.ConstDBBasetypeVarchar, true)
	collection.AddColumn("name", db.ConstDBBasetypeVarchar, true)
	collection.AddColumn("short_description", db.ConstDBBasetypeVarchar, false)
	collection.AddColumn("description", db.ConstDBBasetypeText, false)
	collection.AddColumn("default_image", db.ConstDBBasetypeVarchar, false)
	collection.AddColumn("price", db.ConstDBBasetypeMoney, false)
	collection.AddColumn("weight", db.ConstDBBasetypeFloat, false)
	collection.AddColumn("options", db.ConstDBBasetypeJSON, false)

	collection.AddColumn("related_pids", "[]"+db.ConstDBBasetypeID, false)

	return nil
}
