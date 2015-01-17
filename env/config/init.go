package config

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/impex"
)

// init makes package self-initialization routine
func init() {
	instance := &DefaultConfig{
		configValues:     make(map[string]interface{}),
		configTypes:      make(map[string]string),
		configValidators: make(map[string]env.FuncConfigValueValidator)}

	db.RegisterOnDatabaseStart(setupDB)
	db.RegisterOnDatabaseStart(instance.Load)

	api.RegisterOnRestServiceStart(setupAPI)

	env.RegisterConfig(instance)

	impex.RegisterImpexModel("Config", instance)
}

// setupDB prepares system database for package usage
func setupDB() error {
	collection, err := db.GetCollection(ConstCollectionNameConfig)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddColumn("path", db.ConstTypeVarchar, true)
	collection.AddColumn("value", db.ConstTypeText, false)

	collection.AddColumn("type", db.ConstTypeVarchar, false)

	collection.AddColumn("editor", db.ConstTypeVarchar, false)
	collection.AddColumn("options", db.ConstTypeText, false)

	collection.AddColumn("label", db.ConstTypeVarchar, false)
	collection.AddColumn("description", db.ConstTypeText, false)

	collection.AddColumn("image", db.ConstTypeVarchar, false)

	return nil
}
