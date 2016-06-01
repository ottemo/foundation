package product

import (
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetID returns current product id
func (it *DefaultProduct) GetID() string {
	return it.id
}

// SetID sets current product id
func (it *DefaultProduct) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load loads product information from DB
func (it *DefaultProduct) Load(loadID string) error {

	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := collection.LoadByID(loadID)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a671dee4-b95b-11e5-a86b-28cfe917b6c7", "Unable to find product by id; "+loadID)
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Delete removes current product from DB
func (it *DefaultProduct) Delete() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// stock management stuff
	if stockManager := product.GetRegisteredStock(); stockManager != nil {
		stockManager.RemoveProductQty(it.GetID(), make(map[string]interface{}))
	}

	return nil
}

// Save stores current product to DB
func (it *DefaultProduct) Save() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if it.GetName() == "" || it.GetSku() == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ac7cd02e-0722-4ac8-bbe0-ffa74d091a94", "sku and name should be specified")
	}

	valuesToStore := it.ToHashMap()
	if _, present := valuesToStore["qty"]; present {
		delete(valuesToStore, "qty")
	}

	newID, err := collection.Save(valuesToStore)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.SetID(newID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// stock management stuff
	if stockManager := product.GetRegisteredStock(); stockManager != nil {
		for _, qtyOptions := range it.updatedQty {
			if qtyOptions == nil {
				continue
			}

			if qtyValue, present := qtyOptions[""]; present {
				qty := utils.InterfaceToInt(qtyValue)
				delete(qtyOptions, "")

				err := stockManager.SetProductQty(it.GetID(), qtyOptions, qty)
				if err != nil {
					return env.ErrorDispatch(err)
				}
			}
		}
	}

	return nil
}
