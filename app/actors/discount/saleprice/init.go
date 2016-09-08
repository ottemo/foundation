package saleprice

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
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
