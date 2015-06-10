// Package order is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/order" package
package order

import (
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/env"
	"sync"
	"time"
)

// Package global variables
var (
	lastIncrementID      int
	lastIncrementIDMutex sync.Mutex
)

// Package global constants
const (
	ConstCollectionNameOrder      = "orders"
	ConstCollectionNameOrderItems = "order_items"

	ConstIncrementIDFormat = "%0.10d"

	ConstConfigPathLastIncrementID = "internal.order.increment_id"

	ConstErrorModule = "order"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultOrderItem is a default implementer of InterfaceOrderItem
type DefaultOrderItem struct {
	id  string
	idx int

	OrderID string

	ProductID string

	Qty int

	Name string
	Sku  string

	ShortDescription string

	Options map[string]interface{}

	Price  float64
	Weight float64
}

// DefaultOrder is a default implementer of InterfaceOrder
type DefaultOrder struct {
	id string

	IncrementID string
	Status      string

	VisitorID string
	CartID    string

	Description  string
	PaymentInfo  map[string]interface{}
	CustomInfo   map[string]interface{}
	ShippingInfo map[string]interface{}

	BillingAddress  map[string]interface{}
	ShippingAddress map[string]interface{}

	CustomerEmail string
	CustomerName  string

	PaymentMethod  string
	ShippingMethod string

	Subtotal       float64
	Discount       float64
	TaxAmount      float64
	ShippingAmount float64
	GrandTotal     float64

	Taxes     []order.StructTaxRate
	Discounts []order.StructDiscount

	Notes []string

	CreatedAt time.Time
	UpdatedAt time.Time

	Items map[int]order.InterfaceOrderItem

	maxIdx int
}

// DefaultOrderCollection is a default implementer of InterfaceOrderCollection
type DefaultOrderCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}

// DefaultOrderItemCollection is a default implementer of InterfaceOrderItemCollection
type DefaultOrderItemCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
