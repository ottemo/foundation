package stripesubscription

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
)

// GetID returns current ID of the Stripe Subscription
func (it *DefaultStripeSubscription) GetID() string {
	return it.id
}

// SetID sets current ID of the current Stripe Subscription
func (it *DefaultStripeSubscription) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load will retrieve the Stripe Subscription information from database
func (it *DefaultStripeSubscription) Load(ID string) error {

	collection, err := db.GetCollection(ConstCollectionNameStripeSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	values, err := collection.LoadByID(ID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.FromHashMap(values)
	return env.ErrorDispatch(err)
}

// Delete removes current Stripe Subscription from the database
func (it *DefaultStripeSubscription) Delete() error {

	collection, err := db.GetCollection(ConstCollectionNameStripeSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(it.GetID())
	return env.ErrorDispatch(err)
}

// Save stores current Stripe Subscription to the database
func (it *DefaultStripeSubscription) Save() error {

	collection, err := db.GetCollection(ConstCollectionNameStripeSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// update time depend values
	currentTime := time.Now()
	if it.CreatedAt.IsZero() {
		it.CreatedAt = currentTime
	}
	it.UpdatedAt = currentTime

	storingValues := it.ToHashMap()

	// saving operation
	newID, err := collection.Save(storingValues)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return it.SetID(newID)
}
