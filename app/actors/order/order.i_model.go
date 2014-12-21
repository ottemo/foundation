package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
)

// GetModelName returns model name we have implementation for
func (it *DefaultOrder) GetModelName() string {
	return order.ConstModelNameOrder
}

// GetImplementationName returns name of current model implementation
func (it *DefaultOrder) GetImplementationName() string {
	return "Default" + order.ConstModelNameOrder
}

// New makes new instance of model
func (it *DefaultOrder) New() (models.InterfaceModel, error) {
	return &DefaultOrder{Items: make(map[int]order.InterfaceOrderItem)}, nil
}
