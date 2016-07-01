package env

import (
	"errors"
)

// Package global variables
var (
	// variables to hold currently registered services
	registeredConfig    InterfaceConfig
	registeredIniConfig InterfaceIniConfig
	registeredLogger    InterfaceLogger
	registeredErrorBus  InterfaceErrorBus
	registeredEventBus  InterfaceEventBus
	registeredScheduler InterfaceScheduler

	// variables to hold callback functions on configuration services startup
	callbacksOnConfigStart    = []func() error{}
	callbacksOnConfigIniStart = []func() error{}
)

// RegisterOnConfigStart registers new callback on configuration service start
func RegisterOnConfigStart(callback func() error) {
	callbacksOnConfigStart = append(callbacksOnConfigStart, callback)
}

// RegisterOnConfigIniStart registers new callback on ini configuration service start
func RegisterOnConfigIniStart(callback func() error) {
	callbacksOnConfigIniStart = append(callbacksOnConfigIniStart, callback)
}

// OnConfigStart fires config service start event (callback handling)
func OnConfigStart() error {
	for _, callback := range callbacksOnConfigStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

// OnConfigIniStart fires ini config service start event (callback handling)
func OnConfigIniStart() error {
	for _, callback := range callbacksOnConfigIniStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

// RegisterIniConfig registers ini config service in the system
//   - will cause error if there are couple candidates for that role
func RegisterIniConfig(IniConfig InterfaceIniConfig) error {
	if registeredIniConfig == nil {
		registeredIniConfig = IniConfig
	} else {
		return errors.New("There is other ini config already registered")
	}
	return nil
}

// RegisterConfig registers config service in the system
//   - will cause error if there are couple candidates for that role
func RegisterConfig(Config InterfaceConfig) error {
	if registeredConfig == nil {
		registeredConfig = Config
	} else {
		return errors.New("There is other config already registered")
	}
	return nil
}

// RegisterLogger registers logging service in the system
//   - will cause error if there are couple candidates for that role
func RegisterLogger(logger InterfaceLogger) error {
	if registeredLogger == nil {
		registeredLogger = logger
	} else {
		return errors.New("Logger already registered")
	}
	return nil
}

// RegisterEventBus registers event processor in the system
//   - will cause error if there are couple candidates for that role
func RegisterEventBus(eventBus InterfaceEventBus) error {
	if registeredEventBus == nil {
		registeredEventBus = eventBus
	} else {
		return errors.New("Event bus already registered")
	}
	return nil
}

// RegisterErrorBus registers error processor in the system
//   - will cause error if there are couple candidates for that role
func RegisterErrorBus(errorBus InterfaceErrorBus) error {
	if registeredErrorBus == nil {
		registeredErrorBus = errorBus
	} else {
		return errors.New("Error bus already registered")
	}
	return nil
}

// RegisterScheduler registers scheduler in the system
//   - will cause error if there are couple candidates for that role
func RegisterScheduler(scheduler InterfaceScheduler) error {
	if registeredScheduler == nil {
		registeredScheduler = scheduler
	} else {
		return errors.New("Scheduler already registered")
	}
	return nil
}

// GetConfig returns currently the used configuration service implementation or nil
func GetConfig() InterfaceConfig {
	return registeredConfig
}

// GetIniConfig returns currently used ini config service implementation or nil
func GetIniConfig() InterfaceIniConfig {
	return registeredIniConfig
}

// GetLogger returns currently used logging service implementation or nil
func GetLogger() InterfaceLogger {
	return registeredLogger
}

// GetErrorBus returns currently used error processor implementation or nil
func GetErrorBus() InterfaceErrorBus {
	return registeredErrorBus
}

// GetEventBus returns currently used event processor implementation or nil
func GetEventBus() InterfaceEventBus {
	return registeredEventBus
}

// GetScheduler returns currently used scheduler implementation or nil
func GetScheduler() InterfaceScheduler {
	return registeredScheduler
}

// ConfigEmptyValueValidator is a default validator function to accept any value
func ConfigEmptyValueValidator(val interface{}) (interface{}, bool) {
	return val, true
}
