package defaultCategory

import (
	"github.com/ottemo/foundation/models/product"
)

func (it *DefaultCategory) GetName() string { return it.Name }
func (it *DefaultCategory) GetProducts() []product.IProduct {
	return it.Products
}
