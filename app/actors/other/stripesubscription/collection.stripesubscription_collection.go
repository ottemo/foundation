package stripesubscription

import (
	"github.com/ottemo/foundation/app/models/stripesubscription"
	"github.com/ottemo/foundation/db"
)

// GetDBCollection returns database collection for the Stripe Subscription
func (it *DefaultStripeSubscriptionCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListSubscriptions returns list of subscription model items in the Subscription Collection
func (it *DefaultStripeSubscriptionCollection) ListSubscriptions() []stripesubscription.InterfaceStripeSubscription {
	var result []stripesubscription.InterfaceStripeSubscription

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		stripeSubscriptionModel, err := stripesubscription.GetStripeSubscriptionModel()
		if err != nil {
			return result
		}
		stripeSubscriptionModel.FromHashMap(recordData)

		result = append(result, stripeSubscriptionModel)
	}

	return result
}
