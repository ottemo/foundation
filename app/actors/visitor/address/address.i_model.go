package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

func (it *DefaultVisitorAddress) GetModelName() string {
	return visitor.MODEL_NAME_VISITOR_ADDRESS
}

func (it *DefaultVisitorAddress) GetImplementationName() string {
	return "Default" + visitor.MODEL_NAME_VISITOR_ADDRESS
}

func (it *DefaultVisitorAddress) New() (models.I_Model, error) {
	return &DefaultVisitorAddress{}, nil
}
