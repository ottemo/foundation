package order

import "github.com/ottemo/foundation/utils"

// GetID returns order item unique id, or blank string
func (it *DefaultOrderItem) GetID() string {
	return it.id
}

// GetProductID returns product ID which order item represents
func (it *DefaultOrderItem) GetProductID() string {
	return it.ProductID
}

// SetID sets order item unique id
func (it *DefaultOrderItem) SetID(newID string) error {
	it.id = newID
	return nil
}

// GetName returns order item product name
func (it *DefaultOrderItem) GetName() string {
	return it.Name
}

// GetSku returns order item product sku
func (it *DefaultOrderItem) GetSku() string {
	return it.Sku
}

// GetQty returns order line item qty ordered
func (it *DefaultOrderItem) GetQty() int {
	return it.Qty
}

// GetPrice returns order item product price
func (it *DefaultOrderItem) GetPrice() float64 {
	return it.Price
}

// GetWeight returns order item product weight
func (it *DefaultOrderItem) GetWeight() float64 {
	return it.Weight
}

// GetOptions returns order item product options
func (it *DefaultOrderItem) GetOptions() map[string]interface{} {
	return it.Options
}

// GetSelectedOptions returns order item options as a simple map
// optionId: optionValue or optionLabel: optionValueLabel
func (it *DefaultOrderItem) GetSelectedOptions(asLabels bool) map[string]interface{} {
	result := make(map[string]interface{})

	// order items extraction
	if asLabels {
		for _, value := range it.GetOptions() {
			option := utils.InterfaceToMap(value)
			optionValue := option["value"]
			optionLabel := utils.InterfaceToString(option["label"])
			result[optionLabel] = optionValue

			if options, present := option["options"]; present {
				optionsMap := utils.InterfaceToMap(options)
				if optionValueParameters, ok := optionsMap[utils.InterfaceToString(optionValue)]; ok {
					optionValueParametersMap := utils.InterfaceToMap(optionValueParameters)
					result[optionLabel] = optionValueParametersMap["label"]
				}
			}
		}
		return result
	}

	for key, value := range it.GetOptions() {
		option := utils.InterfaceToMap(value)
		result[key] = option["value"]
	}

	return result
}
