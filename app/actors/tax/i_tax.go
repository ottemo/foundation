package tax

import (
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
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
func processRecords(name string, records []map[string]interface{}, taxableAmount float64, result []checkout.StructTaxRate) []checkout.StructTaxRate {
	priorityValue := ConstPriorityValue
	for _, record := range records {
		taxRate := checkout.StructTaxRate{
			Name:      name,
			Code:      utils.InterfaceToString(record["code"]),
			Amount:    taxableAmount * utils.InterfaceToFloat64(record["rate"]) / 100,
			IsPercent: false,
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

		taxableAmount := currentCheckout.GetSubtotal() + currentCheckout.GetShippingAmount()

		// event which allows to change and/or track taxable cart amount before tax calculation
		eventData := map[string]interface{}{"tax": it, "checkout": currentCheckout, "amount": taxableAmount}
		env.Event("tax.amount", eventData)
		if newAmount := utils.InterfaceToFloat64(eventData["amount"]); newAmount >= 0 && taxableAmount != newAmount {
			taxableAmount = newAmount
		}

		if dbEngine := db.GetDBEngine(); dbEngine != nil {
			if collection, err := dbEngine.GetCollection("Taxes"); err == nil {
				collection.AddFilter("state", "=", "*")
				collection.AddFilter("zip", "=", "*")

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, taxableAmount, result)
				}

				collection.ClearFilters()
				collection.AddFilter("state", "=", state)
				collection.AddFilter("zip", "=", "*")

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, taxableAmount, result)
				}

				collection.ClearFilters()
				collection.AddFilter("state", "=", state)
				collection.AddFilter("zip", "=", zip)

				if records, err := collection.Load(); err == nil {
					result = processRecords(it.GetName(), records, taxableAmount, result)
				}
			}
		}
	}

	return result
}
