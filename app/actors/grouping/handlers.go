package grouping

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// updateCartHandler listening for cart update, on update thru over cart items grouping them by rules defined in config
func updateCartHandler(event string, eventData map[string]interface{}) bool {
	configVAl := env.GetConfig()
	rulesValue := configVAl.GetValue(ConstGroupingConfigPath)

	rules, err := utils.DecodeJSONToStringKeyMap(rulesValue)
	if err != nil {
		env.LogError(err)
		return false
	}

	rulesGroup := utils.InterfaceToArray(rules["group"])
	rulesInto := utils.InterfaceToArray(rules["into"])

	currentCart := eventData["cart"].(cart.InterfaceCart)
	cartChanged := false

	// check all rules and apply them before final version of cart will be created
	for {
		// Go thru all group products and apply possible combination
		for index, group := range rulesGroup {
			ruleInto := utils.InterfaceToArray(rulesInto[index])

			if ruleSetUsage := getGroupQty(currentCart.GetItems(), utils.InterfaceToArray(group)); ruleSetUsage > 0 {
				currentCart = applyGroupRule(currentCart, utils.InterfaceToArray(group), ruleInto, ruleSetUsage)
				cartChanged = true
			}
		}
		if !cartChanged {
			break
		}
		cartChanged = false
	}

	if err := currentCart.Save(); err != nil {
		env.LogError(err)
	}

	return true
}

// getGroupQty check cartItems for presence of product from one rule of grouping
// and calculate possible multiplier for it
func getGroupQty(currentCartItems []cart.InterfaceCartItem, groupProducts []interface{}) int {
	productsInCart := make(map[string]map[string]interface{})
	ruleMultiplier := 999

	for _, cartItem := range currentCartItems {
		productsInCart[cartItem.GetProductID()] = map[string]interface{}{"qty": cartItem.GetQty(), "options": cartItem.GetOptions()}
	}

	for _, groupProduct := range groupProducts {
		groupProduct := utils.InterfaceToMap(groupProduct)

		if value, present := productsInCart[utils.InterfaceToString(groupProduct["pid"])]; present {
			productMultiplier := int(utils.InterfaceToInt(value["qty"]) / utils.InterfaceToInt(groupProduct["qty"]))
			optionsApplly := true

			if groupProduct["options"] != nil && len(utils.InterfaceToMap(groupProduct["options"])) != 0 {
				optionsApplly = utils.MatchMapAValuesToMapB(utils.InterfaceToMap(groupProduct["options"]), utils.InterfaceToMap(value["options"]))
			}

			if productMultiplier >= 1 && optionsApplly {
				if productMultiplier < ruleMultiplier {
					ruleMultiplier = productMultiplier
				}
			} else {
				return 0
			}
		} else {
			return 0
		}
	}
	return ruleMultiplier
}

// applyGroupRule removes products in gruop rule with multiplier and add products from into rule
func applyGroupRule(currentCart cart.InterfaceCart, groupProducts, intoProducts []interface{}, multiplier int) cart.InterfaceCart {

	for _, cartItem := range currentCart.GetItems() {
		productCartID := cartItem.GetProductID()

		for _, product := range groupProducts {
			product := utils.InterfaceToMap(product)
			productID := utils.InterfaceToString(product["pid"])

			if productID == productCartID {
				if productNewQty := cartItem.GetQty() - utils.InterfaceToInt(product["qty"])*multiplier; productNewQty == 0 {
					currentCart.RemoveItem(cartItem.GetIdx())
				} else {
					cartItem.SetQty(productNewQty)
				}
				break
			}
		}
	}

	for _, product := range intoProducts {
		product := utils.InterfaceToMap(product)
		options := utils.InterfaceToMap(product["options"])

		if _, err := currentCart.AddItem(utils.InterfaceToString(product["pid"]), (utils.InterfaceToInt(product["qty"]) * multiplier), options); err != nil {
			env.LogError(err)
		}

	}

	return currentCart
}
