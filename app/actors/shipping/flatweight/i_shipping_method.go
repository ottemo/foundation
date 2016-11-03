package flatweight

import (
	"strings"

	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
)

// GetName return the internal name of the shipping method
func (it ShippingMethod) GetName() string {
	return ConstShippingName
}

// GetCode returns the shipping method's code
func (it ShippingMethod) GetCode() string {
	return ConstShippingCode
}

// IsAllowed Whether the method is enabled
func (it ShippingMethod) IsAllowed(checkoutInstance checkout.InterfaceCheckout) bool {
	return configIsEnabled()
}

// GetRates returns the rates allowed for this order weight, based on config criteria
func (it ShippingMethod) GetRates(checkoutInstance checkout.InterfaceCheckout) []checkout.StructShippingRate {
	var allowedRates []checkout.StructShippingRate

	// Quick exit
	if len(rates) == 0 {
		return allowedRates
	}

	// Calculate order weight
	var orderWeight float64
	for _, cartItem := range checkoutInstance.GetItems() {
		p := cartItem.GetProduct()
		if p != nil { // returns an interface which could be nil
			orderWeight += p.GetWeight() * float64(cartItem.GetQty())
		}
	}

	// Gather allowed rates
	for _, r := range rates {
		if r.validForWeight(orderWeight) && rateIsAllowedCountry(utils.InterfaceToMap(r), checkoutInstance) {
			allowedRates = append(allowedRates, r.toCheckoutStruct())
		}
	}

	return allowedRates
}

// GetAllRates returns a list of all rates that are configured
func (it ShippingMethod) GetAllRates() []checkout.StructShippingRate {
	var allRates []checkout.StructShippingRate
	for _, r := range rates {
		allRates = append(allRates, r.toCheckoutStruct())
	}

	return allRates
}

// rateIsAllowedCountry used to check rate for allowed by country rules and presense of main keys
func rateIsAllowedCountry(shippingRate map[string]interface{}, checkoutObject checkout.InterfaceCheckout) bool {

	if len(shippingRate) > 3 {
		country := "Default"

		shippingAddress := checkoutObject.GetShippingAddress()
		if shippingAddress != nil {
			country = utils.InterfaceToString(shippingAddress.GetCountry())
		}

		for limitingKey, limitingValue := range shippingRate {

			switch strings.ToLower(limitingKey) {
			case "bannedcountries":
				if strings.Contains(utils.InterfaceToString(limitingValue), country) {
					return false
				}
			case "allowedcountries":
				limitingString := utils.InterfaceToString(limitingValue)
				if len(limitingString) > 0 && !strings.Contains(limitingString, country) {
					return false
				}
			}
		}
	}

	return true
}
