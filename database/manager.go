package database

import (
	"errors"
)

var currentDBEngine DBEngine

var callbacksOnDatabaseStart = []func() error{}

// RegisterOnDatabaseStart is a function to add a callback onto the callback
// chain to be executed when a database is started.
func RegisterOnDatabaseStart(callback func() error) {
	callbacksOnDatabaseStart = append(callbacksOnDatabaseStart, callback)
}

// OnDatabaseStart is a function to execute the callback chain when a database
// is started.
func OnDatabaseStart() error {
	for _, callback := range callbacksOnDatabaseStart {
		if err := callback(); err != nil {
			return err
		}
	}
	return nil
}

//RegisterDBEngine registers a new database engine for use.
func RegisterDBEngine(newEngine DBEngine) error {
	if currentDBEngine == nil {
		currentDBEngine = newEngine
	} else {
		return errors.New("Sorry, '" + currentDBEngine.GetName() + "' already registered")
	}
	return nil
}

// GetDBEngine returns the current engine.  SQLite and MongoDB are supported i
// at present.
func GetDBEngine() DBEngine {
	return currentDBEngine
}
