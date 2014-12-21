package media

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	currentMediaStorage          InterfaceMediaStorage // currently registered media storage service in system
	callbacksOnMediaStorageStart = []func() error{}    // set of callback function on media storage service start
)

// RegisterOnMediaStorageStart registers new callback on media storage service start
func RegisterOnMediaStorageStart(callback func() error) {
	callbacksOnMediaStorageStart = append(callbacksOnMediaStorageStart, callback)
}

// OnMediaStorageStart fires media storage service start event (callback handling)
func OnMediaStorageStart() error {
	for _, callback := range callbacksOnMediaStorageStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

// RegisterMediaStorage registers media storage service in the system
//   - will cause error if there are couple candidates for that role
func RegisterMediaStorage(newEngine InterfaceMediaStorage) error {
	if currentMediaStorage == nil {
		currentMediaStorage = newEngine
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "eaf0079512f64ea98ce27cc89088af14", "Sorry, '"+currentMediaStorage.GetName()+"' media storage already registered")
	}
	return nil
}

// GetMediaStorage returns currently used media storage service implementation
func GetMediaStorage() (InterfaceMediaStorage, error) {
	if currentMediaStorage != nil {
		return currentMediaStorage, nil
	}
	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2bac26ac67c64db49d071b813288a1af", "no registered media storage")
}
