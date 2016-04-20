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

// GetPriority returns the code of the current coupon implementation
func (it *DefaultGiftcard) GetPriority() []float64 {
	// adding this first value of priority to make PA that will reduce GT of gift cards by 100% right after subtotal calculation
	return []float64{checkout.ConstCalculateTargetSubtotal, utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathGiftCardApplyPriority)), checkout.ConstCalculateTargetGrandTotal}
}

// Calculate calculates and returns amount and set of applied gift card discounts to given checkout
func (it *DefaultGiftcard) Calculate(checkoutInstance checkout.InterfaceCheckout, currentPriority float64) []checkout.StructPriceAdjustment {
	var result []checkout.StructPriceAdjustment

	if currentPriority == checkout.ConstCalculateTargetSubtotal {
		// discount gift cards on 100%
		return result
	}
	if currentPriority == checkout.ConstCalculateTargetGrandTotal {
		// return gift cards amount
		return result
	}

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
				giftCard, present := giftCardCodes[giftCardCode]
				if !present {
					continue
				}

				giftCardAmount := utils.InterfaceToFloat64(giftCard["amount"])

				if giftCardAmount > 0 {
					result = append(result, checkout.StructPriceAdjustment{
						Code:      utils.InterfaceToString(giftCard["code"]),
						Name:      utils.InterfaceToString(giftCard["name"]),
						Amount:    giftCardAmount * -1,
						IsPercent: false,
						Priority:  priorityValue,
						Labels:    []string{checkout.ConstLabelGiftCard},
					})
					priorityValue += float64(0.0001)
				}
			}
		}
	}

	return result
}
