package coupon

import (
	"strings"
	"time"

	"fmt"
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

// GetPriority returns the code of the current coupon implementation
func (it *Coupon) GetPriority() []float64 {
	return []float64{utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathDiscountApplyPriority))}
}

// Calculate calculates and returns a set of coupons applied to the provided checkout
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
				fmt.Println(discountCode, discountsUsageQty)

				if discountCode != "" && discountsUsageQty > 0 {
					record["usage_qty"] = discountsUsageQty
					discountCodes[discountCode] = record
				}
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

				discountTarget := utils.InterfaceToString(discountCoupon["target"])

				// add discount object for every product id that it can affect
				applicableDiscount := discount{
					Code:     utils.InterfaceToString(discountCoupon["code"]),
					Label:    utils.InterfaceToString(discountCoupon["name"]),
					Amount:   utils.InterfaceToFloat64(discountCoupon["amount"]),
					Percents: utils.InterfaceToFloat64(discountCoupon["percent"]),
					Qty:      utils.InterfaceToInt(discountCoupon["usage_qty"]),
				}

				// case it's a cart discount we just add them to result with calculating amount based on current totals
				if strings.Contains(discountTarget, checkout.ConstDiscountObjectCart) || discountTarget == "" {
					couponPriorityValue += float64(0.000001)

					// build price adjustment for cart coupon discount,
					// one for percent and one for dollar amount value of coupon
					// TODO: this part should be moved in calculate phase to priority 2.2
					currentPriceAdjustment := checkout.StructPriceAdjustment{
						Code:      applicableDiscount.Code,
						Label:     getLabel(applicableDiscount),
						Amount:    applicableDiscount.Percents * -1,
						IsPercent: true,
						Priority:  couponPriorityValue + float64(0.001),
						Types:     []string{checkout.ConstLabelDiscount},
						PerItem:   make(map[int]float64),
					}

					if applicableDiscount.Percents > 0 {
						currentPriceAdjustment.Priority += float64(0.00001)

						result = append(result, currentPriceAdjustment)
					}

					if applicableDiscount.Amount > 0 {
						currentPriceAdjustment.Amount = applicableDiscount.Amount * -1
						currentPriceAdjustment.IsPercent = false
						currentPriceAdjustment.Priority += float64(0.0001)

						result = append(result, currentPriceAdjustment)
					}

					continue
				}

				// collect only discounts for productIDs that are in cart
				for _, productID := range utils.InterfaceToStringArray(discountTarget) {
					if discounts, present := applicableProductDiscounts[productID]; present {
						applicableProductDiscounts[productID] = append(discounts, applicableDiscount)
					}
				}
			}

			// hold price adjustment for every coupon code ( to make total details with right description)
			priceAdjustments := make(map[string]checkout.StructPriceAdjustment)

			// adding to discounts the biggest applicable discount per product
			for _, cartItem := range checkoutInstance.GetItems() {

				if cartProduct := cartItem.GetProduct(); cartProduct != nil {
					cartProduct.ApplyOptions(cartItem.GetOptions())
					productPrice := cartProduct.GetPrice()
					productID := cartItem.GetProductID()
					productQty := cartItem.GetQty()

					// discount will be applied for every single product and grouped per item
					for i := 0; i < productQty; i++ {
						productDiscounts, present := applicableProductDiscounts[productID]
						if !present || len(productDiscounts) <= 0 {
							break
						}

						var biggestAppliedDiscount discount
						var biggestAppliedDiscountIndex int

						// looking for biggest applicable discount for current item
						for index, productDiscount := range productDiscounts {
							if (productDiscount.Qty) > 0 {
								productDiscountableAmount := productDiscount.Amount + productPrice*productDiscount.Percents/100

								// if we have discount that is bigger then a price we will apply it
								if productDiscountableAmount > productPrice {
									biggestAppliedDiscount = productDiscount
									biggestAppliedDiscount.Total = productPrice
									biggestAppliedDiscountIndex = index
									break
								}

								if biggestAppliedDiscount.Total < productDiscountableAmount {
									biggestAppliedDiscount = productDiscount
									biggestAppliedDiscount.Total = productDiscountableAmount
									biggestAppliedDiscountIndex = index
								}
							}
						}

						// update used discount and change qty of chosen discount to number of usage
						discountUsed := 1
						productDiscounts[biggestAppliedDiscountIndex].Qty--
						for i < productQty && productDiscounts[biggestAppliedDiscountIndex].Qty > 0 {
							i++
							productDiscounts[biggestAppliedDiscountIndex].Qty--
							discountUsed++
						}
						biggestAppliedDiscount.Qty = discountUsed

						// remove fully used discount from discounts list
						var newProductDiscounts []discount
						for _, currentDiscount := range productDiscounts {
							if currentDiscount.Qty > 0 {
								newProductDiscounts = append(newProductDiscounts, currentDiscount)
							}
						}
						applicableProductDiscounts[productID] = newProductDiscounts

						// making from discount price adjustment
						// calculating amount that will be discounted from item
						amount := float64(biggestAppliedDiscount.Qty) * biggestAppliedDiscount.Total * -1

						// add this amount to already existing PA (with the same coupon code) or creating new
						if priceAdjustment, present := priceAdjustments[biggestAppliedDiscount.Code]; present {

							if value, present := priceAdjustment.PerItem[cartItem.GetIdx()]; present {
								amount += value
							}

							priceAdjustment.PerItem[cartItem.GetIdx()] = amount
							priceAdjustment.Label = updateLabel(priceAdjustment.Label, biggestAppliedDiscount)

						} else {
							couponPriorityValue += float64(0.000001)
							priceAdjustments[biggestAppliedDiscount.Code] = checkout.StructPriceAdjustment{
								Code:      biggestAppliedDiscount.Code,
								Label:     getLabel(biggestAppliedDiscount),
								Amount:    0,
								IsPercent: false,
								Priority:  couponPriorityValue,
								Types:     []string{checkout.ConstLabelDiscount},
								PerItem: map[int]float64{
									cartItem.GetIdx(): amount,
								},
							}
						}
					}
				}
			}

			// attach price adjustments on products to result
			for _, priceAdjustment := range priceAdjustments {
				result = append(result, priceAdjustment)
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
				} // end of switch
			} // end of loop

			if maxLimitValue, present := limitations["max_usage_qty"]; present {
				if limitingQty := utils.InterfaceToInt(maxLimitValue); limitingQty >= 1 && result > limitingQty {
					result = limitingQty
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

// getLabel should be used to make good description for label
// example 'name |20%x3&15$x3|'
func getLabel(discount discount) string {

	result := discount.Label
	qty := utils.InterfaceToString(discount.Qty)
	flag := false

	if discount.Percents != 0 {
		flag = true
		result += " |" + utils.InterfaceToString(utils.RoundPrice(discount.Percents)) + "%x" + qty
	}

	if discount.Amount != 0 {
		if flag {
			result += "&"
		} else {
			result += " |"
		}

		result += utils.InterfaceToString(utils.RoundPrice(discount.Amount)) + "$x" + qty
		flag = true
	}

	if flag {
		result += "|"
	}

	return result
}

// updateLabel should update label and attach details of discount
func updateLabel(existingLabel string, discount discount) string {

	if len(existingLabel) == 0 || !strings.Contains(existingLabel, "|") {
		return getLabel(discount)
	}

	labelParts := strings.Split(existingLabel, "|")

	if len(labelParts) >= 3 {
		valuePart := labelParts[len(labelParts)-2]

		if multiplier := strings.Index(valuePart, "x"); multiplier < 0 {
			return getLabel(discount)
		}

		dollarAmount := ""
		percentAmount := ""
		dollarQty := 0
		percentQty := 0

		separator := strings.Index(valuePart, "&")

		if index := strings.Index(valuePart, "$x"); index > 0 {
			dollarQty += utils.InterfaceToInt(valuePart[index:])
			dollarAmount = valuePart[0 : index+2]
			if separator > 0 {
				dollarAmount = valuePart[separator+1 : index+2]
			}
		}

		if index := strings.Index(valuePart, "%x"); index > 0 {
			percentAmount = valuePart[0 : index+2]
			if separator > index {
				percentQty += utils.InterfaceToInt(valuePart[index+2 : separator])
			} else {
				percentQty += utils.InterfaceToInt(valuePart[index+2:])
			}
		}

		if discount.Percents != 0 {
			percentQty += discount.Qty
			if percentAmount == "" {
				percentAmount = utils.InterfaceToString(discount.Percents) + "%x"
			}
		}

		if discount.Amount != 0 {
			dollarQty += discount.Qty
			if dollarAmount == "" {
				dollarAmount = utils.InterfaceToString(discount.Amount) + "$x"
			}
		}

		if percentAmount != "" && dollarAmount != "" {
			labelParts[len(labelParts)-2] = percentAmount + utils.InterfaceToString(percentQty) + "&" + dollarAmount + utils.InterfaceToString(dollarQty)
		} else if percentAmount != "" {
			labelParts[len(labelParts)-2] = percentAmount + utils.InterfaceToString(percentQty)
		} else if dollarAmount != "" {
			labelParts[len(labelParts)-2] = dollarAmount + utils.InterfaceToString(dollarQty)
		}

		return strings.Join(labelParts, "|")
	}

	return getLabel(discount)
}
