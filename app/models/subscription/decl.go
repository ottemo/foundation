package subscription

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstModelNameSubscription           = "Subscription"
	ConstModelNameSubscriptionCollection = "SubscriptionCollection"

	ConstErrorModule = "subscription"
	ConstErrorLevel  = env.ConstErrorLevelModel

	ConstSubscriptionOptionName = "Subscription"

	ConstConfigPathSubscription        = "general.subscription"
	ConstConfigPathSubscriptionEnabled = "general.subscription.enabled"

	ConstConfigPathSubscriptionProducts = "general.subscription.products"

	ConstConfigPathSubscriptionEmailSubject  = "general.subscription.emailSubject"
	ConstConfigPathSubscriptionEmailTemplate = "general.subscription.emailTemplate"

	ConstConfigPathSubscriptionStockEmailSubject  = "general.subscription.emailStockSubject"
	ConstConfigPathSubscriptionStockEmailTemplate = "general.subscription.emailStockTemplate"

	ConstSubscriptionLogStorage = "subscription.log"

	ConstSubscriptionStatusSuspended = "suspended"
	ConstSubscriptionStatusConfirmed = "confirmed"
	ConstSubscriptionStatusCanceled  = "canceled"
)

var (
	optionValues = map[string]int{
		"Every 5 days": 5, "Every 30 days": 30, "Every 60 days": 60, "Every 90 days": 90, "Every 120 days": 120,
		"": 0, "10": 10, "30": 30, "60": 60, "90": 90, "hour": -1, "2hours": -2, "day": 1,
		"Just Once": 0, "30 days": 30, "60 days": 60, "90 days": 90, "120 days": 120,
	}
)
