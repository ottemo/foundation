package address

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// module entry point before app start
func init() {
	visitorAddressInstance := new(DefaultVisitorAddress)
	var _ visitor.I_VisitorAddress = visitorAddressInstance
	models.RegisterModel(visitor.ModelNameVisitorAddress, visitorAddressInstance)

	visitorAddressCollectionInstance := new(DefaultVisitorAddressCollection)
	var _ visitor.I_VisitorAddressCollection = visitorAddressCollectionInstance
	models.RegisterModel(visitor.ModelNameVisitorAddressCollection, visitorAddressCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {
	collection, err := db.GetCollection(COLLECTION_NAME_VISITOR_ADDRESS)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("visitor_id", "id", false)
	collection.AddColumn("first_name", "varchar(100)", false)
	collection.AddColumn("last_name", "varchar(100)", false)
	collection.AddColumn("company", "varchar(100)", false)
	collection.AddColumn("address_line1", "varchar(255)", false)
	collection.AddColumn("address_line2", "varchar(255)", false)
	collection.AddColumn("country", "varchar(50)", false)
	collection.AddColumn("state", "varchar(2)", false)
	collection.AddColumn("city", "varchar(100)", false)
	collection.AddColumn("phone", "varchar(100)", false)
	collection.AddColumn("zip_code", "varchar(10)", false)

	return nil
}
