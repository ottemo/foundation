package address

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	visitorAddressInstance := new(DefaultVisitorAddress)
	var _ visitor.InterfaceVisitorAddress = visitorAddressInstance
	models.RegisterModel(visitor.ConstModelNameVisitorAddress, visitorAddressInstance)

	visitorAddressCollectionInstance := new(DefaultVisitorAddressCollection)
	var _ visitor.InterfaceVisitorAddressCollection = visitorAddressCollectionInstance
	models.RegisterModel(visitor.ConstModelNameVisitorAddressCollection, visitorAddressCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameVisitorAddress)
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
