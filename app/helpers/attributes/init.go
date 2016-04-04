package attributes

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models"
)

// init makes package self-initialization routine
func init() {
	var _ models.InterfaceCustomAttributes = new(ModelCustomAttributes)
	var _ models.InterfaceExternalAttributes = new(ModelExternalAttributes)

	db.RegisterOnDatabaseStart(setupDB)
}

// setupDB prepares database for a package
func setupDB() error {

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
