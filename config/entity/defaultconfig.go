package defaultconfig

import (
	"errors"
	"fmt"

	config "github.com/ottemo/foundation/config"
	db "github.com/ottemo/foundation/database"
)

// DefaultConfigEntity is a struct to hold a configuration value and its validator.
// New configuration entities may be created and added at runtime.
type DefaultConfigEntity struct {
	Name      string
	Validator func(interface{}) (interface{}, bool)
	Default   interface{}
	Value     interface{}
}

// DefaultConfig is a map container of configuration entities.
type DefaultConfig struct {
	configValues map[string]*DefaultConfigEntity
}

func init() {
	instance := new(DefaultConfig)

	db.RegisterOnDatabaseStart(instance.Load)
	config.RegisterConfig(new(DefaultConfig))
}

// Register returns nil or an error if an entity has already been registered.
func (dc *DefaultConfig) Register(Name string, Validator func(interface{}) (interface{}, bool), Default interface{}) error {
	if _, present := dc.configValues[Name]; present {
		return errors.New("[ " + Name + " ] is already registered")
	}
	dc.configValues[Name] = &DefaultConfigEntity{Name: Name, Validator: Validator, Default: Default, Value: Default}

	return nil
}

// Destroy returns nil if it is able to delete given entity, otherwise it returns an error.
func (dc *DefaultConfig) Destroy(Name string) error {
	if _, present := dc.configValues[Name]; present {
		delete(dc.configValues, Name)
	} else {
		return errors.New("[ " + Name + " ] does not exist")
	}

	return nil
}

// GetValue returns the value associated with the given string.
func (dc *DefaultConfig) GetValue(Name string) interface{} {
	if configItem, present := dc.configValues[Name]; present {
		return configItem.Value
	}

	return nil
}

//SetValue will set an entity to a given value or return an error.
func (dc *DefaultConfig) SetValue(Name string, Value interface{}) error {
	if configItem, present := dc.configValues[Name]; present {
		if configItem.Validator != nil {
			if newVal, ok := configItem.Validator(Value); ok {
				return errors.New("Attempting to set invalid value to [ " + Name + " ] = " + fmt.Sprintf("%s", Value))
			}
			configItem.Value = newVal
		} else {
			configItem.Value = Value
		}
		return nil
	}
	return errors.New("Cannot find a value for [ " + Name + " ] ")
}

// List returns a list of entities.
func (dc *DefaultConfig) List() []string {
	result := make([]string, len(me.configValues))
	for itemName := range dc.configValues {
		result = append(result, itemName)
	}
	return result
}

// Save will persist the configuration.
func (dc *DefaultConfig) Save() error {
	return nil
}

// Load will executes the registered callbacks for this entity upon start.
func (dc *DefaultConfig) Load() error {

	config.OnConfigStart()

	return nil
}
