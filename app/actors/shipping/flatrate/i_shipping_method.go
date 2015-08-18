package flatrate

import (
	"github.com/ottemo/foundation/app/models/checkout"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// GetName returns name of shipping method
func (it *ShippingMethod) GetName() string {
	return utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathName))
}

// GetCode returns code of shipping method
func (it *ShippingMethod) GetCode() string {
	return ConstShippingCode
}

// IsAllowed checks for method applicability
func (it *ShippingMethod) IsAllowed(checkout checkout.InterfaceCheckout) bool {
	return utils.InterfaceToBool(env.ConfigGetValue(ConstConfigPathEnabled))
}

// GetRates returns rates allowed by shipping method for a given checkout
func (it *ShippingMethod) GetRates(checkoutObject checkout.InterfaceCheckout) []checkout.StructShippingRate {

	result := []checkout.StructShippingRate{}

	if len(additionalRates) > 0 {
		for _, shippingRate := range additionalRates {
			shippingRates := utils.InterfaceToMap(shippingRate)
			if rateIsAllowed(shippingRates, checkoutObject) {
				result = append(result, checkout.StructShippingRate{
					Code:  utils.InterfaceToString(shippingRates["code"]),
					Name:  utils.InterfaceToString(shippingRates["title"]),
					Price: utils.InterfaceToFloat64(shippingRates["price"]),
				})
			}
		}

		if len(result) > 0 {
			return result
		}
	}

	return []checkout.StructShippingRate{
		checkout.StructShippingRate{
			Code:  "default",
			Name:  utils.InterfaceToString(env.ConfigGetValue(ConstConfigPathName)),
			Price: utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathAmount)),
		}}
}

// rateIsAllowed used to check rate for allowed by additional rules and presense of main keys
func rateIsAllowed(shippingRate map[string]interface{}, checkoutObject checkout.InterfaceCheckout) bool {

	if len(shippingRate) > 3 {
		subtotal := checkoutObject.GetSubtotal()
		country := "Default"

		shippingAddress := checkoutObject.GetShippingAddress()
		if shippingAddress != nil {
			country = utils.InterfaceToString(shippingAddress.GetCountry())
		}

		for limitingKey, limitingValue := range shippingRate {

			switch strings.ToLower(limitingKey) {
			case "price_from":
				if subtotal < utils.InterfaceToFloat64(limitingValue) {
					return false
				}
			case "price_to":
				if subtotal > utils.InterfaceToFloat64(limitingValue) {
					return false
				}
			case "banned_countries":
				if strings.Contains(utils.InterfaceToString(limitingValue), country) {
					return false
				}
			case "allowed_countries":
				if !strings.Contains(utils.InterfaceToString(limitingValue), country) {
					return false
				}
			}
		}
	}

	return true
}
