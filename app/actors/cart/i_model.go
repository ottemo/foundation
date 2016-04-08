package cart

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
)

// GetModelName returns model name we have implementation for
func (it *DefaultCart) GetModelName() string {
	return cart.ConstCartModelName
}

// GetImplementationName returns name of current model implementation
func (it *DefaultCart) GetImplementationName() string {
	return "DefaultCart"
}

// New makes new instance of model
func (it *DefaultCart) New() (models.InterfaceModel, error) {
	return &DefaultCart{
		Items:      make(map[int]cart.InterfaceCartItem),
		Info:       make(map[string]interface{}),
		CustomInfo: make(map[string]interface{}),
	}, nil
}
