// Package saleprice is an implementation of discount interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package saleprice

import (
	"github.com/ottemo/foundation/env"
	"time"
	"github.com/ottemo/foundation/db"
)

// Package global constants
const (
	ConstModelNameSalePrice = ""

	ConstCollectionNameSalePrices = "sale_prices"

	ConstConfigPathSalePriceApplyPriority	= "general.discounts.salePrice_apply_priority"

	ConstErrorModule = "saleprice"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// SalePrice is an implementer of InterfaceDiscount
type DefaultSalePrice struct{
	id	string

	amount		float64
	endDatetime	time.Time
	productId 	string
	startDatetime	time.Time
}


// DefaultSalePriceCollection is a default implementer of InterfaceSalePriceCollection
type DefaultSalePriceCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}
