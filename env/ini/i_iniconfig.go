package ini

import (
	"fmt"
	"strings"
)

// returns all ini file values, for current and global sections
func (it *DefaultIniConfig) ListItems() []string {

	flatMap := make(map[string]bool)

	// collection values from global and current section
	for _, sectionName := range []string{it.currentSection, ""} {
		if sectionValues, present := it.iniFileValues[sectionName]; present {
			for itemName, _ := range sectionValues {
				flatMap[itemName] = true
			}
		}
	}

	// making array from collected items
	result := make([]string, 0, len(flatMap))
	for itemName, _ := range flatMap {
		result = append(result, itemName)
	}

	return result
}

// returns specified value from ini file
func (it *DefaultIniConfig) GetValue(valueName string, defaultValue string) string {

	// looking for value in current section and global section
	for _, sectionName := range []string{it.currentSection, ""} {
		if sectionValues, present := it.iniFileValues[sectionName]; present {
			if value, present := sectionValues[valueName]; present {
				return value
			}
		} else {
			it.iniFileValues[sectionName] = make(map[string]string)
		}
	}

	// value was not found - using default
	if strings.HasPrefix(defaultValue, ASK_FOR_VALUE_PREFIX) {
		value := strings.TrimPrefix(defaultValue, ASK_FOR_VALUE_PREFIX)

		fmt.Printf("%s: ", valueName)
		fmt.Scanf("%s", &value)

		if _, present := it.iniFileValues[it.currentSection]; !present {
			it.iniFileValues[it.currentSection] = make(map[string]string)
		}

		it.iniFileValues[it.currentSection][valueName] = value
		it.keysToStore = append(it.keysToStore, valueName)

		return value
	} else {
		return defaultValue
	}
}
