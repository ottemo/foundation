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

// processRecords processes records from database collection
func processRecords(name string, records []map[string]interface{}, cartGrandTotal float64, result []checkout.StructTaxRate) []checkout.StructTaxRate {
	priorityValue := ConstPriorityValue
	for _, record := range records {
		taxRate := checkout.StructTaxRate{
			Name:      name,
			Code:      utils.InterfaceToString(record["code"]),
			Amount:    utils.InterfaceToFloat64(record["rate"]),
			IsPercent: true,
			Priority:  priorityValue,
		}

		priorityValue += float64(0.0001)
		result = append(result, taxRate)
	}

	return result
}

// CalculateTax calculates a taxes for a given checkout
func (it *DefaultTax) CalculateTax(currentCheckout checkout.InterfaceCheckout) []checkout.StructTaxRate {
	var result []checkout.StructTaxRate

	if shippingAddress := currentCheckout.GetShippingAddress(); shippingAddress != nil {
		state := shippingAddress.GetState()
		zip := shippingAddress.GetZipCode()

		cartGrandTotal := currentCheckout.GetSubtotal() + currentCheckout.GetShippingAmount()

		if dbEngine := db.GetDBEngine(); dbEngine != nil {
			if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
				collection.AddFilter("state", "=", "*")
				collection.AddFilter("zip", "=", "*")

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, cartGrandTotal, result)
				}

				collection.ClearFilters()
				collection.AddFilter("state", "=", state)
				collection.AddFilter("zip", "=", "*")

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, cartGrandTotal, result)
				}

				collection.ClearFilters()
				collection.AddFilter("state", "=", state)
				collection.AddFilter("zip", "=", zip)

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, cartGrandTotal, result)
				}
			}
		}
	}

	return result
}
