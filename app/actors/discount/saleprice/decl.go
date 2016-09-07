// Package saleprice is an implementation of discount interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package saleprice

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
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
	SalePrices      []map[string]interface{}
}

// salePriceDelegate variable that is currently used as a stock delegate to extend product attributes
var salePriceDelegate models.InterfaceAttributesDelegate
