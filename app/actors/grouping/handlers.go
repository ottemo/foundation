package grouping

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// updateCartHandler listening for cart update, on update thru over cart items grouping them by rules defined in config
func updateCartHandler(event string, eventData map[string]interface{}) bool {

	currentCart := eventData["cart"].(cart.InterfaceCart)

	for _, ruleValue := range currentRules {
		ruleElement := utils.InterfaceToMap(ruleValue)

		ruleGroup := utils.InterfaceToArray(ruleElement["group"])
		ruleInto := utils.InterfaceToArray(ruleElement["into"])

		applyGroupRule(currentCart, ruleGroup, ruleInto)
	}

	if err := currentCart.Save(); err != nil {
		env.LogError(err)
	}

	return true
}

// getApplyTimesCount returns count of times grouping rule could be applied to cart
func getApplyTimesCount(currentCart cart.InterfaceCart, groupProductsArray []interface{}) int {
	ruleMultiplier := 0 // minimal possible amount

	// loop over grouping rule items
	for _, groupProduct := range groupProductsArray {
		groupProductElement := utils.InterfaceToMap(groupProduct)

		// decoding rule product specification
		groupProductID := utils.InterfaceToString(groupProductElement["pid"])
		groupProductQty := utils.InterfaceToInt(groupProductElement["qty"])
		groupProductOptions := make(map[string]interface{})
		if optionsValue, present := groupProductElement["options"]; present {
			groupProductOptions = utils.InterfaceToMap(optionsValue)
		}

		// looking for match among cart items
		cartItemFoundFlag := false
		for _, cartItem := range currentCart.GetItems() {
			if groupProductID == cartItem.GetProductID() && utils.MatchMapAValuesToMapB(groupProductOptions, cartItem.GetOptions()) {
				possibleAppliesCount := cartItem.GetQty() / groupProductQty

				// determination of minimal possible amount
				if possibleAppliesCount > 0 {
					if ruleMultiplier == 0 || possibleAppliesCount < ruleMultiplier {
						ruleMultiplier = possibleAppliesCount
					}
					cartItemFoundFlag = true
				} else {
					return 0
				}
			}
		}
		if !cartItemFoundFlag {
			return 0
		}
	}
	return ruleMultiplier
}

// applyGroupRule applies grouping rule to current cart (changing cart items)
func applyGroupRule(currentCart cart.InterfaceCart, groupProductsArray []interface{}, intoProductsArray []interface{}) {

	// checking how many times rule could be applied
	ruleMultiplier := getApplyTimesCount(currentCart, groupProductsArray)
	if ruleMultiplier > 0 {

		// modifying current cart with removing items from rule group key element
		for _, groupProductValue := range groupProductsArray {
			groupProductElement := utils.InterfaceToMap(groupProductValue)

			groupProductID := utils.InterfaceToString(groupProductElement["pid"])
			groupProductQty := utils.InterfaceToInt(groupProductElement["qty"])
			groupProductOptions := make(map[string]interface{})
			if optionsValue, present := groupProductElement["options"]; present {
				groupProductOptions = utils.InterfaceToMap(optionsValue)
			}

			for _, cartItem := range currentCart.GetItems() {
				qtyToReduce := ruleMultiplier * groupProductQty
				if groupProductID == cartItem.GetProductID() && utils.MatchMapAValuesToMapB(groupProductOptions, cartItem.GetOptions()) {

					if cartProductNewQty := cartItem.GetQty() - qtyToReduce; cartProductNewQty > 0 {
						cartItem.SetQty(cartProductNewQty)
						break
					} else {
						if cartProductNewQty == 0 {
							currentCart.RemoveItem(cartItem.GetIdx())
							break
						}
					}
				}
			}
		}

		// modifying current cart with increasing items from rule into key element
		for _, intoProductValue := range intoProductsArray {
			intoProductElement := utils.InterfaceToMap(intoProductValue)

			intoProductPID := utils.InterfaceToString(intoProductElement["pid"])
			intoProductQty := utils.InterfaceToInt(intoProductElement["qty"])
			intoProductOptions := make(map[string]interface{})
			if optionsValue, present := intoProductElement["options"]; present {
				intoProductOptions = utils.InterfaceToMap(optionsValue)
			}

			_, err := currentCart.AddItem(intoProductPID, intoProductQty*ruleMultiplier, intoProductOptions)
			if err != nil {
				env.LogError(err)
			}
		}
	}
}
