package defaultproduct

import (
	"github.com/ottemo/foundation/models"
)

func (it *DefaultProductModel) GetModelName() string {
	return "Product"
}

func (it *DefaultProductModel) GetImplementationName() string {
	return "DefaultProduct"
}

func (it *DefaultProductModel) New() (models.Model, error) {

	customAttributes, err := new(custom_attributes.CustomAttributes).Init("product")
	if err != nil {
		return nil, err
	}

	return &DefaultProductModel{CustomAttributes: customAttributes}, nil
}
