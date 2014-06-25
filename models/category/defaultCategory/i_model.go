package defaultCategory

import (
	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/product"
)

func (it *DefaultCategory) GetModelName() string {
	return "Category"
}

func (it *DefaultCategory) GetImplementationName() string {
	return "DefaultCategory"
}

func (it *DefaultCategory) New() (models.IModel, error) {
	return &DefaultCategory{Products: make([]product.IProduct, 0)}, nil
}
