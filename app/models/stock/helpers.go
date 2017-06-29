package stock

import (
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
)

// GetBlogPostModel retrieves current InterfaceBlogPost model implementation
func GetStockModel() (InterfaceStock, error) {
	model, err := models.GetModel(ConstModelNameStock)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stockModel, ok := model.(InterfaceStock)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d2fd8dc4-094e-4a32-8e89-d8be2f3bf1ea", "model "+model.GetImplementationName()+" is not 'InterfaceStock' capable")
	}

	return stockModel, nil
}

// GetProductCollectionModel retrieves current InterfaceProductCollection model implementation
func GetStockCollectionModel() (InterfaceStockCollection, error) {
	model, err := models.GetModel(ConstModelNameStockCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	stockModel, ok := model.(InterfaceStockCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "343218b7-587f-47e5-83a8-a372615116d9", "model "+model.GetImplementationName()+" is not 'InterfaceStockCollection' capable")
	}

	return stockModel, nil
}
