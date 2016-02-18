// Package checkout is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package checkout

import (
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "checkout"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultCheckout is a default implementer of InterfaceCheckout
type DefaultCheckout struct {
	CartID    string
	VisitorID string
	OrderID   string

	SessionID string

	ShippingAddress map[string]interface{}
	BillingAddress  map[string]interface{}

	PaymentMethodCode  string
	ShippingMethodCode string

	ShippingRate checkout.StructShippingRate

	PriceAdjustments []checkout.StructPriceAdjustment

	// should store details about applied adjustments for specific keys
	// 0 - cart, 1,2,3, .. n - index of cart item
	calculationDetailTotals map[int]map[string]float64
	cart                    cart.InterfaceCart

	Taxes     []checkout.StructTaxRate
	Discounts []checkout.StructDiscount

	Info map[string]interface{}

	// amounts used to have static values during calculation
	// out of date?
	subtotalAmount  float64
	shippingAmount  float64
	taxesAmount     float64
	discountsAmount float64

	calculateAmount float64

	// flags enables and disables during calculation to prevent recursion
	taxesCalculateFlag     bool
	discountsCalculateFlag bool
	calculateFlag          bool
}
