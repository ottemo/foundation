package token

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetModelName returns the Visitor Address Model
func (it *DefaultVisitorCard) GetModelName() string {
	return visitor.ConstModelNameVisitorCard
}

// GetImplementationName returns the Implementation name
func (it *DefaultVisitorCard) GetImplementationName() string {
	return "Default" + visitor.ConstModelNameVisitorCard
}

// New creates a new Visitor Address interface
func (it *DefaultVisitorCard) New() (models.InterfaceModel, error) {
	return &DefaultVisitorCard{}, nil
}
