package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
)

// GetCollection returns collection of current instance type
func (it *DefaultOrder) GetCollection() models.InterfaceCollection {
	model, _ := models.GetModel(order.ConstModelNameOrderCollection)
	if result, ok := model.(order.InterfaceOrderCollection); ok {
		return result
	}

	return nil
}
