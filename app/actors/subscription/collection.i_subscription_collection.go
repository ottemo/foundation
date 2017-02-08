package subscription

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/subscription"
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
		if err := subscriptionModel.FromHashMap(recordData); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "da49add1-620a-40b3-940d-075c190f7c9a", err.Error())
		}

		result = append(result, subscriptionModel)
	}

	return result
}
