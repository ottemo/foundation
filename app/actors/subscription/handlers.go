package subscription

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
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

	var checkoutOrder order.InterfaceOrder
	if eventItem, present := eventData["order"]; present {
		if typedItem, ok := eventItem.(order.InterfaceOrder); ok {
			checkoutOrder = typedItem
		}
	}

	if checkoutOrder != nil && currentCheckout != nil {
		go subscriptionCreate(currentCheckout, checkoutOrder)
	}

	return true
}

// subscriptionCreate is a asynchronously used to create subscription based on finished checkout
func subscriptionCreate(currentCheckout checkout.InterfaceCheckout, checkoutOrder order.InterfaceOrder) error {

	currentCart := currentCheckout.GetCart()
	if currentCart == nil {
		return nil
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

	return nil
}

// getOptionsExtend is a handler for product get options event which extend available product options
func getOptionsExtend(event string, eventData map[string]interface{}) bool {
	if value, present := eventData["options"]; present {
		// "Subscription":{"type":"select","required":false,"order":1,"label":"Subscription","options":{"Every 5 days":{"order":1,"label":"Every 5 days"},"Every 15 days":{"order":2,"label":"Every 15 days"},"Every 30 days":{"order":3,"label":"Every 30 days"}}}
		options := utils.InterfaceToMap(value)
		options[optionName] = map[string]interface{}{
			"type":     "select",
			"required": false,
			"order":    1,
			"label":    "Subscription",
			"options":  map[string]interface{}{"Every 5 days": map[string]interface{}{"order": 1, "label": "Every 5 days"}},
		}
	}
	return true
}
