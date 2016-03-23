package shipstation

import "github.com/ottemo/foundation/env"

const (
	ConstErrorModule = "shipstation"
	ConstErrorLevel  = env.ConstErrorLevel
)

// struct goes here
type Orders struct {
	Orders []Order `xml:"Order"`
}

type Order struct {
	OrderId        string
	OrderNumber    string
	OrderDate      string // Set to string because we can't control the date formatting otherwise
	OrderStatus    string
	LastModified   string // Same formatting issue
	OrderTotal     float64
	ShippingAmount float64

	// Containers
	Customer Customer
	Items    []OrderItem `xml:"Items>Item"`
}

type Customer struct {
	CustomerCode string // We use email address here

	// Containers
	BillingAddress  BillingAddress  `xml:"BillTo"`
	ShippingAddress ShippingAddress `xml:"ShipTo"`
}

type BillingAddress struct {
	Name string
}

type ShippingAddress struct {
	Name     string
	Address1 string
	City     string
	State    string
	Country  string
}

type OrderItem struct {
	Sku       string `xml:"SKU"`
	Name      string
	Quantity  int
	UnitPrice float64
}
