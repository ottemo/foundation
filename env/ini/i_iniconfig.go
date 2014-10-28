package ini

import (
	"fmt"
)

func (it *DefaultIniConfig) ListItems() []string {
	result := make([]string, len(it.iniFileValues))
	for itemName, _ := range it.iniFileValues {
		result = append(result, itemName)
	}
	return result
}

func (it *DefaultIniConfig) GetValue(valueName string, defaultValue string) string {

	if value, present := it.iniFileValues[valueName]; present {
		return value
	} else {

		if defaultValue == "?" {
			value = ""

			fmt.Printf("%s: ", valueName)
			fmt.Scanf("%s", &value)

			it.iniFileValues[valueName] = value
			it.keysToStore = append(it.keysToStore, valueName)

			return value
		} else {
			return defaultValue
		}
	}
}
