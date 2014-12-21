package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
)

// GetCollection returns collection of current instance type
func (it *DefaultProduct) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(product.ConstModelNameProductCollection)
	if result, ok := model.(product.InterfaceProductCollection); ok {
		return result
	}

	return nil
}
