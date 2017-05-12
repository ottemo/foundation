package taxcloud

import (
	"github.com/ottemo/foundation/app/actors/tax/taxcloud/model"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

func init() {
	productTicInstance := new(DefaultProductTic)
	var _ model.InterfaceProductTic = productTicInstance
	if err := models.RegisterModel(model.ConstProductTicModelName, productTicInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a9800d4d-9112-4e48-b47b-254969822c25", err.Error())
	}

	taxCloudInstance := new(TaxCloudPriceAdjustment)
	var _ checkout.InterfacePriceAdjustment = taxCloudInstance
	if err := checkout.RegisterPriceAdjustment(taxCloudInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7c78f9e6-5f47-48ac-8597-900d72429fc4", err.Error())
	}

	db.RegisterOnDatabaseStart(setupDB)

	ticDelegate = new(TicDelegate)
	env.RegisterOnConfigStart(setupConfig)
}

// setupDB prepares system database for package usage
func setupDB() error {
	dbCollection, err := db.GetCollection(model.ConstProductTicCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := dbCollection.AddColumn(ConstTicIdAttribute, db.ConstTypeInteger, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "cbbad0a4-6ccf-4f9e-b394-53190e8301c6", err.Error())
	}
	if err := dbCollection.AddColumn("product_id", db.ConstTypeID, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7e3b991e-e4c3-4fd8-bf55-0a926a239e3e", err.Error())
	}

	return nil

}
