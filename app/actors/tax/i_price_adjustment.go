package tax

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/utils"
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

	return []float64{ConstPriorityValue}
}

// processRecords processes records from database collection
func processRecords(name string, records []map[string]interface{}, result []checkout.StructPriceAdjustment) []checkout.StructPriceAdjustment {
	for _, record := range records {
		taxRate := checkout.StructPriceAdjustment{
			Code:      utils.InterfaceToString(record["code"]),
			Label:     name,
			Amount:    utils.InterfaceToFloat64(record["rate"]),
			IsPercent: true,
			Priority:  priority,
			Types:     []string{checkout.ConstLabelTax},
			PerItem:   make(map[int]float64),
		}

		priority += float64(0.00001)
		result = append(result, taxRate)
	}

	return result
}

// Calculate calculates a taxes for a given checkout
func (it *DefaultTax) Calculate(currentCheckout checkout.InterfaceCheckout) []checkout.StructPriceAdjustment {
	var result []checkout.StructPriceAdjustment
	priority = ConstPriorityValue

	if shippingAddress := currentCheckout.GetShippingAddress(); shippingAddress != nil {
		state := shippingAddress.GetState()
		zip := shippingAddress.GetZipCode()

		if dbEngine := db.GetDBEngine(); dbEngine != nil {
			if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
				collection.AddFilter("state", "=", "*")
				collection.AddFilter("zip", "=", "*")

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, result)
				}

				collection.ClearFilters()
				collection.AddFilter("state", "=", state)
				collection.AddFilter("zip", "=", "*")

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, result)
				}

				collection.ClearFilters()
				collection.AddFilter("state", "=", state)
				collection.AddFilter("zip", "=", zip)

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, result)
				}
			}
		}
	}

	return result
}
