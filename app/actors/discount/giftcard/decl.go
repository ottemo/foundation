// Package giftcard creates and manage gift cards
package giftcard

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstSessionKeyAppliedGiftCardCodes = "applied_giftcard_codes"
	ConstCollectionNameGiftCard         = "gift_card"

	ConstConfigPathGiftEmailTemplate = "general.discounts.giftCard_email"
	ConstConfigPathGiftEmailSubject  = "general.discounts.giftCard_email_subject"
	ConstConfigPathGiftCardSKU       = "general.discounts.giftCard_SKU_code"

	ConstConfigPathGiftCardApplyPriority = "general.discounts.giftCard_apply_priority"

	ConstErrorModule = "giftcard"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstGiftCardStatusNew          = "new"
	ConstGiftCardStatusApplied      = "applied"
	ConstGiftCardStatusUsed         = "used"
	ConstGiftCardStatusOverCredited = "negative"
	ConstGiftCardStatusRefilled     = "refilled"
	ConstGiftCardStatusCancelled    = "cancelled"
	ConstGiftCardStatusDelivered    = "delivered"
)

// DefaultGiftcard is a default implementer of InterfaceDiscount
type DefaultGiftcard struct{}

type GiftcardShipping struct{}
