package ini

import (
	"fmt"
	"strings"
)

// ListItems returns all ini file values, for current and global sections
func (it *DefaultIniConfig) ListItems() []string {

	flatMap := make(map[string]bool)

	// collection values from global and current section
	for _, sectionName := range []string{it.currentSection, ConstIniGlobalSection} {
		if sectionValues, present := it.iniFileValues[sectionName]; present {
			for itemName := range sectionValues {
				flatMap[itemName] = true
			}
		}
	}

	// making array from collected items
	result := make([]string, 0, len(flatMap))
	for itemName := range flatMap {
		result = append(result, itemName)
	}

	return result
}

// GetValue returns specified value from ini file, looks for value in current section then in global
func (it *DefaultIniConfig) GetValue(valueName string, defaultValue string) string {

	// looking for value in current section and global section
	for _, sectionName := range []string{it.currentSection, ConstIniGlobalSection} {
		if sectionValues, present := it.iniFileValues[sectionName]; present {
			if value, present := sectionValues[valueName]; present {
				return value
			}
		} else {
			it.iniFileValues[sectionName] = make(map[string]string)
		}
	}

	// value was not found - using default
	if strings.HasPrefix(defaultValue, ConstAskForValuePrefix) {
		value := strings.TrimPrefix(defaultValue, ConstAskForValuePrefix)

		fmt.Printf("%s: ", valueName)
		fmt.Scanf("%s", &value)

		if _, present := it.iniFileValues[it.currentSection]; !present {
			it.iniFileValues[it.currentSection] = make(map[string]string)
		}

		it.iniFileValues[it.currentSection][valueName] = value
		it.keysToStore[valueName] = true

		return value
	}
	it.iniFileValues[ConstIniGlobalSection][valueName] = defaultValue

	return defaultValue
}

// SetWorkingSection changes working ini section to specified
func (it *DefaultIniConfig) SetWorkingSection(sectionName string) error {
	it.currentSection = sectionName
	return nil
}

// GetSectionValue returns value assigned to specified ini section only, or [defaultValue] if not assigned
func (it *DefaultIniConfig) GetSectionValue(sectionName string, valueName string, defaultValue string) string {
	if value, present := it.iniFileValues[sectionName][valueName]; present {
		return value
	}
	return defaultValue
}

// ListSections enumerates currently used ini sections
func (it *DefaultIniConfig) ListSections() []string {
	result := make([]string, 0, len(it.iniFileValues))
	for sectionName := range it.iniFileValues {
		result = append(result, sectionName)
	}
	return result
}

// ListSectionItems enumerates value names within specified ini section
func (it *DefaultIniConfig) ListSectionItems(sectionName string) []string {
	var result []string
	if sectionValues, present := it.iniFileValues[sectionName]; present {
		for valueName := range sectionValues {
			result = append(result, valueName)
		}
	}
	return result
}

// SetValue sets/updates current working section value, modified value marks for saving
func (it *DefaultIniConfig) SetValue(valueName string, value string) error {
	if _, present := it.iniFileValues[it.currentSection]; !present {
		it.iniFileValues[it.currentSection] = make(map[string]string)
	}

	it.iniFileValues[it.currentSection][valueName] = value
	it.keysToStore[valueName] = true

	return nil
}
