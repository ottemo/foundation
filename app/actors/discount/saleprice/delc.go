// Package saleprice is an implementation of discount interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package saleprice

import (
	"github.com/ottemo/foundation/env"
	"time"
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

//type discount struct {
//	//Code     string
//	//Name     string
//	//Total    float64
//	//Amount   float64
//	//Percents float64
//	//Qty      int
//}