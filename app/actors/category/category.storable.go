package category

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

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

	junctionCollection.AddFilter("category_id", "=", it.GetID())
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

	delete(storingValues, "products")

	categoryCollection, err := db.GetCollection(ConstCollectionNameCategory)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// saving category
	if newID, err := categoryCollection.Save(storingValues); err == nil {
		if it.GetID() != newID {
			it.SetID(newID)
			it.updatePath()
			it.Save()
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
	junctionCollection.AddFilter("category_id", "=", it.GetID())
	_, err = junctionCollection.Delete()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// adding new assignments
	for _, categoryProductID := range it.ProductIds {
		junctionCollection.Save(map[string]interface{}{"category_id": it.GetID(), "product_id": categoryProductID})
	}

	return nil
}
