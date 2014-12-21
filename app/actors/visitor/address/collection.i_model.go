package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns the Visitor Address model
func (it *DefaultVisitorAddressCollection) GetModelName() string {
	return visitor.ConstModelNameVisitorAddress
}

// GetImplementationName returns the Visitor Address implementation name
func (it *DefaultVisitorAddressCollection) GetImplementationName() string {
	return "Default" + visitor.ConstModelNameVisitorAddress
}

// New creates a new Visitor Address Collection
func (it *DefaultVisitorAddressCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameVisitorAddress)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitorAddressCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
