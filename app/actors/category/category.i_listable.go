package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
)

// GetCollection returns collection of current instance type
func (it *DefaultCategory) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(category.ConstModelNameCategoryCollection)
	if result, ok := model.(category.InterfaceCategoryCollection); ok {
		return result
	}

	return nil
}
