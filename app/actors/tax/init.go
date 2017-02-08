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

	if err := checkout.RegisterPriceAdjustment(instance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "e5f12dd7-70ca-4563-aa51-1bfd913ac797", err.Error())
	}

	db.RegisterOnDatabaseStart(setupDB)
	api.RegisterOnRestServiceStart(setupAPI)
}

// setupDB prepares system database for package usage
func setupDB() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
			if err := collection.AddColumn("code", db.ConstTypeVarchar, true); err != nil {
				return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "6585ff4b-e079-4e9e-9784-d1479814ae19", err.Error())
			}
			if err := collection.AddColumn("country", db.ConstTypeVarchar, true); err != nil {
				return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "ee6fdcd8-be30-480f-8f81-5a258a14c36e", err.Error())
			}
			if err := collection.AddColumn("state", db.ConstTypeVarchar, true); err != nil {
				return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "efce92d6-0f42-4188-88d3-bca828bda089", err.Error())
			}
			if err := collection.AddColumn("zip", db.ConstTypeVarchar, false); err != nil {
				return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "1a93db1c-4d6a-43cb-92f0-5cf59e920c9b", err.Error())
			}
			if err := collection.AddColumn("rate", db.ConstTypeDecimal, false); err != nil {
				return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "d4cea0ed-a7f2-4f2d-9c05-ee7d158be092", err.Error())
			}
		} else {
			return env.ErrorDispatch(err)
		}
	} else {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "e132f4b0-d69d-4900-b2de-24b96d0fc1ce", "Can't get database engine")
	}

	return nil
}
