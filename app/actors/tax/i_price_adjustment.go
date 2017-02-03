package tax

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/checkout"
)

// GetName returns name of current tax implementation
func (it *DefaultTax) GetName() string {
	return "Tax"
}

// GetCode returns code of current tax implementation
func (it *DefaultTax) GetCode() string {
	return "tax"
}

// GetPriority returns the code of the current coupon implementation
func (it *DefaultTax) GetPriority() []float64 {

	return []float64{checkout.ConstCalculateTargetSubtotal, ConstPriorityValue, checkout.ConstCalculateTargetGrandTotal}
}

// processRecords processes records from database collection
func (it *DefaultTax) processRecords(records []map[string]interface{}, result []checkout.StructPriceAdjustment) []checkout.StructPriceAdjustment {
	for _, record := range records {
		amount := utils.InterfaceToFloat64(record["rate"])

		taxRate := checkout.StructPriceAdjustment{
			Code:      utils.InterfaceToString(record["code"]),
			Name:      it.GetName(),
			Amount:    amount,
			IsPercent: true,
			Priority:  priority,
			Labels:    []string{checkout.ConstLabelTax},
		}

		priority += float64(0.00001)
		result = append(result, taxRate)
	}

	return result
}

// Calculate calculates a taxes for a given checkout
func (it *DefaultTax) Calculate(currentCheckout checkout.InterfaceCheckout, currentPriority float64) []checkout.StructPriceAdjustment {
	var result []checkout.StructPriceAdjustment
	priority = ConstPriorityValue

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

	if shippingAddress := currentCheckout.GetShippingAddress(); shippingAddress != nil {
		state := shippingAddress.GetState()
		zip := shippingAddress.GetZipCode()

		if dbEngine := db.GetDBEngine(); dbEngine != nil {
			if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
				collection.AddFilter("state", "=", "*")
				collection.AddFilter("zip", "=", "*")

				if records, err := collection.Load(); err == nil {
					result = it.processRecords(records, result)
				}

				collection.ClearFilters()
				collection.AddFilter("state", "=", state)
				collection.AddFilter("zip", "=", "*")

				if records, err := collection.Load(); err == nil {
					result = it.processRecords(records, result)
				}

				collection.ClearFilters()
				collection.AddFilter("state", "=", state)
				collection.AddFilter("zip", "=", zip)

				if records, err := collection.Load(); err == nil {
					result = it.processRecords(records, result)
				}
			}
		}
	}

	return result
}
