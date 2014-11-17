package visitor

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// returns model name
func (it *DefaultVisitorCollection) GetModelName() string {
	return visitor.MODEL_NAME_VISITOR_COLLECTION
}

// returns model implementation name
func (it *DefaultVisitorCollection) GetImplementationName() string {
	return "Default" + visitor.MODEL_NAME_VISITOR_COLLECTION
}

// returns new instance of model implementation object
func (it *DefaultVisitorCollection) New() (models.I_Model, error) {
	dbCollection, err := db.GetCollection(COLLECTION_NAME_VISITOR)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitorCollection{listCollection: dbCollection, listExtraAtributes: make([]string, 0)}, nil
}
