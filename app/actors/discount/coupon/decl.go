// Package coupon is a default implementation of discount interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package coupon

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstSessionKeyCurrentRedemptions  = "current_redemption_codes"
	ConstCollectionNameCouponDiscounts = "coupon_discounts"

	ConstConfigPathDiscounts             = "general.discounts"
	ConstConfigPathDiscountApplyPriority = "general.discounts.discount_apply_priority"

	ConstErrorModule = "coupon"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// Coupon is a default implementer of InterfaceDiscount
type Coupon struct{}

// usedCoupons contains used coupon codes with visitorsId's, initialize from orders and updated on checkout success
var usedCoupons map[string][]string

type discount struct {
	Code     string
	Name     string
	Total    float64
	Amount   float64
	Percents float64
	Qty      int
}
