package product

import (
	"github.com/ottemo/foundation/utils"
	"github.com/ottemo/foundation/app/models/product"
)

func UpdateProductOptions(product product.InterfaceProduct) map[string]interface{}{

	newOptions := make(map[string]interface{})

	for optionsName, currentOption := range product.GetOptions() {
		currentOption := utils.InterfaceToMap(currentOption)

		if option, present := currentOption["options"]; present {
			newOptionValues := make(map[string]interface{})

			for key, value := range utils.InterfaceToMap(option) {
				newOptionValues[utils.StrToSnakeCase(key)] = value
			}

			currentOption["options"] = newOptionValues
		}
		newOptions[utils.StrToSnakeCase(optionsName)] = currentOption
	}

	return newOptions
}

