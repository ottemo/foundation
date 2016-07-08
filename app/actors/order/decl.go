// Package order is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/order" package
package order

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"

	"sync"
	"time"

	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	lastIncrementID      int
	lastIncrementIDMutex sync.Mutex

	// used when exporting orders in CSV format
	orderFields = [][]string{

		{
			"Order No", "Order Date", "Order Time", "Email", "Payment Status", "Order Status", "Discount", "Sub Total", "Shipping Cost", "Tax", "Total", "Billing FirstName", "Billing LastName", "Billing Company", "Billing Address 1", "Billing Address 2", "Billing City", "Billing State", "Billing Country", "Billing Zip", "Billing Phone", "Shipping FirstName", "Shipping LastName", "Shipping Company", "Shipping Address 1", "Shipping Address 2", "Shipping City", "Shipping State", "Shipping Country", "Shipping Zip", "Shipping Phone", "Shipping Method", "Shipping Carrier", "Payment Method", "Transaction ID", "Credit Card No", "Card Holder Name", "Card Type", "SKU", "Product Name", "Item Price", "Item Quantity", "Item Weight", "Item Amount", "Item Order No",
		},
	}

	dataSet = [][]interface{}{

		{
			"$increment_id",
			func(record map[string]interface{}) string {
				timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))

				createdAt := utils.InterfaceToTime(record["created_at"])
				convertedDate, _ := utils.MakeUTCOffsetTime(createdAt, timeZone)
				return convertedDate.Format("01/02/2006")
			},
			func(record map[string]interface{}) string {
				timeZone := utils.InterfaceToString(env.ConfigGetValue(app.ConstConfigPathStoreTimeZone))

				createdAt := utils.InterfaceToTime(record["created_at"])
				convertedDate, _ := utils.MakeUTCOffsetTime(createdAt, timeZone)
				return convertedDate.Format("15:04:05")
			},

			"$customer_email", "$status", "$status", "$discount", "$subtotal", "$shipping_amount", "$tax_amount", "$grand_total",
			"billing.first_name", "billing.last_name", "billing.company", "billing.address_line1", "billing.address_line2",
			"billing.city", "billing.state", "billing.country", "billing.zip_code", "billing.phone",
			"shipping.first_name", "shipping.last_name", "shipping.company", "shipping.address_line1", "shipping.address_line2",
			"shipping.city", "shipping.state", "shipping.country", "shipping.zip_code", "shipping.phone",

			// "Shipping Method" shipping_info.shipping_method_name
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["shipping_info"])["shipping_method_name"])
			},

			"$shipping_method",

			// "Payment Method", payment_info.payment_method_name
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["payment_info"])["payment_method_name"])
			},

			// "Transaction ID", payment_info.transactionID
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["payment_info"])["transactionID"])
			},

			// "Credit Card No", payment_info.creditCardNumbers
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["payment_info"])["creditCardNumbers"])
			},

			// "Card Holder Name"
			func(record map[string]interface{}) string {
				address := utils.InterfaceToMap(record["billing_address"])
				return utils.InterfaceToString(address["first_name"]) + " " + utils.InterfaceToString(address["last_name"])
			},

			// "Card Type",, payment_info.creditCardType
			func(record map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToMap(record["payment_info"])["creditCardType"])
			},
			"item.sku", "item.name", "item.price", "item.qty", "item.weight",
			//			"item.amount",
			func(orderItemIndex int, orderItem map[string]interface{}) string {
				return utils.InterfaceToString(utils.InterfaceToFloat64(orderItem["price"]) * utils.InterfaceToFloat64(orderItem["qty"]))
			},
			"item.order_id",
		},
	}

	blocksHeaders = map[string][][]string{
		"customers": {{}},
		"invoices":  orderFields,
	}

	blocksRecords = map[string][][]interface{}{
		"customers": {{}},
		"invoices":  dataSet,
	}
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

	SessionID string
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
