package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns the Visitor address
func (it *DefaultVisitorAddressCollection) GetModelName() string {
	return visitor.ModelNameVisitorAddress
}

// GetImplementationName returns the default Visitor address
func (it *DefaultVisitorAddressCollection) GetImplementationName() string {
	return "Default" + visitor.ModelNameVisitorAddress
}

// New creates a new address for a Visitor
func (it *DefaultVisitorAddressCollection) New() (models.I_Model, error) {
	dbCollection, err := db.GetCollection(CollectionNameVisitorAddress)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitorAddressCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
