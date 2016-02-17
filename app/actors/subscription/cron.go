package subscription

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"strings"
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
	subscriptionCollection.AddFilter("status", "=", subscription.ConstSubscriptionStatusConfirmed)

	//	get subscriptions with current day date and do action
	subscriptionsOnSubmit, err := subscriptionCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if !subscriptionEnabled {
		if len(subscriptionsOnSubmit) > 0 {
			env.Log(subscription.ConstSubscriptionLogStorage, "Warn", "Subscription is turned off, some records are not processed for this hour "+currentHourBeginning.String())
		}
		return nil
	}

	for _, record := range subscriptionsOnSubmit {

		subscriptionInstance, err := subscription.GetSubscriptionModel()
		if err != nil {
			handleSubscriptionError(subscriptionInstance, err)
			continue
		}

		err = subscriptionInstance.FromHashMap(record)
		if err != nil {
			handleSubscriptionError(subscriptionInstance, err)
			continue
		}

		checkoutInstance, err := subscriptionInstance.GetCheckout()
		if err != nil {
			handleSubscriptionError(subscriptionInstance, err)
			continue
		}

		checkoutInstance.SetInfo("subscription_id", subscriptionInstance.GetID())

		// need to check for unreached payment
		// to send email to user in case of low balance on credit card
		_, err = checkoutInstance.Submit()
		if err != nil {
			handleCheckoutError(subscriptionInstance, checkoutInstance, err)
			continue
		}

		subscriptionInstance.Set("last_submit", time.Now())

		// save new action date for current subscription
		if err = subscriptionInstance.UpdateActionDate(); err != nil {
			handleSubscriptionError(subscriptionInstance, err)
			continue
		}

		if err = subscriptionInstance.Save(); err != nil {
			handleSubscriptionError(subscriptionInstance, err)
		}
	}

	return nil
}

// handleCheckoutError will do the required actions with subscription
func handleCheckoutError(subscriptionInstance subscription.InterfaceSubscription, checkoutInstance checkout.InterfaceCheckout, err error) {

	errorMessage := err.Error()

	// handle notification of customers
	if strings.Contains(errorMessage, checkout.ConstPaymentErrorDeclined) {
		if emailError := sendNotificationEmail(subscriptionInstance); emailError != nil {
			env.ErrorDispatch(emailError)
			env.Log(subscription.ConstSubscriptionLogStorage, "Notification Error", subscriptionInstance.GetID()+": "+emailError.Error())
		}

		if internalError := subscriptionInstance.SetStatus(subscription.ConstSubscriptionStatusSuspended); internalError != nil {
			env.ErrorDispatch(internalError)
		}

		if internalError := subscriptionInstance.Save(); internalError != nil {
			env.ErrorDispatch(internalError)
		}

		return
	}

	handleSubscriptionError(subscriptionInstance, err)
}

func handleSubscriptionError(subscriptionInstance subscription.InterfaceSubscription, err error) {
	env.LogError(env.ErrorDispatch(err))

	if subscriptionInstance != nil {
		env.Log(subscription.ConstSubscriptionLogStorage, "Error", subscriptionInstance.GetID()+": "+err.Error())
	} else {
		env.Log(subscription.ConstSubscriptionLogStorage, "Error", err.Error())
	}
}
