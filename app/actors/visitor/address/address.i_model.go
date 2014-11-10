package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetModelName returns the Visitor address
func (it *DefaultVisitorAddress) GetModelName() string {
	return visitor.ModelNameVisitorAddress
}

// GetImplementationName returns the default Visitor address
func (it *DefaultVisitorAddress) GetImplementationName() string {
	return "Default" + visitor.ModelNameVisitorAddress
}

// New will create a new Visitor address or return an error
func (it *DefaultVisitorAddress) New() (models.I_Model, error) {
	return &DefaultVisitorAddress{}, nil
}
