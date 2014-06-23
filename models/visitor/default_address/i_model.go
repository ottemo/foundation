package default_address

import (
	"github.com/ottemo/foundation/models"
)

func (dva *DefaultVisitorAddress) GetModelName() string {
	return "VisitorAddress"
}

func (dva *DefaultVisitorAddress) GetImplementationName() string {
	return "DefaultVisitorAddress"
}

func (dva *DefaultVisitorAddress) New() (models.Model, error) {
	return &DefaultVisitorAddress{}, nil
}
