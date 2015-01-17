package order

import (
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
