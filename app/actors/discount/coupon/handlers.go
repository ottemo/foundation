package coupon

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// checkoutSuccessHandler find applied coupons for this session and add them to session used coupons values
func checkoutSuccessHandler(event string, eventData map[string]interface{}) bool {

	orderPlaced, ok := eventData["order"].(order.InterfaceOrder)
	if !ok {
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4238e657-ed8e-44ed-89cf-0e66b91fecbd", "order can't be used"))
		return false
	}

	// check is discounts are applied to this order if they make change of session used discounts
	orderAppliedDiscounts := orderPlaced.GetDiscounts()
	if len(orderAppliedDiscounts) > 0 && eventData["session"] != nil {

		session, ok := eventData["session"].(api.InterfaceSession)
		if !ok {
			env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "55b4054a-fe1a-4b5a-a05e-65fd6c0c2103", "session can't be used"))
			return false
		}

		usedDiscounts := utils.InterfaceToStringArray(session.Get(ConstSessionKeyUsedDiscountCodes))

		for _, discount := range orderAppliedDiscounts {
			usedDiscounts = append(usedDiscounts, discount.Code)
		}

		session.Set(ConstSessionKeyUsedDiscountCodes, usedDiscounts)
	}

	return true
}
