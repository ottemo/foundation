package stripesubscription

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetModelName returns model name for the Subscription Collection
func (it *DefaultStripeSubscriptionCollection) GetModelName() string {
	return ConstModelNameStripeSubscriptionCollection
}

// GetImplementationName returns model implementation name for the Subscription Collection
func (it *DefaultStripeSubscriptionCollection) GetImplementationName() string {
	return "Default" + ConstModelNameStripeSubscriptionCollection
}

// New returns new instance of model implementation object for the Subscription Collection
func (it *DefaultStripeSubscriptionCollection) New() (models.InterfaceModel, error) {
	dbCollection, err := db.GetCollection(ConstCollectionNameStripeSubscription)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return &DefaultStripeSubscriptionCollection{listCollection: dbCollection, listExtraAttributes: make([]string, 0)}, nil
}
