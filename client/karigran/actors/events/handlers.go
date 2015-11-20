package events

import (
	"encoding/json"
	"fmt"
	"github.com/ottemo/foundation/app/actors/other/mailchimp"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {
	//FIXME: remove when code works
	fmt.Println("checkout success handler called")
	var checkoutOrder order.InterfaceOrder
	if eventItem, present := eventData["order"]; present {
		if typedItem, ok := eventItem.(order.InterfaceOrder); ok {
			checkoutOrder = typedItem
		}
	}

	var subscriptionConfigs []MailSubscriptionConfig
	if cfgJson := utils.InterfaceToString(env.ConfigGetValue(ConstKGMailSubscriptionEvent)); cfgJson != nil {
		if err := json.Unmarshal([]byte(cfgJson), &subscriptionConfigs); err != nil {
			env.ErrorDispatch(err)
			return false
		}
	}

	//FIXME: remove when code works
	fmt.Println("subscription configs loaded" + subscriptionConfigs[0].ListId)

	if checkoutOrder != nil { //Is it possible at this point for the checkoutOrder to be nil?
		go processOrder(checkoutOrder, subscriptionConfigs)
	}

	return true

}

func processOrder(order order.InterfaceOrder, subscriptionConfigs []MailSubscriptionConfig) {
	for _, config := range subscriptionConfigs {
		if orderContainsSku := containsItem(order, config.Sku); orderContainsSku {
			if ok, err := mailchimp.Subscribe(order, config.Sku); err != nil && !ok {
				//TODO: How to send mail? If the call to subscribe fails, KG would like to get an email
			}
		} else {
			//FIXME: remove when code works
			fmt.Println("In process order, no valid orders found")
		}
	}
}

func containsItem(order order.InterfaceOrder, sku string) bool {
	for _, item := range order.GetItems() {
		if item.GetSku() == sku {
			return true
		}
	}
	return false
}
