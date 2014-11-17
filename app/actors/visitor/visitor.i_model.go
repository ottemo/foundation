package visitor

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
)

// returns model name
func (it *DefaultVisitor) GetModelName() string {
	return visitor.MODEL_NAME_VISITOR
}

// returns model implementation name
func (it *DefaultVisitor) GetImplementationName() string {
	return "Default" + visitor.MODEL_NAME_VISITOR
}

// returns new instance of model implementation object
func (it *DefaultVisitor) New() (models.I_Model, error) {

	customAttributes, err := new(attributes.CustomAttributes).Init(visitor.MODEL_NAME_VISITOR, COLLECTION_NAME_VISITOR)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitor{CustomAttributes: customAttributes}, nil
}
