package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns model name
func (it *DefaultOrderCollection) GetModelName() string {
	return order.ConstModelNameOrderCollection
}

// GetImplementationName returns model implementation name
func (it *DefaultOrderCollection) GetImplementationName() string {
	return "Default" + order.ConstModelNameOrderCollection
}

// New returns new instance of model implementation object
func (it *DefaultOrderCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameOrder)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultOrderCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
