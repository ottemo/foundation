package subscription

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// checkoutSuccessHandler is a handler for checkout success event which sends order information to TrustPilot
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	var currentCheckout checkout.InterfaceCheckout
	if eventItem, present := eventData["checkout"]; present {
		if typedItem, ok := eventItem.(checkout.InterfaceCheckout); ok {
			currentCheckout = typedItem
		}
	}

	subscriptionID := currentCheckout.GetInfo("subscription")
	if subscriptionID == nil {
		return true
	}

	var checkoutOrder order.InterfaceOrder
	if eventItem, present := eventData["order"]; present {
		if typedItem, ok := eventItem.(order.InterfaceOrder); ok {
			checkoutOrder = typedItem
		}
	}

	if checkoutOrder != nil && currentCheckout != nil {
		go subscriptionUpdate(checkoutOrder, utils.InterfaceToString(subscriptionID))
	}

	return true
}

// subscriptionUpdate is a asynchronously update subscription with new state of the order
func subscriptionUpdate(checkoutOrder order.InterfaceOrder, subscriptionID string) error {

	subscriptionCollection, err := db.GetCollection(ConstCollectionNameSubscription)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	subscriptionCollection.AddFilter("_id", "=", subscriptionID)

	dbRecords, err := subscriptionCollection.Load()

	if len(dbRecords) == 0 {
		return env.ErrorDispatch(env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "5a4bd9ee-6ba7-4f7b-9e1a-24c254541582", "subscription not found"))
	}

	subscription := utils.InterfaceToMap(dbRecords[0])
	subscriptionAction := utils.InterfaceToString(subscription["action"])
	//	subscriptionDate := utils.InterfaceToTime(subscription["action_date"])
	//	subscriptionPeriod := utils.InterfaceToInt(subscription["period"])

	// update subscription with new order info (set new order id
	if subscriptionAction != ConstSubscriptionActionSubmit {

		//		subscriptionNextDate := subscriptionDate.AddDate(0, subscriptionPeriod, 0)
		//		subscription["action_date"] = subscriptionNextDate
		subscription["status"] = ConstSubscriptionStatusSuspended
		subscription["order_id"] = checkoutOrder.GetID()
		subscription["action"] = ConstSubscriptionActionUpdate

		if paymentInfo := utils.InterfaceToMap(checkoutOrder.Get("payment_info")); paymentInfo != nil {
			if _, present := paymentInfo["transactionID"]; present {
				subscription["action"] = ConstSubscriptionActionSubmit
			}
		}

		_, err = subscriptionCollection.Save(subscription)
		if err != nil {
			env.LogError(err)
		}
	}

	return nil
}
