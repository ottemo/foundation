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
		if err := it.Set(attribute, value); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e4770f57-e1f9-4b73-a23a-7b2627c659ee", err.Error())
		}
	}

	it.Items = make(map[int]order.InterfaceOrderItem)
	it.maxIdx = 0

	// loading order items
	orderItemsCollection, err := db.GetCollection(ConstCollectionNameOrderItems)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := orderItemsCollection.AddFilter("order_id", "=", it.GetID()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5b4bd653-2b14-4e6a-ba2f-4f3fc083a608", err.Error())
	}
	orderItems, err := orderItemsCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, orderItemValues := range orderItems {
		orderItem := new(DefaultOrderItem)

		for attribute, value := range orderItemValues {
			if err := orderItem.Set(attribute, value); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8047073d-164b-4a5c-bb29-d076c8bd3065", err.Error())
			}
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
	if err := it.SetID(newID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9e1abe0d-20a0-4bea-8f81-40fddaacde3f", err.Error())
	}

	// storing order items
	for _, orderItem := range it.GetItems() {
		if err := orderItem.Set("order_id", newID); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "90202b42-0cd0-4e08-a11a-6f1f5981d246", err.Error())
		}
		orderItemStoringValues := orderItem.ToHashMap()

		newID, err := orderItemsCollection.Save(orderItemStoringValues)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		if err := orderItem.SetID(newID); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "95eb6404-df35-4386-8a74-268df577194d", err.Error())
		}
	}

	return nil
}
