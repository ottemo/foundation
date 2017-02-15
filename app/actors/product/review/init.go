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

		if err := collection.AddColumn("product_id", db.ConstTypeID, true); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d9ed39ea-d31d-4a56-9dd6-b5792cb2fcc1", err.Error())
		}
		if err := collection.AddColumn("visitor_id", db.ConstTypeID, true); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8cf09bcd-d8a7-4470-968f-66b465986971", err.Error())
		}
		if err := collection.AddColumn("username", db.TypeWPrecision(db.ConstTypeVarchar, 100), true); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1622dd0a-389e-4e2b-bb2a-b475352fa1e1", err.Error())
		}
		if err := collection.AddColumn("rating", db.ConstTypeInteger, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bd776ff5-ffdc-45cf-8de1-69bc954f900b", err.Error())
		}
		if err := collection.AddColumn("review", db.ConstTypeText, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4381bafb-e72d-49e2-b1ad-83419fc14d1c", err.Error())
		}
		if err := collection.AddColumn("created_at", db.ConstTypeDatetime, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5eff2436-7067-40ee-b54a-2a35d2764379", err.Error())
		}
		if err := collection.AddColumn("approved", db.ConstTypeBoolean, true); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1f422b61-9ba4-4b9f-a0eb-fd02d1a2f3d0", err.Error())
		}

		if shouldFillApprovedField {
			env.Log(ConstErrorModule, env.ConstLogPrefixInfo, "Field 'approved' have been added. Make all reviews approved.")
			if err := fillApprovedField(); err != nil {
				return env.ErrorDispatch(err)
			}
		}
	} else {
		return env.ErrorDispatch(err)
	}

	if collection, err := dbEngine.GetCollection("rating"); err == nil {
		if err := collection.AddColumn("product_id", db.ConstTypeID, true); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "aafc3eb0-85ae-4a89-96b5-3d4414c1f3d6", err.Error())
		}
		if err := collection.AddColumn("stars_1", db.ConstTypeInteger, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f99c1c16-952b-4d26-8bf3-ece9e6ffa698", err.Error())
		}
		if err := collection.AddColumn("stars_2", db.ConstTypeInteger, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "64f7c52c-8c25-4c43-9a28-d9999685c698", err.Error())
		}
		if err := collection.AddColumn("stars_3", db.ConstTypeInteger, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "91a44165-b7e8-45a7-bc90-ec130d917c70", err.Error())
		}
		if err := collection.AddColumn("stars_4", db.ConstTypeInteger, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dfcce367-89d7-4d99-8395-2d90e0ae56e3", err.Error())
		}
		if err := collection.AddColumn("stars_5", db.ConstTypeInteger, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "67fec2dd-7138-44d1-8c56-b56c37cddb0f", err.Error())
		}
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
