package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"time"
)

// GetOrderCollectionModel retrieves current InterfaceOrderCollection model implementation
func GetOrderCollectionModel() (InterfaceOrderCollection, error) {
	model, err := models.GetModel(ConstModelNameOrderCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderCollectionModel, ok := model.(InterfaceOrderCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1f016c45-174b-4bb6-bdf4-00f2ec942f64", "model "+model.GetImplementationName()+" is not 'InterfaceOrderCollection' capable")
	}

	return orderCollectionModel, nil
}

// GetOrderItemCollectionModel retrieves current InterfaceOrderCollection model implementation
func GetOrderItemCollectionModel() (InterfaceOrderItemCollection, error) {
	model, err := models.GetModel(ConstModelNameOrderItemCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderItemCollectionModel, ok := model.(InterfaceOrderItemCollection)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "54bb3316-68f9-4999-97a3-8549dde0940a", "model "+model.GetImplementationName()+" is not 'InterfaceOrderItemCollection' capable")
	}

	return orderItemCollectionModel, nil
}

// GetOrderModel retrieves current InterfaceOrder model implementation
func GetOrderModel() (InterfaceOrder, error) {
	model, err := models.GetModel(ConstModelNameOrder)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderModel, ok := model.(InterfaceOrder)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "58fe1c56-cea4-40f3-bf66-3c51c03838c0", "model "+model.GetImplementationName()+" is not 'InterfaceOrder' capable")
	}

	return orderModel, nil
}

// GetOrderModelAndSetID retrieves current InterfaceOrder model implementation and sets its ID to some value
func GetOrderModelAndSetID(orderID string) (InterfaceOrder, error) {

	orderModel, err := GetOrderModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = orderModel.SetID(orderID)
	if err != nil {
		return orderModel, env.ErrorDispatch(err)
	}

	return orderModel, nil
}

// LoadOrderByID loads order data into current InterfaceOrder model implementation
func LoadOrderByID(orderID string) (InterfaceOrder, error) {

	orderModel, err := GetOrderModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = orderModel.Load(orderID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return orderModel, nil
}

// GetOrdersCreatedBetween Get the orders `created_at` a certain date range
func GetOrdersCreatedBetween(startDate time.Time, endDate time.Time) []models.StructListItem {
	oModel, _ := GetOrderCollectionModel()
	oModel.GetDBCollection().AddFilter("created_at", ">=", startDate)
	oModel.GetDBCollection().AddFilter("created_at", "<", endDate)
	oModel.ListAddExtraAttribute("created_at") // If you are filtering on created_at you probably want that too
	foundOrders, _ := oModel.List()            // This is the lite response StructListItem

	return foundOrders
}

// GetFullOrdersUpdatedBetween db query for getting all orders, with expensive details
func GetFullOrdersUpdatedBetween(startDate time.Time, endDate time.Time) []InterfaceOrder {
	oModel, _ := GetOrderCollectionModel()
	oModel.GetDBCollection().AddFilter("updated_at", ">=", startDate)
	oModel.GetDBCollection().AddFilter("updated_at", "<", endDate)
	result := oModel.ListOrders()

	return result
}

// GetItemsForOrders Get the relavent order items given a slice of orders
func GetItemsForOrders(orderIds []string) []map[string]interface{} {
	oiModel, _ := GetOrderItemCollectionModel()
	oiDB := oiModel.GetDBCollection()
	oiDB.AddFilter("order_id", "in", orderIds)
	oiResults, _ := oiDB.Load()

	return oiResults
}
