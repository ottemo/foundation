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

	ConstOrderStatusCancelled = "cancelled"
	ConstOrderStatusNew       = "new"
	ConstOrderStatusPending   = "pending"
	ConstOrderStatusCompleted = "completed"

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

	GetShippingAddress() visitor.InterfaceVisitorAddress
	GetBillingAddress() visitor.InterfaceVisitorAddress

	GetShippingMethod() string
	GetPaymentMethod() string

	GetStatus() string
	SetStatus(status string) error

	Proceed() error
	Cancel() error

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
