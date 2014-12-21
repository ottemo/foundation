package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns model name
func (it *DefaultOrderItemCollection) GetModelName() string {
	return order.ConstModelNameOrderItemCollection
}

// GetImplementationName returns model implementation name
func (it *DefaultOrderItemCollection) GetImplementationName() string {
	return "Default" + order.ConstModelNameOrderItemCollection
}

// New returns new instance of model implementation object
func (it *DefaultOrderItemCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameOrderItems)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultOrderItemCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
