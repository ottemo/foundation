package visitor

import (
	"github.com/ottemo/foundation/app/helpers/attributes"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns model name
func (it *DefaultVisitor) GetModelName() string {
	return visitor.ModelNameVisitor
}

// GetImplementationName returns model implementation name
func (it *DefaultVisitor) GetImplementationName() string {
	return "Default" + visitor.ModelNameVisitor
}

// New returns new instance of model implementation object
func (it *DefaultVisitor) New() (models.I_Model, error) {

	customAttributes, err := new(attributes.CustomAttributes).Init(visitor.ModelNameVisitor, CollectionNameVisitor)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultVisitor{CustomAttributes: customAttributes}, nil
}
