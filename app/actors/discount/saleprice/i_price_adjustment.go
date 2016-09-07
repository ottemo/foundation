// Implementation of github.com/ottemo/foundation/app/models/checkout/interfaces InterfacePriceAdjustment
package saleprice

import (
	"time"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetName returns name of current sale price implementation
func (it *DefaultSalePrice) GetName() string {
	return "SalePriceDiscount"
}

// GetCode returns code of current sale price implementation
func (it *DefaultSalePrice) GetCode() string {
	return "saleprice_discount"
}

// GetPriority returns the priority of sale price adjustment during checkout calculation
func (it *DefaultSalePrice) GetPriority() []float64 {
	return []float64{
		checkout.ConstCalculateTargetSubtotal,
		utils.InterfaceToFloat64(
			env.ConfigGetValue(ConstConfigPathSalePriceApplyPriority)),
		checkout.ConstCalculateTargetGrandTotal}
}

// Calculate calculates and returns amount and set of applied discounts to given checkout
func (it *DefaultSalePrice) Calculate(checkoutInstance checkout.InterfaceCheckout, currentPriority float64) []checkout.StructPriceAdjustment {
	var result []checkout.StructPriceAdjustment

	if currentPriority == checkout.ConstCalculateTargetSubtotal {

		salePriceCollection, err := db.GetCollection(saleprice.ConstModelNameSalePriceCollection)
		if err != nil {
			return result
		}

		today := time.Now()
		items := checkoutInstance.GetItems()
		perItem := make(map[string]float64)
		for _, item := range items {
			productItem := item.GetProduct()
			if productItem == nil {
				return result
			}

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

			for _, salePrice := range salePrices {
				if utils.InterfaceToTime(salePrice["start_datetime"]).Before(today) &&
					utils.InterfaceToTime(salePrice["end_datetime"]).After(today) {
					perItem[utils.InterfaceToString(item.GetIdx())] =
						-(utils.InterfaceToFloat64(item.GetQty()) * utils.InterfaceToFloat64(salePrice["amount"]))
				}
			}
		}

		if perItem == nil || len(perItem) == 0 {
			return result
		}
		result = append(result, checkout.StructPriceAdjustment{
			Code:      it.GetCode(),
			Name:      it.GetName(),
			Amount:    0,
			IsPercent: false,
			Priority:  checkout.ConstCalculateTargetSubtotal,
			Labels:    []string{checkout.ConstLabelSalePriceAdjustment},
			PerItem:   perItem,
		})

		return result
	}

	return nil
}
