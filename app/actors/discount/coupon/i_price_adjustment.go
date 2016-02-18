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
	return "Coupon"
}

// GetCode returns the code of the current coupon implementation
func (it *Coupon) GetCode() string {
	return "coupon"
}

// CalculateDiscount calculates and returns a set of coupons applied to the provided checkout
// flow:
// add new variable above the loop for coupons - hold items - amount of discount to prevent from staking
// when discount is applied to cart, just add it to array
// when there some product target:
// take a look into var for value and compare if this bigger then update it
// two possible value - percent and amount we should fond a sum for them to compare amounts
func (it *Coupon) Calculate(checkoutInstance checkout.InterfaceCheckout) []checkout.StructPriceAdjustment {

	var result []checkout.StructPriceAdjustment

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

			applicableProductDiscounts := make(map[string][]discount)
			// collect products to one map, that holds productID: qty and used to get apply qty
			productsInCart := make(map[string]int)
			for _, productInCart := range checkoutInstance.GetItems() {
				productID := productInCart.GetProductID()
				productQty := productInCart.GetQty()

				if qty, present := productsInCart[productID]; present {
					productsInCart[productID] = qty + productQty
					continue
				}
				productsInCart[productID] = productQty
				applicableProductDiscounts[productID] = make([]discount, 0)
			}

			// use coupon map to hold the correct application order and ignore previously used coupons
			discountCodes := make(map[string]map[string]interface{})
			for _, record := range records {

				discountsUsageQty := getCouponApplyQty(productsInCart, record)
				discountCode := utils.InterfaceToString(record["code"])

				if discountCode != "" && discountsUsageQty > 0 {
					record["usage_qty"] = discountsUsageQty
					discountCodes[discountCode] = record
				}
			}

			var discountableCartTotal float64
			cartTotal := checkoutInstance.GetItemTotals(0)
			if value, present := cartTotal[checkout.ConstLabelSubtotal]; present {
				discountableCartTotal = value
			}

			couponPriorityValue := utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathDiscountApplyPriority))

			// accumulation of coupon discounts for cart to result and for products to applicableProductDiscounts
			for appliedCodesIdx, discountCode := range redeemedCodes {
				discountCoupon, present := discountCodes[discountCode]
				if !present {
					continue
				}

				validStart := isValidStart(discountCoupon["since"])
				validEnd := isValidEnd(discountCoupon["until"])

				// to be applicable coupon should satisfy following conditions:
				//   [begin] >= currentTime <= [end] if set
				if !validStart || !validEnd {
					// we have not applicable coupon - removing it from applied coupons list
					newRedemptions := make([]string, 0, len(redeemedCodes)-1)
					for idx, value := range redeemedCodes {
						if idx != appliedCodesIdx {
							newRedemptions = append(newRedemptions, value)
						}
					}
					currentSession.Set(ConstSessionKeyCurrentRedemptions, newRedemptions)
				}

				// calculating coupon discount amount
				discountAmount := utils.InterfaceToFloat64(discountCoupon["amount"])
				discountPercent := utils.InterfaceToFloat64(discountCoupon["percent"])

				discountTarget := utils.InterfaceToString(discountCoupon["target"])
				discountUsageQty := utils.InterfaceToInt(discountCoupon["usage_qty"])

				// make some change to name to show values that where used
				discountLabel := getLabel(utils.InterfaceToString(discountCoupon["name"]), discountAmount, discountPercent)
				discountCode := utils.InterfaceToString(discountCoupon["code"])

				// case it's a cart discount we just add them to result with calculating amount based on current totals
				if strings.Contains(discountTarget, checkout.ConstDiscountObjectCart) || discountTarget == "" {
					couponPriorityValue += float64(0.0001)
					priceAdjustmentAmount := discountAmount + discountPercent*discountableCartTotal

					// build price adjustment for cart coupon discount,
					currentPriceAdjustment := checkout.StructPriceAdjustment{
						Code:      discountCode,
						Label:     discountLabel,
						Amount:    priceAdjustmentAmount * -1,
						IsPercent: false,
						Priority:  couponPriorityValue,
						Types:     []string{checkout.ConstLabelSubtotal},
						PerItem:   make(map[int]float64),
					}

					result = append(result, currentPriceAdjustment)

					continue
				}

				// add discount object for every product id that it can affect
				applicableProductDiscount := discount{
					Code:     discountCode,
					Label:    discountLabel,
					Amount:   discountAmount,
					Percents: discountPercent,
					Qty:      discountUsageQty,
				}

				// collect only discounts for productIDs that are in cart
				// applicableProductDiscounts already have all pids keys that are exist
				for _, productID := range utils.InterfaceToArray(discountTarget) {
					if discounts, present := applicableProductDiscounts[productID]; present {
						applicableProductDiscounts[productID] = append(discounts, applicableProductDiscount)
					}
				}
			}

			// adding to discounts the biggest applicable discount per product
			for _, cartItem := range checkoutInstance.GetItems() {

				if cartProduct := cartItem.GetProduct(); cartProduct != nil {
					cartProduct.ApplyOptions(cartItem.GetOptions())
					productPrice := cartProduct.GetPrice()
					productID := cartItem.GetProductID()
					productQty := cartItem.GetQty()

					// discount will be applied for every single product and grouped per item
					for i := 0; i < productQty; i++ {
						if productDiscounts, present := applicableProductDiscounts[productID]; present && len(productDiscounts) > 0 {
							var biggestAppliedDiscount discount
							var biggestAppliedDiscountIndex int

							for index, productDiscount := range productDiscounts {
								productDiscountableAmount := productDiscount.Amount + productPrice*productDiscount.Percents/100

								if biggestAppliedDiscount.Amount < productDiscountableAmount {
									biggestAppliedDiscount = productDiscount
									biggestAppliedDiscount.Amount = productDiscountableAmount
									biggestAppliedDiscountIndex = index
								}
							}

							var newProductDiscounts []discount
							// remove biggest discount from the list
							for index, productDiscount := range productDiscounts {
								if index != biggestAppliedDiscountIndex {
									newProductDiscounts = append(newProductDiscounts, productDiscount)
								}
							}

							applicableProductDiscounts[productID] = newProductDiscounts

							// TODO: place logic to add discount to Price adjustment using
						}
					}

					currentPriceAdjustment := checkout.StructPriceAdjustment{
						Code:      discountCode,
						Label:     discountLabel,
						Amount:    priceAdjustmentAmount * -1,
						IsPercent: false,
						Priority:  couponPriorityValue,
						Types:     []string{checkout.ConstLabelSubtotal},
						PerItem:   make(map[int]float64),
					}

					result = append(result, currentPriceAdjustment)
				}

				// attach price adjustment value to result
			}
		}
	}

	return result
}

// check coupon limitation parameters for correspondence to current checkout values
// return qty of usages if coupon is allowed for current checkout and satisfies all conditions
func getCouponApplyQty(productsInCart map[string]int, couponDiscount map[string]interface{}) int {

	result := -1
	if limits, present := couponDiscount["limits"]; present {
		limitations := utils.InterfaceToMap(limits)
		if len(limitations) > 0 {
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

// validStart returns a boolean value of the datetame passed is valid
func isValidStart(start interface{}) bool {

	couponStart := utils.InterfaceToTime(start)
	currentTime := time.Now()

	isValidStart := (utils.IsZeroTime(couponStart) || couponStart.Unix() <= currentTime.Unix())

	return isValidStart
}

// validEnd returns a boolean value of the datetame passed is valid
func isValidEnd(end interface{}) bool {

	couponEnd := utils.InterfaceToTime(end)
	currentTime := time.Now()

	// to be applicable coupon should satisfy following conditions:
	isValidEnd := (utils.IsZeroTime(couponEnd) || couponEnd.Unix() >= currentTime.Unix())

	return isValidEnd
}

func getLabel(code string, dollarAmount float64, percentAmount float64) string {

	return ""
}
