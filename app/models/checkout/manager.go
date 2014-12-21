package checkout

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	registeredShippingMethods = make([]InterfaceShippingMethod, 0)
	registeredPaymentMethods  = make([]InterfacePaymentMethod, 0)

	registeredTaxes     = make([]InterfaceTax, 0)
	registeredDiscounts = make([]InterfaceDiscount, 0)
)

// RegisterShippingMethod registers given shipping method in system
func RegisterShippingMethod(shippingMethod InterfaceShippingMethod) error {
	for _, registeredMethod := range registeredShippingMethods {
		if registeredMethod == shippingMethod {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7b532e87b8ca4b9b90ec4a09b37bf7e2", "shipping method already registered")
		}
	}

	registeredShippingMethods = append(registeredShippingMethods, shippingMethod)

	return nil
}

// RegisterPaymentMethod registers given payment method in system
func RegisterPaymentMethod(paymentMethod InterfacePaymentMethod) error {
	for _, registeredMethod := range registeredPaymentMethods {
		if registeredMethod == paymentMethod {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "177091a820294dd7a2a2e09fa6efa0c8", "payment method already registered")
		}
	}

	registeredPaymentMethods = append(registeredPaymentMethods, paymentMethod)

	return nil
}

// RegisterTax registers given tax calculator in system
func RegisterTax(tax InterfaceTax) error {
	for _, registeredTax := range registeredTaxes {
		if registeredTax == tax {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "355841ed700e4c8fbc8da2d95844d393", "tax already registered")
		}
	}

	registeredTaxes = append(registeredTaxes, tax)

	return nil
}

// RegisterDiscount registers given discount calculator in system
func RegisterDiscount(discount InterfaceDiscount) error {
	for _, registeredDiscount := range registeredDiscounts {
		if registeredDiscount == discount {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "56234554d230403cb21b8393e7e138d4", "discount already registered")
		}
	}

	registeredDiscounts = append(registeredDiscounts, discount)

	return nil
}

// GetRegisteredShippingMethods returns list of registered shipping methods
func GetRegisteredShippingMethods() []InterfaceShippingMethod {
	return registeredShippingMethods
}

// GetRegisteredPaymentMethods returns list of registered payment methods
func GetRegisteredPaymentMethods() []InterfacePaymentMethod {
	return registeredPaymentMethods
}

// GetRegisteredTaxes returns list of registered tax calculators
func GetRegisteredTaxes() []InterfaceTax {
	return registeredTaxes
}

// GetRegisteredDiscounts returns list of registered tax calculators
func GetRegisteredDiscounts() []InterfaceDiscount {
	return registeredDiscounts
}
