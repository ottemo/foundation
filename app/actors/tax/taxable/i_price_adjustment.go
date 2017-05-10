package taxable

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/utils"
)



// GetName returns name of current tax implementation
func (it *DefaultTaxable) GetName() string {
	return "Taxable"
}

// GetCode returns code of current tax implementation
func (it *DefaultTaxable) GetCode() string {
	return "taxable"
}

// GetPriority returns the code of the current coupon implementation
func (it *DefaultTaxable) GetPriority() []float64 {
	return []float64{checkout.ConstCalculateTargetSubtotal, checkout.ConstCalculateTargetGrandTotal}
}

// Calculate calculates a taxes for a given checkout
func (it *DefaultTaxable) Calculate(currentCheckout checkout.InterfaceCheckout, currentPriority float64) []checkout.StructPriceAdjustment {
	var result []checkout.StructPriceAdjustment

	if currentPriority == checkout.ConstCalculateTargetSubtotal || currentPriority == checkout.ConstCalculateTargetGrandTotal {
		var cartItems = currentCheckout.GetItems()
		perItem := make(map[string]float64)

		for _, cartItem := range cartItems {
			var productItem = cartItem.GetProduct()
			if productItem == nil {
				return result
			}

			var productMap = productItem.ToHashMap()
			if _, present := productMap[ConstProductTaxableAttribute]; !present {
				continue
			}

			var attributesInfo = productItem.GetAttributesInfo()
			if attributesInfo == nil || len(attributesInfo) == 0 {
				return result
			}

			for _, attributeInfo := range attributesInfo {
				if attributeInfo.Attribute == ConstProductTaxableAttribute &&
					attributeInfo.Type == utils.ConstDataTypeBoolean &&
					!utils.InterfaceToBool(productItem.Get(ConstProductTaxableAttribute)) {

					var itemIndex = utils.InterfaceToString(cartItem.GetIdx())

					if currentPriority == checkout.ConstCalculateTargetSubtotal {
						// discount "non taxable" on 100%, so they wouldn't be discounted or taxed
						perItem[itemIndex] = -100 // -100%
					} else if currentPriority == checkout.ConstCalculateTargetGrandTotal {
						// restore "non taxable" amounts as they basic subtotal
						perItem[itemIndex] = currentCheckout.GetItemSpecificTotal(itemIndex, checkout.ConstLabelSubtotal)
					}
				}
			}
		}

		if perItem == nil || len(perItem) == 0 {
			return result
		}

		if currentPriority == checkout.ConstCalculateTargetSubtotal {
			// discount "non taxable" on 100%, so they wouldn't be discounted or taxed
			result = append(result, checkout.StructPriceAdjustment{
				Code:      "TN", // Taxable? - No
				Name:      it.GetName(),
				Amount:    -100,
				IsPercent: true,
				Priority:  checkout.ConstCalculateTargetSubtotal,
				Labels:    []string{checkout.ConstLabelTax},
				PerItem:   perItem,
			})
		} else if currentPriority == checkout.ConstCalculateTargetGrandTotal {
			// restore "non taxable" amounts as they basic subtotal
			result = append(result, checkout.StructPriceAdjustment{
				Code:      "TN", // Taxable? - No
				Name:      it.GetName(),
				Amount:    0,
				IsPercent: false,
				Priority:  checkout.ConstCalculateTargetGrandTotal,
				Labels:    []string{checkout.ConstLabelTax},
				PerItem:   perItem,
			})
		}

		return result
	}

	return result
}

