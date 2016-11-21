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
		var shouldFillApprovedField = !collection.HasColumn("approved")

		collection.AddColumn("product_id", db.ConstTypeID, true)
		collection.AddColumn("visitor_id", db.ConstTypeID, true)
		collection.AddColumn("username", db.TypeWPrecision(db.ConstTypeVarchar, 100), true)
		collection.AddColumn("rating", db.ConstTypeInteger, false)
		collection.AddColumn("review", db.ConstTypeText, false)
		collection.AddColumn("created_at", db.ConstTypeDatetime, false)
		collection.AddColumn("approved", db.ConstTypeBoolean, true)

		if shouldFillApprovedField {
			env.Log(ConstErrorModule, env.ConstLogPrefixInfo, "Field 'approved' have been added. Make all reviews approved.")
			if err := fillApprovedField(); err != nil {
				return env.ErrorDispatch(err)
			}
		} else {
			env.Log(ConstErrorModule, env.ConstLogPrefixInfo, "'approved' value need no update.")
		}
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

// fillApprovedField makes all reviews approved on adding "approved" column
func fillApprovedField() error {
	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "f906658c-fbfb-4e2a-a266-9d703bc199ab", "Can't get database engine")
	}

	// get reviews collection
	reviewCollection, err := dbEngine.GetCollection("review")
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// update reviews
	reviewMaps, err := reviewCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, currentReview := range reviewMaps {
		currentReview["approved"] = true

		_, err := reviewCollection.Save(currentReview)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}
