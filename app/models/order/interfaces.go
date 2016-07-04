// Package order represents abstraction of business layer purchase order object
package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstModelNameOrder           = "Order"
	ConstModelNameOrderCollection = "OrderCollection"

	ConstModelNameOrderItemCollection = "OrderItemCollection"

	ConstOrderStatusDeclined  = "declined"  // order was created and then declined
	ConstOrderStatusNew       = "new"       // order created but not paid
	ConstOrderStatusPending   = "pending"   // order was submitted and currently in processing
	ConstOrderStatusProcessed = "processed" // order was authorized and funds collected
	ConstOrderStatusCompleted = "completed" // order was completed by retailer
	ConstOrderStatusCancelled = "cancelled" // order was cancelled by retailer

	ConstErrorModule = "order"
	ConstErrorLevel  = env.ConstErrorLevelModel
)

// InterfaceOrderItem represents interface to access business layer implementation of purchase order item object
type InterfaceOrderItem interface {
	GetID() string
	SetID(newID string) error

	GetProductID() string

	GetName() string
	GetSku() string

	GetQty() int

	GetPrice() float64

	GetWeight() float64

	GetOptions() map[string]interface{}

	GetSelectedOptions(asLabels bool) map[string]interface{}

	models.InterfaceObject
}

// InterfaceOrder represents interface to access business layer implementation of purchase order object
type InterfaceOrder interface {
	GetItems() []InterfaceOrderItem

	AddItem(productID string, qty int, productOptions map[string]interface{}) (InterfaceOrderItem, error)
	RemoveItem(itemIdx int) error

	CalculateTotals() error

	NewIncrementID() error

	GetIncrementID() string
	SetIncrementID(incrementID string) error

	GetSubtotal() float64
	GetGrandTotal() float64

	GetDiscountAmount() float64
	GetTaxAmount() float64
	GetShippingAmount() float64

	GetTaxes() []StructTaxRate
	GetDiscounts() []StructDiscount

	GetShippingAddress() visitor.InterfaceVisitorAddress
	GetBillingAddress() visitor.InterfaceVisitorAddress

	GetShippingMethod() string
	GetPaymentMethod() string

	GetStatus() string
	SetStatus(status string) error

	Proceed() error
	Rollback() error

	DuplicateOrder(params map[string]interface{}) (interface{}, error)
	SendShippingStatusUpdateEmail() error
	SendOrderConfirmationEmail() error

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
	models.InterfaceListable
}

// InterfaceOrderCollection represents interface to access business layer implementation of purchase order collection
type InterfaceOrderCollection interface {
	ListOrders() []InterfaceOrder

	models.InterfaceCollection
}

// InterfaceOrderItemCollection represents interface to access business layer implementation of purchase order item collection
type InterfaceOrderItemCollection interface {
	models.InterfaceCollection
}

// StructTaxRate represents type to hold tax rate information generated by implementation of InterfaceTax
type StructTaxRate struct {
	Name   string
	Code   string
	Amount float64
}

// StructDiscount represents type to hold discount information generated by implementation of InterfaceDiscount
type StructDiscount struct {
	Name   string
	Code   string
	Amount float64
}
