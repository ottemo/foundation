// Package saleprice is an implementation of discount interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package saleprice

import (
	"time"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/app/models/product"
)

// Package global constants
const (
	ConstConfigPathGroup                  = "general.sale_price"
	ConstConfigPathEnabled                = "general.sale_price.enabled"
	ConstConfigPathSalePriceApplyPriority = "general.sale_price.priority"

	ConstErrorModule = "saleprice"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultSalePrice is an implementer of InterfaceDiscount
type DefaultSalePrice struct {
	id string

	amount        float64
	endDatetime   time.Time
	productID     string
	startDatetime time.Time
}

// DefaultSalePriceCollection is a default implementer of InterfaceSalePriceCollection
type DefaultSalePriceCollection struct {
	listCollection     db.InterfaceDBCollection
	listExtraAtributes []string
}

// SalePriceDelegate type implements InterfaceAttributesDelegate and have handles
// on InterfaceStorable methods which should have call-back on model method call
// in order to test it we are pushing the callback status to model instance
type SalePriceDelegate struct {
	productInstance product.InterfaceProduct
	SalePrices      []saleprice.InterfaceSalePrice
}

// salePriceDelegate variable that is currently used as a sale price delegate to extend product attributes
var salePriceDelegate models.InterfaceAttributesDelegate
