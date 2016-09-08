package saleprice

import (
	"time"

	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/discount/saleprice"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

//----------------------------------------------------------------------------------------------------------------------
// Implementation of github.com/ottemo/foundation/app/models/checkout/interfaces InterfacePriceAdjustment
//----------------------------------------------------------------------------------------------------------------------

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
	return []float64{utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathSalePriceApplyPriority))}
}

// Calculate calculates and returns amount and set of applied discounts to given checkout
func (it *DefaultSalePrice) Calculate(checkoutInstance checkout.InterfaceCheckout, currentPriority float64) []checkout.StructPriceAdjustment {
	var result []checkout.StructPriceAdjustment

	allowedApplyPriority := utils.InterfaceToFloat64(env.ConfigGetValue(ConstConfigPathSalePriceApplyPriority))
	cardGrandTotal := checkoutInstance.GetItemSpecificTotal(0, checkout.ConstLabelGrandTotal)

	if currentPriority == allowedApplyPriority && cardGrandTotal > 0 {

		salePriceCollection, err := db.GetCollection(saleprice.ConstSalePriceDbCollectionName)
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

			itemGrandTotal := checkoutInstance.GetItemSpecificTotal(item.GetIdx(), checkout.ConstLabelGrandTotal)

			for _, salePrice := range salePrices {
				suggestedSalePrice := utils.InterfaceToFloat64(item.GetQty()) * utils.InterfaceToFloat64(salePrice["amount"])

				// do not use sale price if it greater than current item calculated total
				if itemGrandTotal < suggestedSalePrice {
					continue
				}

				if utils.InterfaceToTime(salePrice["start_datetime"]).Before(today) &&
					utils.InterfaceToTime(salePrice["end_datetime"]).After(today) {
					perItem[utils.InterfaceToString(item.GetIdx())] = -suggestedSalePrice

					// Because of time ranges are not overlapped, first found sale price is
					// acceptable
					break
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
			Priority:  allowedApplyPriority,
			Labels:    []string{checkout.ConstLabelSalePriceAdjustment},
			PerItem:   perItem,
		})

		return result
	}

	return nil
}
