package coupon

import (
	"time"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// GetName returns name of current discount implementation
func (it *DefaultDiscount) GetName() string {
	return "CouponDiscount"
}

// GetCode returns code of current discount implementation
func (it *DefaultDiscount) GetCode() string {
	return "coupon_discount"
}

// CalculateDiscount calculates and returns a set of discounts applied to given checkout
func (it *DefaultDiscount) CalculateDiscount(checkoutInstance checkout.InterfaceCheckout) []checkout.StructDiscount {

	var result []checkout.StructDiscount

	// checking session for applied coupon codes
	if currentSession := checkoutInstance.GetSession(); currentSession != nil {

		appliedCodes := utils.InterfaceToStringArray(currentSession.Get(ConstSessionKeyAppliedDiscountCodes))
		usedCodes := utils.InterfaceToStringArray(currentSession.Get(ConstSessionKeyUsedDiscountCodes))

		if len(appliedCodes) > 0 {

			// loading information about applied discounts
			collection, err := db.GetCollection(ConstCollectionNameCouponDiscounts)
			if err != nil {
				return result
			}
			err = collection.AddFilter("code", "in", appliedCodes)
			if err != nil {
				return result
			}

			records, err := collection.Load()
			if err != nil || len(records) == 0 {
				return result
			}

			// making coupon code map for right apply order ignoring used coupons
			discountCodes := make(map[string]map[string]interface{})
			for _, record := range records {
				if discountCode := utils.InterfaceToString(record["code"]); discountCode != "" &&
					!utils.IsInArray(discountCode, usedCodes) && !isLimited(checkoutInstance, record) {
					discountCodes[discountCode] = record
				}
			}

			priorityValue := utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathDiscountApplyPriority))
			checkoutSubtotal := checkoutInstance.GetSubtotal()
			var biggestDiscountApplied float64

			// applying coupon codes
			for appliedCodesIdx, discountCode := range appliedCodes {
				if discountCoupon, ok := discountCodes[discountCode]; ok {

					applyTimes := utils.InterfaceToInt(discountCoupon["times"])
					workSince := utils.InterfaceToTime(discountCoupon["since"])
					workUntil := utils.InterfaceToTime(discountCoupon["until"])

					currentTime := time.Now()

					// to be applicable coupon should satisfy following conditions:
					//   [applyTimes] should be -1 or >0 and [workSince] >= currentTime <= [workUntil] if set
					if (applyTimes == -1 || applyTimes > 0) &&
						(utils.IsZeroTime(workSince) || workSince.Unix() <= currentTime.Unix()) &&
						(utils.IsZeroTime(workUntil) || workUntil.Unix() >= currentTime.Unix()) {

						// calculating coupon discount amount
						discountAmount := utils.InterfaceToFloat64(discountCoupon["amount"])
						discountPercent := utils.InterfaceToFloat64(discountCoupon["percent"])

						possibleDiscount := discountAmount + (discountPercent / 100 * checkoutSubtotal)

						// only the biggest coupon discount will be returned
						if possibleDiscount > biggestDiscountApplied {
							biggestDiscountApplied = possibleDiscount

							result = []checkout.StructDiscount{}

							if discountPercent > 0 {
								result = append(result, checkout.StructDiscount{
									Name:      utils.InterfaceToString(discountCoupon["name"]),
									Code:      utils.InterfaceToString(discountCoupon["code"]),
									Amount:    discountPercent,
									IsPercent: true,
									Priority:  priorityValue,
								})
								priorityValue += float64(0.0001)
							}

							if discountAmount > 0 {
								result = append(result, checkout.StructDiscount{
									Name:      utils.InterfaceToString(discountCoupon["name"]),
									Code:      utils.InterfaceToString(discountCoupon["code"]),
									Amount:    discountAmount,
									IsPercent: false,
									Priority:  priorityValue,
								})
								priorityValue += float64(0.0001)
							}
						}

					} else {
						// we have not applicable coupon - removing it from applied coupons list
						newAppliedCodes := make([]string, 0, len(appliedCodes)-1)
						for idx, value := range appliedCodes {
							if idx != appliedCodesIdx {
								newAppliedCodes = append(newAppliedCodes, value)
							}
						}
						currentSession.Set(ConstSessionKeyAppliedDiscountCodes, newAppliedCodes)
					}
				}
			}
		}
	}

	return result
}

// checks discount limiting parameters and correspondence of current checkout to their values
func isLimited(checkoutInstance checkout.InterfaceCheckout, couponDiscount map[string]interface{}) bool {

	if limits, present := couponDiscount["limits"]; present {
		limitations := utils.InterfaceToMap(limits)
		if len(limitations) > 0 {

			var productsIDs []string
			for _, productInCart := range checkoutInstance.GetCart().GetItems() {
				productsIDs = append(productsIDs, productInCart.GetProductID())
			}

			for key, limit := range limitations {

				switch strings.ToLower(key) {
				case "product_in_cart":
					allowedProducts := utils.InterfaceToArray(limit)
					for index, productID := range allowedProducts {
						if utils.IsInArray(productID, productsIDs) {
							break
						}
						if index == (len(allowedProducts) - 1) {
							return true
						}

					}
				}
			}
		}
	}

	return false
}
