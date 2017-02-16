package order

import (
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetDBCollection returns database collection
func (it *DefaultOrderCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListOrders returns array of products in model instance form
func (it *DefaultOrderCollection) ListOrders() []order.InterfaceOrder {
	var result []order.InterfaceOrder

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, dbRecordData := range dbRecords {
		orderModel, err := order.GetOrderModel()
		if err != nil {
			return result
		}
		if err := orderModel.FromHashMap(dbRecordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "29635564-a788-46d7-9485-f286ae481533", err.Error())
		}

		result = append(result, orderModel)
	}

	return result
}
