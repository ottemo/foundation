package visitor

import (
	"github.com/ottemo/foundation/app/actors/visitor/address"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetID returns current product id
func (it *DefaultVisitor) GetID() string {
	return it.id
}

// SetID sets current product id
func (it *DefaultVisitor) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load will request the Visitor information from DB
func (it *DefaultVisitor) Load(ID string) error {

	collection, err := db.GetCollection(CollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	values, err := collection.LoadByID(ID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.FromHashMap(values)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Delete removes current visitor from DB
func (it *DefaultVisitor) Delete() error {

	collection, err := db.GetCollection(CollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	addressCollection, err := db.GetCollection(address.CollectionNameVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	addressCollection.AddFilter("visitorID", "=", it.GetID())
	if _, err := addressCollection.Delete(); err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteById(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Save stores current visitor to DB
func (it *DefaultVisitor) Save() error {

	collection, err := db.GetCollection(CollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if it.GetID() == "" {
		collection.AddFilter("email", "=", it.GetEmail())
		n, err := collection.Count()
		if err != nil {
			return env.ErrorDispatch(err)
		}
		if n > 0 {
			return env.ErrorNew("email already exists")
		}
	}

	storableValues := it.ToHashMap()

	delete(storableValues, "billing_address")
	delete(storableValues, "shipping_address")

	/*if it.Password == "" {
		return env.ErrorNew("password can't be blank")
	}*/

	storableValues["facebook_id"] = it.FacebookID
	storableValues["google_id"] = it.GoogleID
	storableValues["password"] = it.Password
	storableValues["validate"] = it.ValidateKey

	// shipping address save
	if it.ShippingAddress != nil {
		err := it.ShippingAddress.Save()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		storableValues["shipping_address_id"] = it.ShippingAddress.GetID()
	}

	// billing address save
	if it.BillingAddress != nil {
		err := it.BillingAddress.Save()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		storableValues["billing_address_id"] = it.BillingAddress.GetID()
	}

	// saving visitor
	if newID, err := collection.Save(storableValues); err == nil {
		it.Set("_id", newID)
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
