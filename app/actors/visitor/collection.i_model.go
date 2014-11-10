package visitor

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns model name
func (it *DefaultVisitorCollection) GetModelName() string {
	return visitor.ModelNameVisitorCollection
}

// GetImplementationName returns model implementation name
func (it *DefaultVisitorCollection) GetImplementationName() string {
	return "Default" + visitor.ModelNameVisitorCollection
}

// New returns new instance of model implementation object
func (it *DefaultVisitorCollection) New() (models.I_Model, error) {
	dbCollection, err := db.GetCollection(CollectionNameVisitor)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitorCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
