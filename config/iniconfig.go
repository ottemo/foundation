package config

import (
	"errors"
)

var registeredIniConfig IniConfig
var callbacksOnConfigIniStart = []func() error{}

// IniConfig is an initialization interface for reading INI file values
type IniConfig interface {
	GetValue(Name string) string
	List() []string
}

// RegisterOnConfigIniStart will register the ini file upon application start.
func RegisterOnConfigIniStart(callback func() error) {
	callbacksOnConfigIniStart = append(callbacksOnConfigIniStart, callback)
}

// OnConfigIniStart executes the registered callbacks to be run when the INI file has been initialized.
func OnConfigIniStart() error {
	for _, callback := range callbacksOnConfigIniStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

// RegisterIniConfig registers the new INI configuration file.
func RegisterIniConfig(IniConfig IniConfig) error {
	if registeredIniConfig == nil {
		registeredIniConfig = IniConfig
	} else {
		return errors.New("Configuration file already registered. Unable to register configuration.")
	}
	return nil
}

// GetIniConfig returns the registered INI configuration that has been initialized.
func GetIniConfig() IniConfig { return registeredIniConfig }

// ConfigEmptyValueValidator initializes an empty value validator.
func ConfigEmptyValueValidator(val interface{}) (interface{}, bool) { return val, true }
