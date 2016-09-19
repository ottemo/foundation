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
	models.RegisterModel(saleprice.ConstModelNameSalePrice, salePriceInstance)

	salePriceCollectionInstance := new(DefaultSalePriceCollection)
	var _ saleprice.InterfaceSalePriceCollection = salePriceCollectionInstance
	models.RegisterModel(saleprice.ConstSalePriceDbCollectionName, salePriceCollectionInstance)

	var _ checkout.InterfacePriceAdjustment = salePriceInstance
	checkout.RegisterPriceAdjustment(salePriceInstance)

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

	dbSalePriceCollection.AddColumn("amount", db.ConstTypeMoney, false)
	dbSalePriceCollection.AddColumn("end_datetime", db.ConstTypeDatetime, true)
	dbSalePriceCollection.AddColumn("product_id", db.ConstTypeID, true)
	dbSalePriceCollection.AddColumn("start_datetime", db.ConstTypeDatetime, true)

	return nil

}
