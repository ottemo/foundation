package address

import (
	"github.com/ottemo/foundation/models"
)

func (it *DefaultVisitorAddress) GetModelName() string {
	return "VisitorAddress"
}

func (it *DefaultVisitorAddress) GetImplementationName() string {
	return "DefaultVisitorAddress"
}

func (it *DefaultVisitorAddress) New() (models.IModel, error) {
	return &DefaultVisitorAddress{}, nil
}
