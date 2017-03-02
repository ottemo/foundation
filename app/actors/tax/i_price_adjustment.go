package tax

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/env"

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
				if err := collection.AddFilter("state", "=", "*"); err != nil {
					_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "6ff4f7b9-204d-4cef-933e-1b50c0a7810f", err.Error())
				}
				if err := collection.AddFilter("zip", "=", "*"); err != nil {
					_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "56e7bcc3-da2f-4b14-a9d5-207144bab513", err.Error())
				}

				if records, err := collection.Load(); err == nil {
					result = it.processRecords(records, result)
				}

				if err := collection.ClearFilters(); err != nil {
					_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "45f36a5a-7039-4383-af0a-5d7014fd2972", err.Error())
				}
				if err := collection.AddFilter("state", "=", state); err != nil {
					_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "48cf0d87-3335-4276-9a9c-97d400fe3229", err.Error())
				}
				if err := collection.AddFilter("zip", "=", "*"); err != nil {
					_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "0dcafa80-5e1a-4848-8761-4aad8b9e7fff", err.Error())
				}

				if records, err := collection.Load(); err == nil {
					result = it.processRecords(records, result)
				}

				if err := collection.ClearFilters(); err != nil {
					_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "7c2fb5ab-15f9-4d61-a0e7-a4611796d783", err.Error())
				}
				if err := collection.AddFilter("state", "=", state); err != nil {
					_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "500125dd-5e8f-4130-9877-b5df489373dd", err.Error())
				}
				if err := collection.AddFilter("zip", "=", zip); err != nil {
					_ = env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "e5f12dd7-70ca-4563-aa51-1bfd913ac797", err.Error())
				}

				if records, err := collection.Load(); err == nil {
					result = it.processRecords(records, result)
				}
			}
		}
	}

	return result
}
