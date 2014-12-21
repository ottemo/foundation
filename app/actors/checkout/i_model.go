package checkout

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
)

// GetModelName returns model name we have implementation for
func (it *DefaultCheckout) GetModelName() string {
	return checkout.ConstCheckoutModelName
}

// GetImplementationName returns name of current model implementation
func (it *DefaultCheckout) GetImplementationName() string {
	return "Default" + checkout.ConstCheckoutModelName
}

// New makes new instance of model
func (it *DefaultCheckout) New() (models.InterfaceModel, error) {
	return &DefaultCheckout{Info: make(map[string]interface{})}, nil
}
