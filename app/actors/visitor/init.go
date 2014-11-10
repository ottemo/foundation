package visitor

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// package self initializer
func init() {
	visitorInstance := new(DefaultVisitor)
	var _ visitor.I_Visitor = visitorInstance
	models.RegisterModel(visitor.ModelNameVisitor, visitorInstance)

	visitorCollectionInstance := new(DefaultVisitorCollection)
	var _ visitor.I_VisitorCollection = visitorCollectionInstance
	models.RegisterModel(visitor.ModelNameVisitorCollection, visitorCollectionInstance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setups database tables for model usage
func setupDB() error {

	collection, err := db.GetCollection(CollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("email", "id", true)
	collection.AddColumn("validate", "varchar(128)", false)
	collection.AddColumn("password", "varchar(128)", false)
	collection.AddColumn("first_name", "varchar(50)", true)
	collection.AddColumn("last_name", "varchar(50)", true)

	collection.AddColumn("facebook_id", "varchar(100)", true)
	collection.AddColumn("google_id", "varchar(100)", true)

	collection.AddColumn("billing_address_id", "id", false)
	collection.AddColumn("shipping_address_id", "id", false)

	collection.AddColumn("birthday", "datetime", false)
	collection.AddColumn("is_admin", "bool", false)
	collection.AddColumn("created_at", "datetime", false)

	return nil
}
