package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetModelName returns the Visitor Address Model
func (it *DefaultVisitorAddress) GetModelName() string {
	return visitor.ConstModelNameVisitorAddress
}

// GetImplementationName returns the Implementation name
func (it *DefaultVisitorAddress) GetImplementationName() string {
	return "Default" + visitor.ConstModelNameVisitorAddress
}

// New creates a new Visitor Address interface
func (it *DefaultVisitorAddress) New() (models.InterfaceModel, error) {
	return &DefaultVisitorAddress{}, nil
}
