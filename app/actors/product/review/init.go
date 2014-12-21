package review

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine before app start
func init() {
	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// DB preparations for current model implementation
func setupDB() error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "da2fa1c95ca346dda5ea3fae1e7b9614", "Can't get database engine")
	}

	if collection, err := dbEngine.GetCollection("review"); err == nil {
		collection.AddColumn("product_id", "id", true)
		collection.AddColumn("visitor_id", "id", true)
		collection.AddColumn("username", "varchar(100)", true)
		collection.AddColumn("rating", "int", false)
		collection.AddColumn("review", "text", false)
		collection.AddColumn("created_at", "datetime", false)
	} else {
		return env.ErrorDispatch(err)
	}

	if collection, err := dbEngine.GetCollection("rating"); err == nil {
		collection.AddColumn("product_id", "id", true)
		collection.AddColumn("stars_1", "int", false)
		collection.AddColumn("stars_2", "int", false)
		collection.AddColumn("stars_3", "int", false)
		collection.AddColumn("stars_4", "int", false)
		collection.AddColumn("stars_5", "int", false)
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
