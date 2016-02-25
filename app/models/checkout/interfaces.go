// Package checkout represents abstraction of business layer checkout object
package checkout

import (
	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
)

// InterfaceCheckout represents interface to access business layer implementation of checkout object
type InterfaceCheckout interface {
	SetShippingAddress(address visitor.InterfaceVisitorAddress) error
	GetShippingAddress() visitor.InterfaceVisitorAddress

	SetBillingAddress(address visitor.InterfaceVisitorAddress) error
	GetBillingAddress() visitor.InterfaceVisitorAddress

	SetPaymentMethod(paymentMethod InterfacePaymentMethod) error
	GetPaymentMethod() InterfacePaymentMethod

	SetInfo(key string, value interface{}) error
	GetInfo(key string) interface{}

	SetShippingMethod(shippingMethod InterfaceShippingMethod) error
	GetShippingMethod() InterfaceShippingMethod

	SetShippingRate(shippingRate StructShippingRate) error
	GetShippingRate() *StructShippingRate

	GetItems() []cart.InterfaceCartItem

	GetItemSpecificTotal(idx int, label string) float64

	GetPriceAdjustments(label string) []StructPriceAdjustment

	GetTaxes() []StructPriceAdjustment
	GetTaxAmount() float64

	GetDiscounts() []StructPriceAdjustment
	GetDiscountAmount() float64

	GetSubtotal() float64
	GetShippingAmount() float64

	CalculateAmount(calculateTarget float64) float64
	GetGrandTotal() float64

	SetCart(checkoutCart cart.InterfaceCart) error
	GetCart() cart.InterfaceCart

	SetVisitor(checkoutVisitor visitor.InterfaceVisitor) error
	GetVisitor() visitor.InterfaceVisitor

	SetSession(api.InterfaceSession) error
	GetSession() api.InterfaceSession

	SetOrder(checkoutOrder order.InterfaceOrder) error
	GetOrder() order.InterfaceOrder

	CheckoutSuccess(checkoutOrder order.InterfaceOrder, session api.InterfaceSession) error
	SendOrderConfirmationMail() error

	IsSubscription() bool

	Submit() (interface{}, error)

	SubmitFinish(map[string]interface{}) (interface{}, error)

	models.InterfaceModel
	models.InterfaceObject
}

// InterfaceShippingMethod represents interface to access business layer implementation of checkout shipping method
type InterfaceShippingMethod interface {
	GetName() string
	GetCode() string

	IsAllowed(checkoutInstance InterfaceCheckout) bool

	GetRates(checkoutInstance InterfaceCheckout) []StructShippingRate
}

// InterfacePaymentMethod represents interface to access business layer implementation of checkout payment method
type InterfacePaymentMethod interface {
	GetName() string
	GetCode() string
	GetType() string

	IsAllowed(checkoutInstance InterfaceCheckout) bool
	IsTokenable(checkoutInstance InterfaceCheckout) bool

	Authorize(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error)
	Capture(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error)
	Refund(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error)
	Void(orderInstance order.InterfaceOrder, paymentInfo map[string]interface{}) (interface{}, error)
}

// InterfacePriceAdjustment represents interface to access business layer implementation of checkout calculation elements
type InterfacePriceAdjustment interface {
	GetName() string
	GetCode() string
	GetPriority() []float64

	Calculate(checkoutInstance InterfaceCheckout) []StructPriceAdjustment
}

// StructShippingRate represents type to hold shipping rate information generated by implementation of InterfaceShippingMethod
type StructShippingRate struct {
	Name  string
	Code  string
	Price float64
}

// StructPriceAdjustment represents type to hold  information generated by implementation of InterfacePriceAdjustment (calculating entities of checkout)
type StructPriceAdjustment struct {
	Code      string          `json:"Code"`
	Label     string          `json:"Label"`
	Priority  float64         `json:"Priority"`
	Amount    float64         `json:"Amount"`
	IsPercent bool            `json:"IsPercent,bool"`
	Types     []string        `json:"Types"`
	PerItem   map[int]float64 `json:"PerItem,string"`
}
