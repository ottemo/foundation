package product

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns model name
func (it *DefaultProduct) GetModelName() string {
	return product.ConstModelNameProduct
}

// GetImplementationName returns model implementation name
func (it *DefaultProduct) GetImplementationName() string {
	return "Default" + product.ConstModelNameProduct
}

// New returns new instance of model implementation object
func (it *DefaultProduct) New() (models.InterfaceModel, error) {

	customAttributes, err := new(attributes.CustomAttributes).Init(product.ConstModelNameProduct, ConstCollectionNameProduct)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultProduct{CustomAttributes: customAttributes}, nil
}
