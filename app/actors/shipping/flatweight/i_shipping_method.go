package flatweight

import (
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
	for _, item := range checkoutInstance.GetItems() {
		p := item.GetProduct()
		if p != nil {
			orderWeight += p.GetWeight()
		}
	}

	// Gather allowed rates
	for _, r := range rates {
		if r.validForWeight(orderWeight) {
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
