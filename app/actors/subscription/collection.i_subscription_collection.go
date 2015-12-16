package subscription

import (
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/db"
)

// GetDBCollection returns database collection for the Subscription
func (it *DefaultSubscriptionCollection) GetDBCollection() db.InterfaceDBCollection {
	return it.listCollection
}

// ListSubscriptions returns list of subscription model items in the Subscription Collection
func (it *DefaultSubscriptionCollection) ListSubscriptions() []subscription.InterfaceSubscription {
	var result []subscription.InterfaceSubscription

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result
	}

	for _, recordData := range dbRecords {
		subscriptionModel, err := subscription.GetSubscriptionModel()
		if err != nil {
			return result
		}
		subscriptionModel.FromHashMap(recordData)

		result = append(result, subscriptionModel)
	}

	return result
}
