package tax

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/checkout"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultTax)

	checkout.RegisterPriceAdjustment(instance)

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
			collection.AddColumn("code", db.ConstTypeVarchar, true)
			collection.AddColumn("country", db.ConstTypeVarchar, true)
			collection.AddColumn("state", db.ConstTypeVarchar, true)
			collection.AddColumn("zip", db.ConstTypeVarchar, false)
			collection.AddColumn("rate", db.ConstTypeDecimal, false)
		} else {
			return env.ErrorDispatch(err)
		}
	} else {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "e132f4b0-d69d-4900-b2de-24b96d0fc1ce", "Can't get database engine")
	}

	return nil
}
