package grouping

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/env"
)

func updateCartHandler(event string, eventData map[string]interface{}) bool {
	configVAl := env.GetConfig()
	rulesValue := configVAl.GetValue(ConstGroupingConfigPath)

	rules, err := utils.DecodeJSONToArray(rulesValue)
	if err != nil {
		env.LogError(err)
	}
	if rules == nil {
		return true
	}

	rulesGroup := utils.InterfaceToArray(rules[0])
	rulesInto := utils.InterfaceToArray(rules[1])

	currentCart := eventData["cart"].(cart.InterfaceCart)

	// Go thru all group products and apply possible combination
		for index, _ := range rulesGroup {
			group := utils.InterfaceToArray(rulesGroup[index])
			ruleInto := utils.InterfaceToArray(rulesInto[index])

			if ruleSetUsage := getGroupQty(currentCart.GetItems(), group); ruleSetUsage > 0 {
				currentCart = applyGroupRule(currentCart, group, ruleInto, ruleSetUsage)
			}
		}
	if err := currentCart.Save(); err != nil {
		env.LogError(err)
	}

	return true
}

// getGroupQty check cartItems for presence of product from one rule of grouping
// and calculate possible multiplier for it
func getGroupQty(currentCartItems []cart.InterfaceCartItem, groupProducts []interface {}) int {
	productsInCart := make(map[string]int)
	ruleMultiplier := 999

	for _, cartItem := range currentCartItems {
		productsInCart[cartItem.GetProductID()] = cartItem.GetQty()
	}

	for key, _ := range groupProducts {
		groupProduct := utils.InterfaceToMap(groupProducts[key])

		if value, present := productsInCart[utils.InterfaceToString(groupProduct["pid"])]; present {
			if productMultiplier := int(value/utils.InterfaceToInt(groupProduct["qty"])); productMultiplier >= 1 {
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
func applyGroupRule (currentCart cart.InterfaceCart, groupProducts, intoProducts []interface {}, multiplier int) cart.InterfaceCart {

	for _, cartItem := range currentCart.GetItems() {
		productCartId := cartItem.GetProductID()

		for key, _ := range groupProducts {
			product := utils.InterfaceToMap(groupProducts[key])
			productID := utils.InterfaceToString(product["pid"])

			if productID == productCartId {
				if productNewQty := cartItem.GetQty() - utils.InterfaceToInt(product["qty"]) * multiplier; productNewQty == 0 {
					currentCart.RemoveItem(cartItem.GetIdx())
				} else	{
					cartItem.SetQty(productNewQty)
				}
				break
			}
		}
	}

	for key, _ := range intoProducts {
		product := utils.InterfaceToMap(intoProducts[key])
		options := utils.InterfaceToMap(product["options"])

		if _, err := currentCart.AddItem(utils.InterfaceToString(product["pid"]), utils.InterfaceToInt(product["qty"])*multiplier, options); err != nil {
			env.LogError(err)
		}
	}

	return currentCart
}
