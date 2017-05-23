package stock

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/stock"
)

// init makes package self-initialization routine
func init() {

	instance := new(DefaultStock)
	var _ stock.InterfaceStock = instance
	if err := models.RegisterModel(stock.ConstModelNameStock, instance); err != nil {
		_ = env.ErrorDispatch(err)
	}

	stockCollectionInstance := new(DefaultStockCollection)
	var _ stock.InterfaceStockCollection = stockCollectionInstance
	if err := models.RegisterModel(stock.ConstModelNameStockCollection, stockCollectionInstance); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "42552c35-a4ef-40e6-aa38-5d69d4e92578", err.Error())
	}

	stockDelegate = new(StockDelegate)
	api.RegisterOnRestServiceStart(setupAPI)
	db.RegisterOnDatabaseStart(setupDB)
	env.RegisterOnConfigStart(setupConfig)
}

// setupDB prepares system database for package usage
func setupDB() error {

	if collection, err := db.GetCollection(ConstCollectionNameStock); err == nil {
		if err := collection.AddColumn("product_id", db.ConstTypeID, true); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d7641743-7d3a-4e1e-9627-fc4b1fad85d1", err.Error())
		}
		if err := collection.AddColumn("options", db.ConstTypeJSON, true); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7516f60e-83f6-4b72-baa8-ae5c698f1a81", err.Error())
		}
		if err := collection.AddColumn("qty", db.ConstTypeInteger, false); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a1ce04d7-7a61-4318-a50f-5e3113b6183d", err.Error())
		}
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
