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

	collection.AddColumn("code", db.ConstTypeVarchar, true)
	collection.AddColumn("name", db.ConstTypeVarchar, false)
	collection.AddColumn("amount", db.ConstTypeDecimal, false)
	collection.AddColumn("percent", db.ConstTypeDecimal, false)
	collection.AddColumn("times", db.ConstTypeInteger, false)
	collection.AddColumn("since", db.ConstTypeDatetime, false)
	collection.AddColumn("until", db.ConstTypeDatetime, false)

	return nil
}
