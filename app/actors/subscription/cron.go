package subscription

import (
	"fmt"
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
	fmt.Println(time.Now(), time.Now().Unix())

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
		fmt.Println(subscriptionInstance.GetID())

		checkoutInstance, err := subscriptionInstance.GetCheckout()
		if err != nil {
			fmt.Println(err)
			env.LogError(err)
			continue
		}

		checkoutInstance.SetInfo("subscription_id", subscriptionInstance.GetID())

		// need to check for unreached payment
		// to send email to user in case of low balance on credit card
		_, err = checkoutInstance.Submit()
		if err != nil {
			fmt.Println(err)
			handleCheckoutError(subscriptionInstance, checkoutInstance, err)
			continue
		}
		fmt.Println("sucess")

		// save new action date for current subscription
		subscriptionInstance.Set("last_submit", time.Now())

		subscriptionPeriod := subscriptionInstance.GetPeriod()

		actionDate := subscriptionInstance.GetActionDate()
		if subscriptionPeriod < 0 {
			actionDate = actionDate.Add(time.Hour * time.Duration(subscriptionPeriod*-1))
		} else {
			actionDate = actionDate.Add(ConstTimeDay * time.Duration(subscriptionPeriod))
		}

		if err = subscriptionInstance.SetActionDate(actionDate); err != nil {
			fmt.Println(err)
			env.LogError(err)
		}

		if err = subscriptionInstance.Save(); err != nil {
			fmt.Println(err)
			env.LogError(err)
		}
	}

	return nil
}

func handleCheckoutError(subscriptionInstance subscription.InterfaceSubscription, checkoutInstance checkout.InterfaceCheckout, err error) {
	env.LogError(err)
}
