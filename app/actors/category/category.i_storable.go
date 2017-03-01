package category

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetID returns database storage id of current object
func (it *DefaultCategory) GetID() string {
	return it.id
}

// SetID sets database storage id for current object
func (it *DefaultCategory) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load loads object information from database storage
func (it *DefaultCategory) Load(ID string) error {

	// loading category
	categoryCollection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := categoryCollection.LoadByID(ID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.updatePath()

	// loading category product ids
	junctionCollection, err := db.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := junctionCollection.AddFilter("category_id", "=", it.GetID()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c176c200-bc35-41ea-81d6-453f1e8c01bb", err.Error())
	}
	junctedProducts, err := junctionCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, junctionRecord := range junctedProducts {
		it.ProductIds = append(it.ProductIds, utils.InterfaceToString(junctionRecord["product_id"]))
	}

	return nil
}

// Delete removes current object from database storage
func (it *DefaultCategory) Delete() error {
	//deleting category products join
	junctionCollection, err := db.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = junctionCollection.AddFilter("category_id", "=", it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	_, err = junctionCollection.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// deleting category
	categoryCollection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = categoryCollection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Save stores current object to database storage
func (it *DefaultCategory) Save() error {

	storingValues := it.ToHashMap()

	delete(storingValues, "product_ids")

	categoryCollection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// saving category
	if newID, err := categoryCollection.Save(storingValues); err == nil {
		if it.GetID() != newID {
			if err := it.SetID(newID); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ddd740b0-f866-4901-a7b9-f8e86b34191e", err.Error())
			}
			it.updatePath()
			if err := it.Save(); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f8fff0a3-9c6a-4528-84b3-11fc9ef10d4c", err.Error())
			}
		}
	} else {
		return env.ErrorDispatch(err)
	}

	// saving category products assignment
	junctionCollection, err := db.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// deleting old assigned products
	if err := junctionCollection.AddFilter("category_id", "=", it.GetID()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d59a3a77-c72d-479d-8ed4-e84332dd89be", err.Error())
	}
	_, err = junctionCollection.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// adding new assignments
	for _, categoryProductID := range it.ProductIds {
		if _, err := junctionCollection.Save(map[string]interface{}{"category_id": it.GetID(), "product_id": categoryProductID}); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8548775c-1b1b-41d4-92ac-643bd186dc65", err.Error())
		}
	}

	return nil
}
