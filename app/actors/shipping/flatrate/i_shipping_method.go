package flatrate

import (
	"strings"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
)

// GetName returns name of shipping method
func (it *ShippingMethod) GetName() string {
	return ConstShippingName
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

	if len(flatRates) < 1 {
		_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "8945e737-c264-47cb-a68a-c112aafdfccb", "Unable to parse rates or no rates have been configured")

	} else {
		for _, shippingRate := range flatRates {
			shippingRates := utils.InterfaceToMap(shippingRate)
			if rateIsAllowed(shippingRates, checkoutObject) {
				result = append(result, checkout.StructShippingRate{
					Code:  utils.InterfaceToString(shippingRates["code"]),
					Name:  utils.InterfaceToString(shippingRates["title"]),
					Price: utils.InterfaceToFloat64(shippingRates["price"]),
				})
			}
		}
	}

	return result
}

// rateIsAllowed used to check rate for allowed by additional rules and presence of main keys
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

// GetAllRates returns an unfiltered list of all supported shipping rates using the Flat Rate method.
func (it ShippingMethod) GetAllRates() []checkout.StructShippingRate {
	result := []checkout.StructShippingRate{}
	for _, shippingRate := range flatRates {
		shippingRates := utils.InterfaceToMap(shippingRate)
		rate := checkout.StructShippingRate{
			Code:  utils.InterfaceToString(shippingRates["code"]),
			Name:  utils.InterfaceToString(shippingRates["title"]),
			Price: utils.InterfaceToFloat64(shippingRates["price"]),
		}
		result = append(result, rate)
	}

	return result
}
