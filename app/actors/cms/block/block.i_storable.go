package block

import (
	"time"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetID returns id for cms block
func (it *DefaultCMSBlock) GetID() string {
	return it.id
}

// SetID sets id for cms block
func (it *DefaultCMSBlock) SetID(newID string) error {
	it.id = newID
	return nil
}

// Load loads cms block information from DB
func (it *DefaultCMSBlock) Load(id string) error {
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbValues, err := collection.LoadByID(id)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := it.SetID(utils.InterfaceToString(dbValues["_id"])); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5c2adf9c-070c-4332-b9ef-4c5592d9fe5d", err.Error())
	}

	it.Content = utils.InterfaceToString(dbValues["content"])
	it.Identifier = utils.InterfaceToString(dbValues["identifier"])
	it.CreatedAt = utils.InterfaceToTime(dbValues["created_at"])
	it.UpdatedAt = utils.InterfaceToTime(dbValues["updated_at"])

	return nil
}

// Delete removes current cms block from DB
func (it *DefaultCMSBlock) Delete() error {
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return env.ErrorDispatch(err)
}

// Save stores current cms block to DB
func (it *DefaultCMSBlock) Save() error {
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// packing data before save
	storingValues := make(map[string]interface{})

	storingValues["_id"] = it.GetID()

	storingValues["identifier"] = it.GetIdentifier()
	storingValues["content"] = it.GetContent()

	currentTime := time.Now()

	if it.CreatedAt.IsZero() {
		it.CreatedAt = currentTime
	}
	storingValues["created_at"] = it.CreatedAt

	it.UpdatedAt = currentTime
	storingValues["updated_at"] = it.UpdatedAt

	newID, err := collection.Save(storingValues)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	if err := it.SetID(newID); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2d71c28e-d8b6-4ce0-95b9-73b724bc5d89", err.Error())
	}

	return nil
}
