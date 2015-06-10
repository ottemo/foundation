package product

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/seo"
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

	seo.RegisterSEOType(ConstSEOTypeProduct, "product/", product.ConstModelNameProduct)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("enabled", db.ConstTypeBoolean, true)
	collection.AddColumn("sku", db.ConstTypeVarchar, true)
	collection.AddColumn("name", db.ConstTypeVarchar, true)
	collection.AddColumn("short_description", db.ConstTypeVarchar, false)
	collection.AddColumn("description", db.ConstTypeText, false)
	collection.AddColumn("default_image", db.ConstTypeVarchar, false)
	collection.AddColumn("price", db.ConstTypeMoney, false)
	collection.AddColumn("weight", db.ConstTypeFloat, false)
	collection.AddColumn("options", db.ConstTypeJSON, false)

	collection.AddColumn("related_pids", db.TypeArrayOf(db.ConstTypeID), false)

	return nil
}
