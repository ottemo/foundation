package address

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetID will retrieve the Visitor address ID
func (it *DefaultVisitorAddress) GetID() string {
	return it.id
}

// SetID will set the Visitor address ID
func (it *DefaultVisitorAddress) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load will retrieve the Visitor address from the db
func (it *DefaultVisitorAddress) Load(ID string) error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(CollectionNameVisitorAddress); err == nil {

			if values, err := collection.LoadById(ID); err == nil {
				if err := it.FromHashMap(values); err != nil {
					return env.ErrorDispatch(err)
				}
			} else {
				return env.ErrorDispatch(err)
			}

		} else {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

// Delete will remove the Visitor address from the db
func (it *DefaultVisitorAddress) Delete() error {
	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(CollectionNameVisitorAddress); err == nil {
			err := collection.DeleteById(it.GetID())
			if err != nil {
				return env.ErrorDispatch(err)
			}
		} else {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

// Save will persist the Visitor address to the db
func (it *DefaultVisitorAddress) Save() error {

	if dbEngine := db.GetDBEngine(); dbEngine != nil {
		if collection, err := dbEngine.GetCollection(CollectionNameVisitorAddress); err == nil {

			//if it.ZipCode== "" {
			//	return env.ErrorNew("Zip code for address - required")
			//}

			if newID, err := collection.Save(it.ToHashMap()); err == nil {
				it.Set("_id", newID)
				return env.ErrorDispatch(err)
			}

			return env.ErrorDispatch(err)

		}

		// return env.ErrorDispatch(err)
	}

	return nil
}
