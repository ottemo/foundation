package subscription

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
)

// Function for every hour check subscriptions to place an order
// placeOrders used to place orders for subscriptions
func placeOrders(params map[string]interface{}) error {

	currentHourBeginning := time.Now().Truncate(time.Hour)

	subscriptionCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	subscriptionCollection.AddFilter("action_date", ">=", currentHourBeginning)
	subscriptionCollection.AddFilter("action_date", "<", currentHourBeginning.Add(time.Hour))
	subscriptionCollection.AddFilter("status", "=", ConstSubscriptionStatusConfirmed)

	//	get subscriptions with current day date and do action
	subscriptionsOnSubmit, err := subscriptionCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, record := range subscriptionsOnSubmit {

		subscriptionInstance, err := subscription.GetSubscriptionModel()
		if err != nil {
			env.LogError(err)
			continue
		}

		err = subscriptionInstance.FromHashMap(record)
		if err != nil {
			env.LogError(err)
			continue
		}

		checkoutInstance, err := subscriptionInstance.GetCheckout()
		if err != nil {
			env.LogError(err)
			continue
		}

		// need to check for unreached payment
		// to send email to user in case of low balance on credit card
		_, err = checkoutInstance.Submit()
		if err != nil {
			handleCheckoutError(subscriptionInstance, checkoutInstance, err)
			continue
		}

		// save new action date for current subscription
		subscriptionInstance.Set("last_submit", time.Now())
		subscriptionInstance.SetActionDate(subscriptionInstance.GetActionDate().AddDate(0, subscriptionInstance.GetPeriod(), 0))
		err = subscriptionInstance.Save()
		if err != nil {
			env.LogError(err)
		}
	}

	return nil
}

func handleCheckoutError(subscriptionInstance subscription.InterfaceSubscription, checkoutInstance checkout.InterfaceCheckout, err error) {
	env.LogError(err)
}
