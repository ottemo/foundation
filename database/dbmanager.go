package database

import (
	"errors"
)

// DBEngine is the interface that holds the necessary methods required for a database engine.
type DBEngine interface {
	GetName() string

	CreateCollection(Name string) error
	GetCollection(Name string) (DBCollection, error)
	HasCollection(Name string) bool
}

// DBCollection is the interface that contains all the necessary methods for a
// collection when working with extendable objects.
type DBCollection interface {
	Load() ([]map[string]interface{}, error)
	LoadByID(id string) (map[string]interface{}, error)

	Save(map[string]interface{}) (string, error)

	Delete() (int, error)
	DeleteByID(id string) error

	AddFilter(ColumnName string, Operator string, Value string) error
	ClearFilters() error

	ListColumns() map[string]string
	HasColumn(ColumnName string) bool

	AddColumn(ColumnName string, ColumnType string, indexed bool) error
	RemoveColumn(ColumnName string) error
}

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
