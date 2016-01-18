package giftcard

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetName returns name of current discount implementation
func (it *DefaultGiftcard) GetName() string {
	return "GiftCardDiscount"
}

// GetCode returns code of current discount implementation
func (it *DefaultGiftcard) GetCode() string {
	return "giftcard_discount"
}

// CalculateDiscount calculates and returns amount and set of applied gift card discounts to given checkout
func (it *DefaultGiftcard) CalculateDiscount(checkoutInstance checkout.InterfaceCheckout) []checkout.StructDiscount {
	var result []checkout.StructDiscount

	// checking session for applied gift cards codes
	if currentSession := checkoutInstance.GetSession(); currentSession != nil {

		appliedCodes := utils.InterfaceToStringArray(currentSession.Get(ConstSessionKeyAppliedGiftCardCodes))

		if len(appliedCodes) > 0 {

			// loading information about applied discounts
			collection, err := db.GetCollection(ConstCollectionNameGiftCard)
			if err != nil {
				return result
			}
			err = collection.AddFilter("code", "in", appliedCodes)
			if err != nil {
				return result
			}

			records, err := collection.Load()
			if err != nil {
				return result
			}

			// making gift cards code map for right apply order
			giftCardCodes := make(map[string]map[string]interface{})

			for _, record := range records {
				if giftCardCode := utils.InterfaceToString(record["code"]); giftCardCode != "" {
					giftCardCodes[giftCardCode] = record
				}
			}
			priorityValue := utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathGiftCardApplyPriority))

			// applying gift card discount codes
			for _, giftCardCode := range appliedCodes {
				if giftCard, ok := giftCardCodes[giftCardCode]; ok {

					giftCardAmount := utils.InterfaceToFloat64(giftCard["amount"])

					// to be applicable gift card should satisfy following conditions:
					// have positive amount and we have amount that will be discounted
					result = append(result, checkout.StructDiscount{
						Name:      utils.InterfaceToString(giftCard["name"]),
						Code:      utils.InterfaceToString(giftCard["code"]),
						Amount:    giftCardAmount,
						IsPercent: false,
						Priority:  priorityValue,
						Object:    checkout.ConstDiscountObjectCart,
						Type:      it.GetCode(),
					})
					priorityValue += float64(0.0001)

				}
			}
		}
	}

	return result
}
