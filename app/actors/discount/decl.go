// Package discount is a default implementation of discount interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package discount

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstSessionKeyAppliedDiscountCodes = "applied_discount_codes"
	ConstCollectionNameCouponDiscounts  = "coupon_discounts"

	ConstErrorModule = "discount"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultDiscount is a default implementer of InterfaceDiscount
type DefaultDiscount struct{}
