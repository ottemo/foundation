package subscription

import (
	"strings"
	"time"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// Function for every hour check subscriptions to place an order
// placeOrders used to place orders for subscriptions
func placeOrders(params map[string]interface{}) error {

	if !subscriptionEnabled {
		return nil
	}

	currentHourBeginning := time.Now().Truncate(time.Hour)

	subscriptionCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if err := subscriptionCollection.AddFilter("action_date", ">=", currentHourBeginning); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7cd48d8b-5cc2-4475-8f64-93d88d5ddd02", err.Error())
	}
	if err := subscriptionCollection.AddFilter("action_date", "<", currentHourBeginning.Add(time.Hour)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "79a6da16-5c69-4ad7-8f3b-0ec10f9d4bca", err.Error())
	}
	if err := subscriptionCollection.AddFilter("status", "=", subscription.ConstSubscriptionStatusConfirmed); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b467c195-fa47-4d1d-b543-9f52efc82a60", err.Error())
	}

	//	get subscriptions with current day date and do action
	subscriptionsOnSubmit, err := subscriptionCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
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

		if err := checkoutInstance.SetInfo("subscription_id", subscriptionInstance.GetID()); err != nil {
			handleSubscriptionError(subscriptionInstance, err)
			continue
		}

		// need to check for unreached payment
		// to send email to user in case of low balance on credit card
		_, err = checkoutInstance.Submit()
		if err != nil {
			handleCheckoutError(subscriptionInstance, checkoutInstance, err)
			continue
		}

		if err := subscriptionInstance.Set("last_submit", time.Now()); err != nil {
			handleSubscriptionError(subscriptionInstance, err)
			continue
		}

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
			_ = env.ErrorDispatch(emailError)
			env.Log(subscription.ConstSubscriptionLogStorage, "Notification Error", subscriptionInstance.GetID()+": "+emailError.Error())
		}

		if internalError := subscriptionInstance.SetStatus(subscription.ConstSubscriptionStatusCanceled); internalError != nil {
			_ = env.ErrorDispatch(internalError)
		}

		if internalError := subscriptionInstance.Save(); internalError != nil {
			_ = env.ErrorDispatch(internalError)
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
