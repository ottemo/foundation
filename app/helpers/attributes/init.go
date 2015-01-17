package attributes

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	db.RegisterOnDatabaseStart(SetupDB)
}

// SetupDB prepares system database for package usage
func SetupDB() error {

	if collection, err := db.GetCollection("custom_attributes"); err == nil {
		collection.AddColumn("model", db.ConstTypeVarchar, true)
		collection.AddColumn("collection", db.ConstTypeVarchar, true)
		collection.AddColumn("attribute", db.ConstTypeVarchar, true)
		collection.AddColumn("type", db.ConstTypeVarchar, false)
		collection.AddColumn("required", db.ConstTypeBoolean, false)
		collection.AddColumn("label", db.ConstTypeVarchar, true)
		collection.AddColumn("group", db.ConstTypeVarchar, false)
		collection.AddColumn("editors", db.ConstTypeVarchar, false)
		collection.AddColumn("options", db.ConstTypeText, false)
		collection.AddColumn("default", db.ConstTypeText, false)
		collection.AddColumn("validators", db.ConstTypeVarchar, false)
		collection.AddColumn("layered", db.ConstTypeBoolean, false)
		collection.AddColumn("public", db.ConstTypeBoolean, false)

	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
