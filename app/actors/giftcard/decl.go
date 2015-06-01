// Package giftcard creates and manage giftCards
package giftcard

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstSessionKeyAppliedGiftCardCodes = "applied_giftcard_codes"
	ConstCollectionNameGiftCard         = "gift_card"

	ConstConfigPathDiscounts   = "general.discounts"
	ConstConfigPathGiftEmail   = "general.discounts.giftCard_email"
	ConstConfigPathGiftCardSKU = "general.discounts.giftCard_SKU_code"

	ConstErrorModule = "giftcard"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstGiftCardStatusNew      = "new"
	ConstGiftCardStatusApplied  = "applied"
	ConstGiftCardStatusUsed     = "used"
	ConstGiftCardStatusOverUsed = "over-used"
	ConstGiftCardStatusRefilled = "refilled"
	ConstGiftCardStatusCanceled = "canceled"
)

// DefaultDiscount is a default implementer of InterfaceDiscount
type DefaultDiscount struct{}
