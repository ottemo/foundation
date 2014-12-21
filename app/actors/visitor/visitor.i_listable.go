package visitor

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
)

// GetCollection returns collection of current instance type
func (it *DefaultVisitor) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(visitor.ConstModelNameVisitorCollection)
	if result, ok := model.(visitor.InterfaceVisitorCollection); ok {
		return result
	}

	return nil
}
