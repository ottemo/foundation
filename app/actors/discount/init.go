package discount

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultDiscount)
	var _ checkout.InterfaceDiscount = instance
	checkout.RegisterDiscount(instance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {

	collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
	if err != nil {
		env.ErrorDispatch(err)
	}

	collection.AddColumn("code", "text", true)
	collection.AddColumn("name", "text", false)
	collection.AddColumn("amount", "decimal", false)
	collection.AddColumn("percent", "decimal", false)
	collection.AddColumn("times", "int", false)
	collection.AddColumn("since", "datetime", false)
	collection.AddColumn("until", "datetime", false)

	return nil
}
