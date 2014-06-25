package defaultProduct

import (
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/attribute"
)

func (it *ProductModel) GetModelName() string {
	return "Product"
}

func (it *ProductModel) GetImplementationName() string {
	return "DefaultProduct"
}

func (it *ProductModel) New() (models.IModel, error) {

	customAttributes, err := new(attribute.CustomAttributes).Init("product")
	if err != nil {
		return nil, err
	}

	return &ProductModel{CustomAttributes: customAttributes}, nil
}
