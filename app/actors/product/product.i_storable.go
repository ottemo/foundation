package product

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// returns current product id
func (it *DefaultProduct) GetId() string {
	return it.id
}

// sets current product id
func (it *DefaultProduct) SetId(NewId string) error {
	it.id = NewId
	return nil
}

// loads product information from DB
func (it *DefaultProduct) Load(loadId string) error {

	collection, err := db.GetCollection(COLLECTION_NAME_PRODUCT)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := collection.LoadById(loadId)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// removes current product from DB
func (it *DefaultProduct) Delete() error {
	collection, err := db.GetCollection(COLLECTION_NAME_PRODUCT)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteById(it.GetId())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// stores current product to DB
func (it *DefaultProduct) Save() error {
	collection, err := db.GetCollection(COLLECTION_NAME_PRODUCT)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	newId, err := collection.Save(it.ToHashMap())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.SetId(newId)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
