package config

import (
	"github.com/ottemo/foundation/utils"
)

// Import imports value in config through Impex
func (it *DefaultConfig) Import(item map[string]interface{}, testMode bool) (map[string]interface{}, error) {
	var path string
	var value interface{}

	if pathValue, present := item["path"]; present {
		path = utils.InterfaceToString(pathValue)
	} else if keyValue, present := item["key"]; present {
		path = utils.InterfaceToString(keyValue)
	}

	value, _ = item["value"]

	if testMode == false {
		err := it.SetValue(path, value)
		return item, err
	}

	return item, nil
}

// Export exports config values through Impex
func (it *DefaultConfig) Export(iterator func(map[string]interface{}) bool) error {
	for itemPath, itemValue := range it.configValues {
		continueFlag := iterator(map[string]interface{}{"path": itemPath, "value": itemValue})
		if continueFlag == false {
			break
		}
	}

	return nil
}
