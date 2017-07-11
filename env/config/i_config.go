package config

import (
	"sort"
	"strings"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// ListPathes enumerates registered pathes for config
func (it *DefaultConfig) ListPathes() []string {
	var result []string
	for key := range it.configValues {
		result = append(result, key)
	}
	sort.Strings(result)

	return result
}

// RegisterItem registers new config value in system
func (it *DefaultConfig) RegisterItem(Item env.StructConfigItem, Validator env.FuncConfigValueValidator) error {

	// registering new config item
	if _, present := it.configValues[Item.Path]; !present {

		collection, err := db.GetCollection(ConstCollectionNameConfig)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		if Item.Type == env.ConstConfigTypeSecret {
			stringValue := utils.InterfaceToString(Item.Value)

			// it should not modify value if string was not encrypted
			stringValue = utils.DecryptString(stringValue)
			Item.Value = utils.EncryptString(stringValue)
		}

		recordValues := make(map[string]interface{})

		recordValues["path"] = Item.Path
		recordValues["value"] = Item.Value
		recordValues["type"] = Item.Type
		recordValues["editor"] = Item.Editor
		recordValues["options"] = Item.Options
		recordValues["label"] = Item.Label
		recordValues["description"] = Item.Description

		_, err = collection.Save(recordValues)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		it.configValues[Item.Path] = Item.Value
		it.configTypes[Item.Path] = Item.Type
	}

	// registering validator
	if _, present := it.configValidators[Item.Path]; Validator != nil && !present {
		it.configValidators[Item.Path] = Validator
	}

	// validating current set value
	if validator, present := it.configValidators[Item.Path]; present && validator != nil {
		newValue, err := validator(it.configValues[Item.Path])
		if err != nil {
			if err := it.SetValue(Item.Path, newValue); err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}

	return nil
}

// UnregisterItem removes config value from system
func (it *DefaultConfig) UnregisterItem(Path string) error {

	if _, present := it.configValues[Path]; present {

		collection, err := db.GetCollection(ConstCollectionNameConfig)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = collection.AddFilter("path", "LIKE", Path+"%")
		if err != nil {
			return env.ErrorDispatch(err)
		}

		_, err = collection.Delete()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		return it.Reload()
	}

	return nil
}

// GetValue returns value for config item of nil if not present
func (it *DefaultConfig) GetValue(Path string) interface{} {
	if value, present := it.configValues[Path]; present {

		if it.configTypes[Path] == env.ConstConfigTypeSecret {
			stringValue := utils.InterfaceToString(value)
			return utils.DecryptString(stringValue)
		}

		return value
	}
	return nil
}

// SetValue updates config item with new value, returns error if not possible
func (it *DefaultConfig) SetValue(Path string, Value interface{}) error {
	if _, present := it.configValues[Path]; present {

		// updating value on GO side
		//--------------------------
		if validator, present := it.configValidators[Path]; present {

			newVal, err := validator(Value)
			if err != nil {
				return env.ErrorDispatch(err)
			}

			it.configValues[Path] = newVal

		} else {
			it.configValues[Path] = Value
		}

		if it.configTypes[Path] == env.ConstConfigTypeSecret {
			stringValue := utils.InterfaceToString(it.configValues[Path])

			// it should not modify value if string was not encrypted
			stringValue = utils.DecryptString(stringValue)
			it.configValues[Path] = utils.EncryptString(stringValue)
		}

		// updating value in DB
		//---------------------
		collection, err := db.GetCollection(ConstCollectionNameConfig)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		err = collection.AddFilter("path", "=", Path)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		records, err := collection.Load()
		if err != nil {
			return env.ErrorDispatch(err)
		}

		if len(records) == 0 {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1f29fe4c-6062-48b5-a106-e3ac94339cda", "config item '"+Path+"' is not registered")
		}

		record := records[0]

		record["value"] = it.configValues[Path]

		_, err = collection.Save(record)
		if err != nil {
			return env.ErrorDispatch(err)
		}

	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6984f1ce-1fb1-40d5-b674-9d88956164c0", "can not find config item '"+Path+"' ")
	}

	return nil
}

// GetGroupItems returns information about config items with type [ConstConfigTypeGroup]
func (it *DefaultConfig) GetGroupItems() []env.StructConfigItem {

	var result []env.StructConfigItem

	collection, err := db.GetCollection(ConstCollectionNameConfig)
	if err != nil {
		return result
	}

	err = collection.AddFilter("type", "=", env.ConstConfigTypeGroup)
	if err != nil {
		return result
	}

	records, err := collection.Load()
	if err != nil {
		return result
	}

	for _, record := range records {

		valueType := utils.InterfaceToString(record["type"])
		valuePath := utils.InterfaceToString(record["path"])
		configItem := env.StructConfigItem{
			Path: valuePath,
			Type: valueType,

			Editor:  utils.InterfaceToString(record["editor"]),
			Options: record["options"],

			Label:       utils.InterfaceToString(record["label"]),
			Description: utils.InterfaceToString(record["description"]),

			Image: utils.InterfaceToString(record["image"]),
		}

		if valueType == env.ConstConfigTypeSecret {
			configItem.Value = it.GetValue(valuePath)
		} else {
			configItem.Value = db.ConvertTypeFromDbToGo(record["value"], valueType)
		}

		result = append(result, configItem)
	}

	return result
}

// GetItemsInfo returns information about config items with given path
// 	- use '*' to list sub-items (like "paypal.*" or "paypal*" if group item also needed)
func (it *DefaultConfig) GetItemsInfo(Path string) []env.StructConfigItem {
	var result []env.StructConfigItem

	collection, err := db.GetCollection(ConstCollectionNameConfig)
	if err != nil {
		return result
	}

	err = collection.AddFilter("path", "LIKE", strings.Replace(Path, "*", "%", -1))
	if err != nil {
		return result
	}

	records, err := collection.Load()
	if err != nil {
		return result
	}

	for _, record := range records {

		valueType := utils.InterfaceToString(record["type"])
		valuePath := utils.InterfaceToString(record["path"])
		configItem := env.StructConfigItem{
			Path: valuePath,
			Type: valueType,

			Editor:  utils.InterfaceToString(record["editor"]),
			Options: record["options"],

			Label:       utils.InterfaceToString(record["label"]),
			Description: utils.InterfaceToString(record["description"]),

			Image: utils.InterfaceToString(record["image"]),
		}

		if valueType == env.ConstConfigTypeSecret {
			configItem.Value = it.GetValue(valuePath)
		} else {
			configItem.Value = db.ConvertTypeFromDbToGo(record["value"], valueType)
		}

		result = append(result, configItem)
	}

	return result
}

// Load loads config data from DB on app startup
//   - calls env.OnConfigStart() after
func (it *DefaultConfig) Load() error {

	err := it.Reload()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = env.OnConfigStart()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Reload updates all config values from database
func (it *DefaultConfig) Reload() error {
	it.configValues = make(map[string]interface{})
	it.configTypes = make(map[string]string)

	collection, err := db.GetCollection(ConstCollectionNameConfig)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.SetResultColumns("path", "type", "value")
	if err != nil {
		return env.ErrorDispatch(err)
	}

	records, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	for _, record := range records {
		valuePath := utils.InterfaceToString(record["path"])
		valueType := utils.InterfaceToString(record["type"])

		if valueType == env.ConstConfigTypeSecret {
			it.configValues[valuePath] = utils.InterfaceToString(record["value"])
		} else {
			it.configValues[valuePath] = db.ConvertTypeFromDbToGo(record["value"], valueType)
		}
		it.configTypes[valuePath] = valueType
	}

	return nil
}
