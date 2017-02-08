package order

import (
	"strings"
	"time"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
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
	oModel, err := GetOrderCollectionModel()
	if err != nil {
		_ = env.ErrorDispatch(err)
	}
	if err := oModel.GetDBCollection().AddFilter("created_at", ">=", startDate); err != nil {
		_ = env.ErrorDispatch(err)
	}
	if err := oModel.GetDBCollection().AddFilter("created_at", "<", endDate); err != nil {
		_ = env.ErrorDispatch(err)
	}
	if err := oModel.ListAddExtraAttribute("created_at"); err != nil { // If you are filtering on created_at you probably want that too
		_ = env.ErrorDispatch(err)
	}
	foundOrders, err := oModel.List()            // This is the lite response StructListItem
	if err != nil {
		_ = env.ErrorDispatch(err)
	}

	return foundOrders
}

// GetFullOrdersUpdatedBetween db query for getting all orders, with expensive details
func GetFullOrdersUpdatedBetween(startDate time.Time, endDate time.Time) []InterfaceOrder {
	oModel, err := GetOrderCollectionModel()
	if err != nil {
		_ = env.ErrorDispatch(err)
	}
	if err := oModel.GetDBCollection().AddFilter("updated_at", ">=", startDate); err != nil {
		_ = env.ErrorDispatch(err)
	}
	if err := oModel.GetDBCollection().AddFilter("updated_at", "<", endDate); err != nil {
		_ = env.ErrorDispatch(err)
	}
	result := oModel.ListOrders()

	return result
}

// GetItemsForOrders Get the relavent order items given a slice of orders
func GetItemsForOrders(orderIds []string) []map[string]interface{} {
	oiModel, err := GetOrderItemCollectionModel()
	if err != nil {
		_ = env.ErrorDispatch(err)
	}
	oiDB := oiModel.GetDBCollection()
	if err := oiDB.AddFilter("order_id", "in", orderIds); err != nil {
		_ = env.ErrorDispatch(err)
	}
	oiResults, err := oiDB.Load()
	if err != nil {
		_ = env.ErrorDispatch(err)
	}

	return oiResults
}

// SplitFullName will take a fullname as a string and split it into first name and last names
func SplitFullName(name string) (string, string) {

	var firstName, lastName string

	fullName := strings.SplitN(name, " ", 2)

	if len(fullName) == 2 {
		firstName = fullName[0]
		lastName = fullName[1]
	} else if len(fullName) == 1 {
		firstName = fullName[0]
		lastName = ""
	} else {
		firstName = ""
		lastName = ""
	}

	return firstName, lastName
}
