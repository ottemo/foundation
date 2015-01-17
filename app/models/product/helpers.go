package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// GetProductCollectionModel retrieves current InterfaceProductCollection model implementation
func GetProductCollectionModel() (InterfaceProductCollection, error) {
	model, err := models.GetModel(ConstModelNameProductCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	productModel, ok := model.(InterfaceProductCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "343218b7-587f-47e5-83a8-a372615116d9", "model "+model.GetImplementationName()+" is not 'InterfaceProductCollection' capable")
	}

	return productModel, nil
}

// GetProductModel retrieves current InterfaceProduct model implementation
func GetProductModel() (InterfaceProduct, error) {
	model, err := models.GetModel(ConstModelNameProduct)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	productModel, ok := model.(InterfaceProduct)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a21db343-9d5d-4ae3-819d-29cf32154eb6", "model "+model.GetImplementationName()+" is not 'InterfaceProduct' capable")
	}

	return productModel, nil
}

// GetProductModelAndSetID retrieves current InterfaceProduct model implementation and sets its ID to some value
func GetProductModelAndSetID(productID string) (InterfaceProduct, error) {

	productModel, err := GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.SetID(productID)
	if err != nil {
		return productModel, env.ErrorDispatch(err)
	}

	return productModel, nil
}

// LoadProductByID loads product data into current InterfaceProduct model implementation
func LoadProductByID(productID string) (InterfaceProduct, error) {

	productModel, err := GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.Load(productID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return productModel, nil
}
