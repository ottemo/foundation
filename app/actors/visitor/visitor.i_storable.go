package visitor

import (
	"github.com/ottemo/foundation/app/actors/visitor/address"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetID returns current ID of the Visitor
func (it *DefaultVisitor) GetID() string {
	return it.id
}

// SetID sets current ID of the current Visitor
func (it *DefaultVisitor) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load will retrieve the Visitor information from database
func (it *DefaultVisitor) Load(ID string) error {

	collection, err := db.GetCollection(ConstCollectionNameVisitor)
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

// Delete removes current Visitor from the database
func (it *DefaultVisitor) Delete() error {

	collection, err := db.GetCollection(ConstCollectionNameVisitor)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	addressCollection, err := db.GetCollection(address.ConstCollectionNameVisitorAddress)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	addressCollection.AddFilter("visitor_id", "=", it.GetID())
	if _, err := addressCollection.Delete(); err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Save stores current Visitor to the database
func (it *DefaultVisitor) Save() error {

	collection, err := db.GetCollection(ConstCollectionNameVisitor)
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
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "29be1531cb6b44cfa78ef1bf9aae1163", "email already exists")
		}
	}

	storableValues := it.ToHashMap()

	delete(storableValues, "billing_address")
	delete(storableValues, "shipping_address")

	/*if it.Password == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "56a1aaaa897b445c853627cd95ecc705", "password can't be blank")
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
