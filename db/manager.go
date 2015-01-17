package db

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	currentDBEngine          InterfaceDBEngine  // currently registered database service in system
	callbacksOnDatabaseStart = []func() error{} // set of callback function on database service start
)

// RegisterOnDatabaseStart registers new callback on database service start
func RegisterOnDatabaseStart(callback func() error) {
	callbacksOnDatabaseStart = append(callbacksOnDatabaseStart, callback)
}

// OnDatabaseStart fires database service start event (callback handling)
func OnDatabaseStart() error {
	for _, callback := range callbacksOnDatabaseStart {
		if err := callback(); err != nil {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

// RegisterDBEngine registers database service in the system
//   - will cause error if there are couple candidates for that role
func RegisterDBEngine(newEngine InterfaceDBEngine) error {
	if currentDBEngine == nil {
		currentDBEngine = newEngine
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c8588998-a99b-406a-91cb-cb5766e34f7e", "Sorry, '"+currentDBEngine.GetName()+"' already registered")
	}
	return nil
}

// GetDBEngine returns currently used database service implementation
func GetDBEngine() InterfaceDBEngine {
	return currentDBEngine
}
