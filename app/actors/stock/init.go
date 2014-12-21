package stock

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {

	instance := new(DefaultStock)
	var _ product.InterfaceStock = instance

	api.RegisterOnRestServiceStart(setupAPI)
	db.RegisterOnDatabaseStart(setupDB)
	env.RegisterOnConfigStart(setupConfig)
}

// setupDB prepares system database for package usage
func setupDB() error {

	if collection, err := db.GetCollection(ConstCollectionNameStock); err == nil {
		collection.AddColumn("product_id", db.ConstDBBasetypeID, true)
		collection.AddColumn("options", db.ConstDBBasetypeJSON, true)
		collection.AddColumn("qty", db.ConstDBBasetypeInteger, false)
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
