package visitor

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns model name for the Visitor
func (it *DefaultVisitor) GetModelName() string {
	return visitor.ConstModelNameVisitor
}

// GetImplementationName returns model implementation name for the Visitor
func (it *DefaultVisitor) GetImplementationName() string {
	return "Default" + visitor.ConstModelNameVisitor
}

// New returns new instance of model implementation object for the Visitor
func (it *DefaultVisitor) New() (models.InterfaceModel, error) {

	customAttributes, err := attributes.CustomAttributes(visitor.ConstModelNameVisitor, ConstCollectionNameVisitor)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitor{ModelCustomAttributes: customAttributes}, nil
}
