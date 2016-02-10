package coupon

import (
	"strings"
	"time"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetName returns the name of the current coupon implementation
func (it *Coupon) GetName() string {
	return "CouponDiscount"
}

// GetCode returns the code of the current coupon implementation
func (it *Coupon) GetCode() string {
	return "coupon_discount"
}

// CalculateDiscount calculates and returns a set of coupons applied to the provided checkout
func (it *Coupon) CalculateDiscount(checkoutInstance checkout.InterfaceCheckout) []checkout.StructDiscount {

	var result []checkout.StructDiscount

	// check session for applied coupon codes
	if currentSession := checkoutInstance.GetSession(); currentSession != nil {

		redeemedCodes := utils.InterfaceToStringArray(currentSession.Get(ConstSessionKeyCurrentRedemptions))

		if len(redeemedCodes) > 0 {

			// loading information about applied coupons
			collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
			if err != nil {
				return result
			}
			err = collection.AddFilter("code", "in", redeemedCodes)
			if err != nil {
				return result
			}

			records, err := collection.Load()
			if err != nil || len(records) == 0 {
				return result
			}

			// use coupon map to hold the correct application order and ignore previously used coupons
			discountCodes := make(map[string]map[string]interface{})
			for _, record := range records {

				discountsUsageQty := discountsUsage(checkoutInstance, record)
				discountCode := utils.InterfaceToString(record["code"])

				if discountCode != "" && discountsUsageQty > 0 {
					record["usage_qty"] = discountsUsageQty
					discountCodes[discountCode] = record
				}
			}

			discountPriorityValue := utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathDiscountApplyPriority))

			// accumulation of coupon discounts to result
			for appliedCodesIdx, discountCode := range redeemedCodes {
				if discountCoupon, ok := discountCodes[discountCode]; ok {

					couponStart := utils.InterfaceToTime(discountCoupon["since"])
					couponEnd := utils.InterfaceToTime(discountCoupon["until"])

					currentTime := time.Now()

					// to be applicable coupon should satisfy following conditions:
					//   [begin] >= currentTime <= [end] if set
					isValidStart := (utils.IsZeroTime(couponStart) || couponStart.Unix() <= currentTime.Unix())
					isValidEnd := (utils.IsZeroTime(couponEnd) || couponEnd.Unix() >= currentTime.Unix())

					if isValidStart && isValidEnd {

						// calculating coupon discount amount
						discountAmount := utils.InterfaceToFloat64(discountCoupon["amount"])
						discountPercent := utils.InterfaceToFloat64(discountCoupon["percent"])
						discountTarget := utils.InterfaceToString(discountCoupon["target"])
						discountUsageQty := utils.InterfaceToInt(discountCoupon["usage_qty"])

						currentDiscount := checkout.StructDiscount{
							Name:      utils.InterfaceToString(discountCoupon["name"]),
							Code:      utils.InterfaceToString(discountCoupon["code"]),
							Amount:    discountPercent,
							IsPercent: true,
							Priority:  discountPriorityValue,
							Object:    checkout.ConstDiscountObjectCart,
							Type:      it.GetCode(),
						}

						// case it's a cart discount we just add them to result
						if strings.Contains(discountTarget, checkout.ConstDiscountObjectCart) || discountTarget == "" {
							discountPercent = discountPercent * float64(discountUsageQty)
							discountAmount = discountAmount * float64(discountUsageQty)

							if discountPercent > 0 {
								currentDiscount.Amount = discountPercent
								currentDiscount.Priority = discountPriorityValue + float64(0.01)
								result = append(result, currentDiscount)
								discountPriorityValue += float64(0.0001)
							}

							if discountAmount > 0 {
								currentDiscount.Amount = discountAmount
								currentDiscount.IsPercent = false
								currentDiscount.Priority = discountPriorityValue + float64(0.01)
								result = append(result, currentDiscount)
								discountPriorityValue += float64(0.0001)
							}

							continue
						}

						// parse target as array of productIDs on which we will return discounts
						if discountPercent > 0 {
							currentDiscount.Amount = discountPercent
							currentDiscount.IsPercent = true

							for _, productID := range utils.InterfaceToArray(discountTarget) {
								currentDiscount.Object = utils.InterfaceToString(productID)
								for index := 0; index < discountUsageQty; index++ {
									result = append(result, currentDiscount)
								}
							}
							discountPriorityValue += float64(0.0001)
						}

						if discountAmount > 0 {
							currentDiscount.Amount = discountAmount
							currentDiscount.IsPercent = false
							currentDiscount.Priority = discountPriorityValue

							for _, productID := range utils.InterfaceToArray(discountTarget) {
								currentDiscount.Object = utils.InterfaceToString(productID)
								for index := 0; index < discountUsageQty; index++ {
									result = append(result, currentDiscount)
								}
							}
							discountPriorityValue += float64(0.0001)
						}

					} else {
						// we have not applicable coupon - removing it from applied coupons list
						newRedemptions := make([]string, 0, len(redeemedCodes)-1)
						for idx, value := range redeemedCodes {
							if idx != appliedCodesIdx {
								newRedemptions = append(newRedemptions, value)
							}
						}
						currentSession.Set(ConstSessionKeyCurrentRedemptions, newRedemptions)
					}
				}
			}
		}
	}

	return result
}

// check coupon limitation parameters for correspondence to current checkout values
// return qty of usages if coupon is allowed for current checkout and satisfies all conditions
func discountsUsage(checkoutInstance checkout.InterfaceCheckout, couponDiscount map[string]interface{}) int {

	result := -1
	if limits, present := couponDiscount["limits"]; present {
		limitations := utils.InterfaceToMap(limits)
		if len(limitations) > 0 {

			productsInCart := make(map[string]int)

			// collect products to one map by ID and qty
			if currentCart := checkoutInstance.GetCart(); currentCart != nil {
				for _, productInCart := range currentCart.GetItems() {
					productID := productInCart.GetProductID()
					productQty := productInCart.GetQty()

					if qty, present := productsInCart[productID]; present {
						productsInCart[productID] = qty + productQty
						continue
					}
					productsInCart[productID] = productQty
				}
			}

			for limitingKey, limitingValue := range limitations {

				switch strings.ToLower(limitingKey) {
				case "product_in_cart":
					requiredProduct := utils.InterfaceToStringArray(limitingValue)
					for index, productID := range requiredProduct {
						if _, present := productsInCart[productID]; present {
							break
						}
						if index == (len(requiredProduct) - 1) {
							return 0
						}
					}

				case "products_in_cart":
					requiredProducts := utils.InterfaceToStringArray(limitingValue)
					for _, productID := range requiredProducts {
						if _, present := productsInCart[productID]; !present {
							return 0
						}
					}

				case "products_in_qty":
					requiredProducts := utils.InterfaceToMap(limitingValue)
					for requiredProductID, requiredQty := range requiredProducts {
						requiredQty := utils.InterfaceToInt(requiredQty)
						if requiredQty < 1 {
							requiredQty = 1
						}
						productQty, present := productsInCart[requiredProductID]
						limitingQty := utils.InterfaceToInt(productQty / requiredQty)

						if !present || limitingQty < 1 {
							return 0
						}

						if result == -1 || limitingQty < result {
							result = limitingQty
						}

					}
				case "max_usage_qty":
					if limitingQty := utils.InterfaceToInt(limitingValue); limitingQty >= 1 && limitingQty < result {
						result = limitingQty
					}
				}
			}
		}
	}
	if result == -1 {
		result = 1
	}

	return result
}
