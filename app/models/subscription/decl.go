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

	ConstSubscriptionOptionName = "subscription"

	ConstConfigPathSubscription        = "general.subscription"
	ConstConfigPathSubscriptionEnabled = "general.subscription.enabled"

	ConstConfigPathSubscriptionProducts = "general.subscription.products"

	// Admin: Out of stock email
	ConstConfigPathSubscriptionStockEmailSubject  = "general.subscription.emailStockSubject"
	ConstConfigPathSubscriptionStockEmailTemplate = "general.subscription.emailStockTemplate"

	// Inusfficient funds email
	ConstConfigPathSubscriptionEmailSubject  = "general.subscription.emailSubject"
	ConstConfigPathSubscriptionEmailTemplate = "general.subscription.emailTemplate"

	// Cancellation Email
	ConstConfigPathSubscriptionCancelEmailSubject  = "general.subscription.emailCancelSubject"
	ConstConfigPathSubscriptionCancelEmailTemplate = "general.subscription.emailCancelTemplate"

	ConstSubscriptionLogStorage = "subscription.log"

	ConstSubscriptionStatusSuspended = "suspended"
	ConstSubscriptionStatusConfirmed = "confirmed"
	ConstSubscriptionStatusCanceled  = "canceled"
)

var (
	optionValues = map[string]int{
		"every_5_days": 5, "every_30_days": 30, "every_60_days": 60, "every_90_days": 90, "every_120_days": 120,
		"": 0, "10": 10, "30": 30, "60": 60, "90": 90, "hour": -1, "2hours": -2, "day": 1,
		"just_once": 0, "30_days": 30, "60_days": 60, "90_days": 90, "120_days": 120,
	}
)
