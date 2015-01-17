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
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "da2fa1c9-5ca3-46dd-a5ea-3fae1e7b9614", "Can't get database engine")
	}

	if collection, err := dbEngine.GetCollection("review"); err == nil {
		collection.AddColumn("product_id", db.ConstTypeID, true)
		collection.AddColumn("visitor_id", db.ConstTypeID, true)
		collection.AddColumn("username", db.TypeWPrecision(db.ConstTypeVarchar, 100), true)
		collection.AddColumn("rating", db.ConstTypeInteger, false)
		collection.AddColumn("review", db.ConstTypeText, false)
		collection.AddColumn("created_at", db.ConstTypeDatetime, false)
	} else {
		return env.ErrorDispatch(err)
	}

	if collection, err := dbEngine.GetCollection("rating"); err == nil {
		collection.AddColumn("product_id", db.ConstTypeID, true)
		collection.AddColumn("stars_1", db.ConstTypeInteger, false)
		collection.AddColumn("stars_2", db.ConstTypeInteger, false)
		collection.AddColumn("stars_3", db.ConstTypeInteger, false)
		collection.AddColumn("stars_4", db.ConstTypeInteger, false)
		collection.AddColumn("stars_5", db.ConstTypeInteger, false)
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
