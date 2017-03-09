package subscription

import (
	"strings"
	"time"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// checkoutSuccessHandler is a handler for checkout success event which creates the subscriptions
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
	if currentCheckout == nil || !currentCheckout.IsSubscription() || currentCheckout.GetInfo("subscription_id") != nil {
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
		go func() {
			if err := subscriptionCreate(currentCheckout, checkoutOrder); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c55060e4-f29a-43e6-9c6b-5345ae1652f1", err.Error())
			}
		}();
	}

	return true
}

// subscriptionCreate is invoked via a go routine to create subscription based on finished checkout
func subscriptionCreate(currentCheckout checkout.InterfaceCheckout, checkoutOrder order.InterfaceOrder) error {

	currentCart := currentCheckout.GetCart()
	if currentCart == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "ae108000-68ff-419f-b443-2df1554dd377", "No cart")
	}

	subscriptionItems := make(map[int]int)
	for _, cartItem := range currentCart.GetItems() {
		itemOptions := cartItem.GetOptions()
		if optionValue, present := itemOptions[subscription.ConstSubscriptionOptionName]; present {
			subscriptionItems[cartItem.GetIdx()] = subscription.GetSubscriptionPeriodValue(utils.InterfaceToString(optionValue))
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
		if err := subscriptionInstance.Set("visitor_id", visitor.GetID()); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3eef50f0-4bda-41f4-9da8-b80e7a4a07f2", err.Error())
		}
		if err := subscriptionInstance.Set("customer_email", visitor.GetEmail()); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "12f4bfb9-0abe-4d6c-8b09-b3ef177144d3", err.Error())
		}
		if err := subscriptionInstance.Set("customer_name", visitor.GetFullName()); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "71049641-67eb-47dc-982d-d2d86475132c", err.Error())
		}
	} else {
		if err := subscriptionInstance.Set("customer_email", currentCheckout.GetInfo("customer_email")); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8c1420f9-10c9-4904-ac3f-7de3037b85a1", err.Error())
		}
		if err := subscriptionInstance.Set("customer_name", currentCheckout.GetInfo("customer_name")); err != nil {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d451050d-9a2c-4b75-a031-c7df7fb887c8", err.Error())
		}
	}

	if err := subscriptionInstance.SetShippingAddress(currentCheckout.GetShippingAddress()); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1ba46a05-bb87-4c52-8340-d371a274dccc", err.Error())
	}
	if err := subscriptionInstance.SetBillingAddress(currentCheckout.GetBillingAddress()); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "015edd41-f549-468c-a129-9387b5f42042", err.Error())
	}

	shippingMethod := currentCheckout.GetShippingMethod()
	var shippingRate checkout.StructShippingRate

	if checkoutShippingRate := currentCheckout.GetShippingRate(); checkoutShippingRate != nil {
		shippingRate.Code = checkoutShippingRate.Code
		shippingRate.Name = checkoutShippingRate.Name
		shippingRate.Price = checkoutShippingRate.Price
	}

	// obtaining values of shipping method and rate from order if they weren't provided in checkout
	if shippingMethod == nil || shippingRate.Code == "" {

		shippingParts := strings.Split(checkoutOrder.GetShippingMethod(), "/")
		orderShippingMethod := checkout.GetShippingMethodByCode(shippingParts[0])

		for _, orderShippingRate := range orderShippingMethod.GetRates(currentCheckout) {
			if shippingParts[1] == orderShippingRate.Code {
				shippingRate = checkout.StructShippingRate{
					Name:  orderShippingRate.Name,
					Code:  orderShippingRate.Code,
					Price: orderShippingRate.Price,
				}
				shippingMethod = orderShippingMethod

				break
			}
		}
	}

	if err := subscriptionInstance.SetShippingMethod(shippingMethod); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f894070c-044c-4258-9902-b51ad4800e93", err.Error())
	}
	if err := subscriptionInstance.SetShippingRate(checkout.StructShippingRate{
		Name:  shippingRate.Name,
		Code:  shippingRate.Code,
		Price: shippingRate.Price,
	}); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6a1b1715-8ab4-40e9-986b-ce2769eb3258", err.Error())
	}

	if err := subscriptionInstance.SetStatus(subscription.ConstSubscriptionStatusConfirmed); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4ceb2aee-cd8c-4251-8080-4e0c20c72cc6", err.Error())
	}
	if err := subscriptionInstance.Set("order_id", checkoutOrder.GetID()); err != nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "25d58d3e-f857-4dd7-bf6b-7e482deeb276", err.Error())
	}

	subscriptionTime := time.Now().Truncate(time.Hour)

	// create unique subscriptions for every subscription product
	for _, cartItem := range currentCart.GetItems() {
		if subscriptionPeriodValue, present := subscriptionItems[cartItem.GetIdx()]; present && subscriptionPeriodValue != 0 {

			if err = subscriptionInstance.SetActionDate(subscriptionTime); err != nil {
				_ = env.ErrorDispatch(err)
				continue
			}

			if err = subscriptionInstance.SetPeriod(subscriptionPeriodValue); err != nil {
				_ = env.ErrorDispatch(err)
				continue
			}

			if err = subscriptionInstance.UpdateActionDate(); err != nil {
				_ = env.ErrorDispatch(err)
				continue
			}

			var items []subscription.StructSubscriptionItem

			// populate the subscription object
			subscriptionItem := subscription.StructSubscriptionItem{
				Name:      "",
				ProductID: cartItem.GetProductID(),
				Qty:       cartItem.GetQty(),
				Options:   cartItem.GetOptions(),
			}

			if product := cartItem.GetProduct(); product != nil {
				subscriptionItem.Name = product.GetName()
				subscriptionItem.Sku = product.GetSku()
				subscriptionItem.Price = product.GetPrice()

				productOptions := make(map[string]interface{})
				// add options to subscription info as description that used to show on FED
				fullOptions := product.GetOptions()
				subscriptionInstance.SetInfo("detail_options", fullOptions)

				for key, value := range fullOptions {
					option := utils.InterfaceToMap(value)
					optionLabel := key
					if labelValue, optionLabelPresent := option["label"]; optionLabelPresent {
						optionLabel = utils.InterfaceToString(labelValue)
					}

					optionValue, optionValuePresent := option["value"]
					productOptions[optionLabel] = optionValue

					// in this case looks like structure of options was changed or it's not a map
					if !optionValuePresent {
						productOptions[optionLabel] = value
						continue
					}

					optionType := ""
					if val, present := option["type"]; present {
						optionType = utils.InterfaceToString(val)
					}
					if options, present := option["options"]; present {
						optionsMap := utils.InterfaceToMap(options)

						if optionType == "multi_select" {
							selectedOptions := ""
							for i, optionValue := range utils.InterfaceToArray(optionValue) {
								if optionValueParameters, ok := optionsMap[utils.InterfaceToString(optionValue)]; ok {
									optionValueParametersMap := utils.InterfaceToMap(optionValueParameters)
									if labelValue, labelValuePresent := optionValueParametersMap["label"]; labelValuePresent {
										productOptions[optionLabel] = labelValue
										if i > 0 {
											selectedOptions = selectedOptions + ", "
										}
										selectedOptions = selectedOptions + utils.InterfaceToString(labelValue)
									}
								}
							}
							productOptions[optionLabel] = selectedOptions

						} else if optionValueParameters, ok := optionsMap[utils.InterfaceToString(optionValue)]; ok {
							optionValueParametersMap := utils.InterfaceToMap(optionValueParameters)
							if labelValue, labelValuePresent := optionValueParametersMap["label"]; labelValuePresent {
								productOptions[optionLabel] = labelValue
							}

						}
					}
				}

				subscriptionInstance.SetInfo("options", productOptions)
			}

			items = append(items, subscriptionItem)

			if err := subscriptionInstance.Set("items", items); err != nil {
				return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "bbff14b1-817f-4321-a9ab-72ba836b8a3a", err.Error())
			}
			if err := subscriptionInstance.SetID(""); err != nil {
				return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelActor, "cf7abf5c-51a6-44e7-8709-a09912f8d2e6", err.Error())
			}

			if err = subscriptionInstance.Save(); err != nil {
				_ = env.ErrorDispatch(err)
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

		// removing subscription option for products that are not in the list
		if len(subscriptionProducts) > 0 {
			if productID, present := eventData["id"]; !present || !utils.IsInListStr(utils.InterfaceToString(productID), subscriptionProducts) {
				delete(options, subscription.ConstSubscriptionOptionName)
				return true
			}
		}

		storedOptions := map[string]interface{}{
			"type":     "select",
			"required": true,
			"order":    1,
			"label":    "Subscription",
			"options": map[string]interface{}{
				"just_once": map[string]interface{}{"order": 1, "label": "Just Once"},
				"30_days":   map[string]interface{}{"order": 2, "label": "30 days"},
				"60_days":   map[string]interface{}{"order": 3, "label": "60 days"},
				"90_days":   map[string]interface{}{"order": 4, "label": "90 days"},
				"120_days":  map[string]interface{}{"order": 5, "label": "120 days"},
			},
		}

		// when we are using getOptions for product after they was applied there add field Value
		if subscriptionOption, present := options[subscription.ConstSubscriptionOptionName]; present {
			subscriptionOptionMap := utils.InterfaceToMap(subscriptionOption)
			if appliedValue, present := subscriptionOptionMap["value"]; present {
				storedOptions["value"] = appliedValue
			}
		}

		options[subscription.ConstSubscriptionOptionName] = storedOptions
	}
	return true
}

func updateCronJob(newExecutionTimeOption interface{}) error {

	executionTimeOption := utils.InterfaceToString(newExecutionTimeOption)
	executionTimeCronExpr := subscription.GetSubscriptionCronExpr(
		subscription.GetSubscriptionPeriodValue(executionTimeOption))

	if scheduler := env.GetScheduler(); scheduler != nil {
		schedules := scheduler.ListSchedules() //[]InterfaceSchedule
		var schedule env.InterfaceSchedule

		for _, schedule = range schedules {
			if schedule.Get("task") == ConstSchedulerTaskName {
				if err := schedule.Set("expr", executionTimeCronExpr); err != nil {
					return env.ErrorDispatch(err)
				}

				if err := schedule.Set("time", time.Now()); err != nil {
					return env.ErrorDispatch(err)
				}
			}
		}
	}
	return nil
}
