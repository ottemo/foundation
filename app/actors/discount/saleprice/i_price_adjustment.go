//Implementation of github.com/ottemo/foundation/app/models/checkout/interfaces InterfacePriceAdjustment
package saleprice

import (
	//"github.com/ottemo/foundation/app/models/checkout"
	//"github.com/ottemo/foundation/env"
	//"github.com/ottemo/foundation/utils"
)
import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/db"
	"time"
)


//// GetName returns name of current discount implementation
func (it *DefaultSalePrice) GetName() string {
	return "SalePriceDiscount"
}

//// GetCode returns code of current discount implementation
func (it *DefaultSalePrice) GetCode() string {
	return "saleprice_discount"
}

//// GetPriority returns the code of the current coupon implementation
func (it *DefaultSalePrice) GetPriority() []float64 {
	// TODO: what is PA? adding this first value of priority to make PA that will reduce GT of gift cards by 100% right after subtotal calculation
	return []float64{
		checkout.ConstCalculateTargetSubtotal,
		utils.InterfaceToFloat64(
			env.ConfigGetValue(ConstConfigPathSalePriceApplyPriority)),
		checkout.ConstCalculateTargetGrandTotal}
}

//// Calculate calculates and returns amount and set of applied discounts to given checkout
func (it *DefaultSalePrice) Calculate(checkoutInstance checkout.InterfaceCheckout, currentPriority float64) []checkout.StructPriceAdjustment {
	logDebugHelper("DefaultSalePrice Calculate "+utils.InterfaceToString(currentPriority))
	var result []checkout.StructPriceAdjustment

	logDebugHelper("checkout.ConstCalculateTargetGrandTotal "+utils.InterfaceToString(checkout.ConstCalculateTargetGrandTotal))
	if currentPriority == checkout.ConstCalculateTargetGrandTotal {
		logDebugHelper("checkout.ConstCalculateTargetGrandTotal")

		// loading information about applied coupons
		salePriceCollection, err := db.GetCollection(ConstCollectionNameSalePrices)
		if err != nil {
			return result
		}
		logDebugHelper("got salePriceCollection")

		today := time.Now()
		items := checkoutInstance.GetItems()
		perItem := make(map[string]float64)
		for _, item := range items {
			productItem := item.GetProduct();
			if  productItem == nil {
				return result
			}
			logDebugHelper("process product "+utils.InterfaceToString(productItem))

			err = salePriceCollection.ClearFilters()
			if err != nil {
				return result
			}

			err = salePriceCollection.AddFilter("product_id", "in", productItem.GetID())
			if err != nil {
				return result
			}

			salePrices, err := salePriceCollection.Load()
			if err != nil || len(salePrices) == 0 {
				return result
			}
			logDebugHelper("salePrices "+utils.InterfaceToString(salePrices))

			for _, salePrice := range salePrices {
				logDebugHelper("salePrice "+utils.InterfaceToString(salePrice))
				logDebugHelper("start_datetime "+utils.InterfaceToString(utils.InterfaceToTime(salePrice["start_datetime"])))
				logDebugHelper("end_datetime "+utils.InterfaceToString(utils.InterfaceToTime(salePrice["end_datetime"])))
				logDebugHelper("today "+utils.InterfaceToString(utils.InterfaceToTime(today)))
				if utils.InterfaceToTime(salePrice["start_datetime"]).Before(today) &&
					utils.InterfaceToTime(salePrice["end_datetime"]).After(today) {
					logDebugHelper("today is in time range")
					logDebugHelper("perItem key "+utils.InterfaceToString(item.GetIdx()))
					perItem[utils.InterfaceToString(item.GetIdx())] =
						-utils.InterfaceToFloat64(salePrice["amount"])
				}
			}
		}

		if perItem == nil || len(perItem) == 0 {
			return result
		}
		logDebugHelper("perItem "+utils.InterfaceToString(perItem))
		result = append(result, checkout.StructPriceAdjustment{
			Code:      it.GetCode(),
			Name:      it.GetName(),
			Amount:    0,
			IsPercent: false,
			Priority:  checkout.ConstCalculateTargetSubtotal,
			Labels:    []string{checkout.ConstLabelSalePriceAdjustment},
			PerItem:   perItem,
		})
		logDebugHelper("result "+utils.InterfaceToString(result))

		return result
	}

	return nil
}

