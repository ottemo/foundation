package config

import (
	"errors"
)

var registeredConfig Config
var callbacksOnConfigStart = []func() error{}

// Config is an interface for working with configuration entities and values
type Config interface {
	Register(Name string, Validator func(interface{}) (interface{}, bool), Default interface{}) error
	Destroy(Name string) error

	GetValue(Name string) interface{}
	SetValue(Name string, Value interface{}) error

	List() []string

	Load() error
	Save() error
}

// RegisterOnConfigStart allows the registration of callbacks to be executed upon application start.
func RegisterOnConfigStart(callback func() error) {
	callbacksOnConfigStart = append(callbacksOnConfigStart, callback)
}

// OnConfigStart executes the registered callbacks upon application start.
func OnConfigStart() error {
	for _, callback := range callbacksOnConfigStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

// RegisterConfig registers the configuration upon application start.
func RegisterConfig(Config Config) error {
	if registeredConfig == nil {
		registeredConfig = Config
	} else {
		return errors.New("Configuration file already registered. Unable to register configuration.")
	}
	return nil
}

// GetConfig returns the registered configuration.
func GetConfig() Config { return registeredConfig }
