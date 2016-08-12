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

// GetOptionValues returns order item options values
// optionId: optionValue or optionLabel: optionValueLabel
func (it *DefaultOrderItem) GetOptionValues(labels bool) map[string]interface{} {
	result := make(map[string]interface{})

	// order items extraction
	if labels {
		// this part is hard version of moving through key's of option just to get one value
		for key, value := range it.GetOptions() {
			option := utils.InterfaceToMap(value)
			optionLabel := key
			if labelValue, optionLabelPresent := option["label"]; optionLabelPresent {
				optionLabel = utils.InterfaceToString(labelValue)
			}

			result[optionLabel] = value
			optionValue, optionValuePresent := option["value"]
			// in this case looks like structure of options was changed or it's not a map
			if !optionValuePresent {
				continue
			}
			result[optionLabel] = optionValue

			optionType := ""

			if val, present := option["type"]; present {
				optionType = utils.InterfaceToString(val)
			}
			if options, present := option["options"]; present {
				optionsMap := utils.InterfaceToMap(options)

				if optionType == "multi_select" {
					selectedOptions := ""
					for i, optionValue := range utils.InterfaceToArray(optionValue) {
						if optionValueParameters, ok := optionsMap[utils.InterfaceToString(optionValue)]; ok {
							optionValueParametersMap := utils.InterfaceToMap(optionValueParameters)
							if labelValue, labelValuePresent := optionValueParametersMap["label"]; labelValuePresent {
								result[optionLabel] = labelValue
								if i > 0 {
									selectedOptions = selectedOptions + ", "
								}
								selectedOptions = selectedOptions + utils.InterfaceToString(labelValue)
							}
						}
					}
					result[optionLabel] = selectedOptions

				} else if optionValueParameters, ok := optionsMap[utils.InterfaceToString(optionValue)]; ok {
					optionValueParametersMap := utils.InterfaceToMap(optionValueParameters)
					if labelValue, labelValuePresent := optionValueParametersMap["label"]; labelValuePresent {
						result[optionLabel] = labelValue
					}

				}
			}
		}
		return result
	}

	for key, value := range it.GetOptions() {
		option := utils.InterfaceToMap(value)
		if optionValue, present := option["value"]; present {
			result[key] = optionValue
		} else {
			result[key] = value
		}
	}

	return result
}
