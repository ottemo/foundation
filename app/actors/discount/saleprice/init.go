package saleprice

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
)

// init makes package self-initialization routine
func init() {
	salePriceInstance := new(DefaultSalePrice)
	var _ saleprice.InterfaceSalePrice = salePriceInstance
	if err := models.RegisterModel(saleprice.ConstModelNameSalePrice, salePriceInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "402b5049-cac4-43ba-a362-a53de170f696", err.Error())
	}

	salePriceCollectionInstance := new(DefaultSalePriceCollection)
	var _ saleprice.InterfaceSalePriceCollection = salePriceCollectionInstance
	if err := models.RegisterModel(saleprice.ConstSalePriceDbCollectionName, salePriceCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6298a58b-b996-42de-9ba8-84f7e336fd7e", err.Error())
	}

	var _ checkout.InterfacePriceAdjustment = salePriceInstance
	if err := checkout.RegisterPriceAdjustment(salePriceInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b1abee9c-f64c-4605-90c2-4cb9365568b3", err.Error())
	}

	db.RegisterOnDatabaseStart(salePriceInstance.setupDB)
	api.RegisterOnRestServiceStart(setupAPI)

	salePriceDelegate = new(SalePriceDelegate)
	env.RegisterOnConfigStart(setupConfig)
}

// setupDB prepares system database for package usage
func (it *DefaultSalePrice) setupDB() error {
	dbSalePriceCollection, err := db.GetCollection(saleprice.ConstSalePriceDbCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := dbSalePriceCollection.AddColumn("amount", db.ConstTypeMoney, false); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7cf84ca1-acec-419a-86b3-eb11b6a1d4cf", err.Error())
	}
	if err := dbSalePriceCollection.AddColumn("end_datetime", db.ConstTypeDatetime, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fed6cf3f-05f9-4651-a8b9-361c404fc020", err.Error())
	}
	if err := dbSalePriceCollection.AddColumn("product_id", db.ConstTypeID, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1333660c-1aa5-4cf2-bb52-89e037b866dc", err.Error())
	}
	if err := dbSalePriceCollection.AddColumn("start_datetime", db.ConstTypeDatetime, true); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "24a2439f-89bc-4606-82ad-992ca6bdd563", err.Error())
	}

	return nil

}
