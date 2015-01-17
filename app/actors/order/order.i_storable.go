package order

import (
	"time"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetID returns id of current order
func (it *DefaultOrder) GetID() string {
	return it.id
}

// SetID sets id for order
func (it *DefaultOrder) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load loads order information from DB
func (it *DefaultOrder) Load(ID string) error {

	// loading order
	orderCollection, err := db.GetCollection(ConstCollectionNameOrder)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	values, err := orderCollection.LoadByID(ID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// initializing DefaultOrder structure
	for attribute, value := range values {
		it.Set(attribute, value)
	}

	it.Items = make(map[int]order.InterfaceOrderItem)
	it.maxIdx = 0

	// loading order items
	orderItemsCollection, err := db.GetCollection(ConstCollectionNameOrderItems)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	orderItemsCollection.AddFilter("order_id", "=", it.GetID())
	orderItems, err := orderItemsCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, orderItemValues := range orderItems {
		orderItem := new(DefaultOrderItem)

		for attribute, value := range orderItemValues {
			orderItem.Set(attribute, value)
		}

		it.Items[orderItem.idx] = orderItem
	}

	return nil
}

// Delete removes current order from DB
func (it *DefaultOrder) Delete() error {
	if it.GetID() == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "05e07871-2c78-4852-ab06-505bf8d708c1", "order id is not set")
	}

	// deleting order items
	orderItemsCollection, err := db.GetCollection(ConstCollectionNameOrderItems)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = orderItemsCollection.AddFilter("order_id", "=", it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	_, err = orderItemsCollection.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// deleting order
	orderCollection, err := db.GetCollection(ConstCollectionNameOrder)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = orderCollection.DeleteByID(it.GetID())

	return env.ErrorDispatch(err)
}

// Save stores current order in DB
func (it *DefaultOrder) Save() error {

	orderCollection, err := db.GetCollection(ConstCollectionNameOrder)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	orderItemsCollection, err := db.GetCollection(ConstCollectionNameOrderItems)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// packing data before save
	orderStoringValues := it.ToHashMap()

	it.UpdatedAt = time.Now()

	newID, err := orderCollection.Save(orderStoringValues)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	it.SetID(newID)

	// storing order items
	for _, orderItem := range it.GetItems() {
		orderItem.Set("order_id", newID)
		orderItemStoringValues := orderItem.ToHashMap()

		newID, err := orderItemsCollection.Save(orderItemStoringValues)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		orderItem.SetID(newID)
	}

	return nil
}
