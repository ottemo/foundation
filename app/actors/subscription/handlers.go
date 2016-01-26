package subscription

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
)

// checkoutSuccessHandler is a handler for checkout success event which sends order information to TrustPilot
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	if !subscriptionEnabled {
		return true
	}

	var currentCheckout checkout.InterfaceCheckout
	if eventItem, present := eventData["checkout"]; present {
		if typedItem, ok := eventItem.(checkout.InterfaceCheckout); ok {
			currentCheckout = typedItem
		}
	}

	// means current order is placed by subscription handler
	if currentCheckout == nil || currentCheckout.GetInfo("subscription_id") != nil {
		return true
	}

	// allows subscription only for registered
	//	if currentCheckout.GetVisitor() == nil {
	//		return true
	//	}

	var checkoutOrder order.InterfaceOrder
	if eventItem, present := eventData["order"]; present {
		if typedItem, ok := eventItem.(order.InterfaceOrder); ok {
			checkoutOrder = typedItem
		}
	}

	if checkoutOrder != nil {
		go subscriptionCreate(currentCheckout, checkoutOrder)
	}

	return true
}

// subscriptionCreate is a asynchronously used to create subscription based on finished checkout
func subscriptionCreate(currentCheckout checkout.InterfaceCheckout, checkoutOrder order.InterfaceOrder) error {

	currentCart := currentCheckout.GetCart()
	if currentCart == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "ae108000-68ff-419f-b443-2df1554dd377", "No cart")
	}

	subscriptionItems := make(map[int]int)
	for _, cartItem := range currentCart.GetItems() {
		itemOptions := cartItem.GetOptions()
		if optionValue, present := itemOptions[optionName]; present {
			subscriptionItems[cartItem.GetIdx()] = getPeriodValue(utils.InterfaceToString(optionValue))
		}
	}

	if len(subscriptionItems) == 0 {
		return nil
	}

	subscriptionInstance, err := subscription.GetSubscriptionModel()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	visitorCreditCard := retrieveCreditCard(currentCheckout, checkoutOrder)
	if visitorCreditCard == nil || visitorCreditCard.GetToken() == "" {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "333d3396-fddc-4aff-a3fe-083e50a2e1a6", "Credit card can't be obtained")
	}

	if err := validateCheckoutToSubscribe(currentCheckout); err != nil {
		return env.ErrorDispatch(err)
	}

	if err = subscriptionInstance.SetCreditCard(visitorCreditCard); err != nil {
		return env.ErrorDispatch(err)
	}

	visitor := currentCheckout.GetVisitor()
	if visitor != nil {
		subscriptionInstance.Set("visitor_id", visitor.GetID())
		subscriptionInstance.Set("customer_email", visitor.GetEmail())
		subscriptionInstance.Set("customer_name", visitor.GetFullName())
	} else {
		subscriptionInstance.Set("customer_email", currentCheckout.GetInfo("customer_email"))
		subscriptionInstance.Set("customer_name", currentCheckout.GetInfo("customer_name"))
	}

	subscriptionInstance.SetShippingAddress(currentCheckout.GetShippingAddress())
	subscriptionInstance.SetBillingAddress(currentCheckout.GetBillingAddress())
	subscriptionInstance.SetShippingMethod(currentCheckout.GetShippingMethod())
	subscriptionInstance.SetStatus(ConstSubscriptionStatusConfirmed)
	subscriptionInstance.Set("order_id", checkoutOrder.GetID())

	subscriptionTime := time.Now().Truncate(time.Hour)

	// create different subscriptions for every subscription product
	for _, cartItem := range currentCart.GetItems() {
		if subscriptionPeriodValue, present := subscriptionItems[cartItem.GetIdx()]; present && subscriptionPeriodValue != 0 {

			if err = subscriptionInstance.SetActionDate(subscriptionTime); err != nil {
				env.LogError(err)
				continue
			}

			if err = subscriptionInstance.SetPeriod(subscriptionPeriodValue); err != nil {
				env.LogError(err)
				continue
			}

			if err = subscriptionInstance.UpdateActionDate(); err != nil {
				env.LogError(err)
				continue
			}

			productCart, err := cart.GetCartModel()
			if err != nil {
				env.LogError(err)
				continue
			}

			if _, err = productCart.AddItem(cartItem.GetProductID(), cartItem.GetQty(), cartItem.GetOptions()); err != nil {
				env.LogError(err)
				continue
			}

			if visitor != nil {
				productCart.SetVisitorID(visitor.GetID())
			}

			if currentSession := currentCheckout.GetSession(); currentSession != nil {
				productCart.SetSessionID(currentSession.GetID())
			}

			if err = productCart.Deactivate(); err != nil {
				env.LogError(err)
				continue
			}

			if err = productCart.Save(); err != nil {
				env.LogError(err)
				continue
			}

			subscriptionInstance.Set("cart_id", productCart.GetID())
			subscriptionInstance.SetID("")

			if err = subscriptionInstance.Save(); err != nil {
				env.LogError(err)
				continue
			}
		}
	}

	return nil
}

// getOptionsExtend is a handler for product get options event which extend available product options
// TODO: create some defined object for options (should explain keys)
func getOptionsExtend(event string, eventData map[string]interface{}) bool {

	if !subscriptionEnabled {
		return true
	}

	if value, present := eventData["options"]; present {
		options := utils.InterfaceToMap(value)

		storedOptions := map[string]interface{}{
			"type":     "select",
			"required": true,
			"order":    1,
			"label":    "Subscription",
			"options": map[string]interface{}{
				"Just Once": map[string]interface{}{"order": 1, "label": "Just Once"},
				"30 days":   map[string]interface{}{"order": 2, "label": "30 days"},
				"60 days":   map[string]interface{}{"order": 3, "label": "60 days"},
				"90 days":   map[string]interface{}{"order": 4, "label": "90 days"},
				"120 days":  map[string]interface{}{"order": 5, "label": "120 days"},
			},
		}

		// when we are using getOptions for product after they was applied there add field Value, but is it all?
		if subscriptionOption, present := options[optionName]; present {
			subscriptionOptionMap := utils.InterfaceToMap(subscriptionOption)
			if appliedValue, present := subscriptionOptionMap["value"]; present {
				storedOptions["value"] = appliedValue
			}
		}

		options[optionName] = storedOptions
	}
	return true
}
