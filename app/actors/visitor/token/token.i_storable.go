package token

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
)

// GetID returns the Default Visitor Token as a string
func (it *DefaultVisitorCard) GetID() string {
	return it.id
}

// SetID takes a string as input and sets the ID on the Visitor Token
func (it *DefaultVisitorCard) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load will take Visitor Token ID and retrieve it from the database
func (it *DefaultVisitorCard) Load(loadID string) error {

	collection, err := db.GetCollection(ConstCollectionNameVisitorToken)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := collection.LoadByID(loadID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return it.FromHashMap(dbRecord)
}

// Delete will remove the Visitor Token from the database
func (it *DefaultVisitorCard) Delete() error {

	collection, err := db.GetCollection(ConstCollectionNameVisitorToken)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return collection.DeleteByID(it.GetID())
}

// Save will persist the Visitor Token to the database
func (it *DefaultVisitorCard) Save() error {

	collection, err := db.GetCollection(ConstCollectionNameVisitorToken)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if it.tokenID == "" || it.Payment == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dbd526dd-e0ae-4f43-a8cd-2da8ff3bba8a", "payment and token should be specified")
	}

	if utils.IsZeroTime(it.CreatedAt) && it.GetID() == "" {
		it.CreatedAt = time.Now()
	}

	storableValues := it.ToHashMap()
	storableValues["token_id"] = it.tokenID

	newID, err := collection.Save(storableValues)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return it.SetID(newID)
}
