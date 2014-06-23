package defaultproduct

import (
	"github.com/ottemo/foundation/models"
)

func (dpm *DefaultProductModel) GetModelName() string {
	return "Product"
}

func (dpm *DefaultProductModel) GetImplementationName() string {
	return "DefaultProduct"
}

func (dpm *DefaultProductModel) New() (models.Model, error) {

	customAttributes, err := new(models.CustomAttribute).Init("product")
	if err != nil {
		return nil, err
	}

	return &DefaultProductModel{CustomAttribute: customAttributes}, nil
}
