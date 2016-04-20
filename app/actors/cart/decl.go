// Package cart is a default implementation of interfaces declared in
// "github.com/ottemo/foundation/app/models/cart" package
package cart

import (
	"time"

	"github.com/ottemo/foundation/app/models/cart"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCartCollectionName      = "cart"
	ConstCartItemsCollectionName = "cart_items"

	ConstErrorModule = "cart"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstEventAPIAdd    = "api.cart.addToCart"
	ConstEventAPIUpdate = "api.cart.update"

	ConstConfigPathCartAbandonEmailSendTime = "general.checkout.abandonEmailSendTime"
	ConstConfigPathCartAbandonEmailTemplate = "general.checkout.abandonEmailTemplate"
)

// DefaultCart is a default implementer of InterfaceCart
type DefaultCart struct {
	id string

	VisitorID string
	SessionID string

	Info       map[string]interface{}
	CustomInfo map[string]interface{}
	Items      map[int]cart.InterfaceCartItem

	Active    bool
	UpdatedAt time.Time
	Subtotal  float64

	maxIdx int
}

// DefaultCartItem is a default implementer of InterfaceCart
type DefaultCartItem struct {
	id        string
	idx       int
	ProductID string
	Qty       int
	Options   map[string]interface{}
	product   product.InterfaceProduct

	Cart *DefaultCart
}

// AbandonCartEmailData is a container for carts and visitors who have items in
// their cart and still have a valid session.
type AbandonCartEmailData struct {
	Visitor AbandonVisitor
	Cart    AbandonCart
}

// AbandonVisitor is a struct to hold the info needed to contact a visitor with
// items in their cart who has not checked out yet.
type AbandonVisitor struct {
	Email     string
	FirstName string
	LastName  string
}

// AbandonCart is a struct holding the ID of the abandoned cart.
type AbandonCart struct {
	ID string
	// Items []AbandonCartItem
}

// type AbandonCartItem struct {
// 	Name  string
// 	SKU   string
// 	Price float64
// 	Image string
// }
