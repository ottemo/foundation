package checkout

import (
	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/app/models/visitor"
)

type I_Checkout interface {
	SetShippingAddress(address visitor.I_VisitorAddress) error
	GetShippingAddress() visitor.I_VisitorAddress

	SetBillingAddress(address visitor.I_VisitorAddress) error
	GetBillingAddress() visitor.I_VisitorAddress

	SetPaymentMethod(paymentMethod I_PaymentMethod) error
	GetPaymentMethod() I_PaymentMethod

	SetInfo(key string, value interface{}) error
	GetInfo(key string) interface{}

	SetShippingMethod(shippingMethod I_ShippingMehod) error
	GetShippingMethod() I_ShippingMehod

	SetShippingRate(shippingRate T_ShippingRate) error
	GetShippingRate() *T_ShippingRate

	GetTaxes() (float64, []T_TaxRate)
	GetDiscounts() (float64, []T_Discount)

	GetGrandTotal() float64

	SetCart(checkoutCart cart.I_Cart) error
	GetCart() cart.I_Cart

	SetVisitor(checkoutVisitor visitor.I_Visitor) error
	GetVisitor() visitor.I_Visitor

	SetSession(api.I_Session) error
	GetSession() api.I_Session

	SetOrder(checkoutOrder order.I_Order) error
	GetOrder() order.I_Order

	SendOrderConfirmationMail() error

	Submit() (interface{}, error)

	models.I_Model
}

type I_ShippingMehod interface {
	GetName() string
	GetCode() string

	IsAllowed(checkoutInstance I_Checkout) bool

	GetRates(checkoutInstance I_Checkout) []T_ShippingRate
}

type I_PaymentMethod interface {
	GetName() string
	GetCode() string
	GetType() string

	IsAllowed(checkoutInstance I_Checkout) bool

	Authorize(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error)
	Capture(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error)
	Refund(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error)
	Void(orderInstance order.I_Order, paymentInfo map[string]interface{}) (interface{}, error)
}

type I_Tax interface {
	GetName() string
	GetCode() string

	CalculateTax(checkoutInstance I_Checkout) []T_TaxRate
}

type I_Discount interface {
	GetName() string
	GetCode() string

	CalculateDiscount(checkoutInstance I_Checkout) []T_Discount
}

type T_ShippingRate struct {
	Name  string
	Code  string
	Price float64
}

type T_TaxRate struct {
	Name   string
	Code   string
	Amount float64
}

type T_Discount struct {
	Name   string
	Code   string
	Amount float64
}
