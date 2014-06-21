package sqlite

import (
	"strings"

	sqlite3 "code.google.com/p/go-sqlite/go1/sqlite3"
	config "github.com/ottemo/foundation/config"
	database "github.com/ottemo/foundation/database"
)

var collections = map[string]database.DBCollection{}

// SQLite structure holds the connection string for the database.
type SQLite struct {
	Connection *sqlite3.Conn
}

func init() {
	instance := new(SQLite)

	config.RegisterOnConfigIniStart(instance.Startup)
	database.RegisterDBEngine(instance)
}

// Startup opens the database connection using the connection string found in
// the INI configuration file.
func (sc *SQLite) Startup() error {

	var uri = "ottemo.db"

	if iniConfig := config.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("db.sqlite3.uri"); iniValue != "" {
			uri = iniValue
		}
	}

	if newConnection, err := sqlite3.Open(uri); err == nil {
		sc.Connection = newConnection
	} else {
		return err
	}

	database.OnDatabaseStart()

	return nil
}

// GetName returns the string for SQLite3.
func (sc *SQLite) GetName() string { return "Sqlite3" }

// HasCollection returns a boolean if the given collection is in the database.
func (sc *SQLite) HasCollection(CollectionName string) bool {
	CollectionName = strings.ToLower(CollectionName)

	SQL := "SELECT name FROM sqlite_master WHERE type='table' AND name='" + CollectionName + "'"
	if _, err := sc.Connection.Query(SQL); err == nil {
		return true
	}
	return false
}

// CreateCollection will return nil if the collection has been created
// successfully or an error if it cannot.
func (sc *SQLite) CreateCollection(CollectionName string) error {
	CollectionName = strings.ToLower(CollectionName)

	SQL := "CREATE TABLE " + CollectionName + "(_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL)"
	if err := sc.Connection.Exec(SQL); err != nil {
		return err
	}
	return nil
}

// GetCollection returns a DBCollection
func (sc *SQLite) GetCollection(CollectionName string) (database.DBCollection, error) {
	CollectionName = strings.ToLower(CollectionName)

	if collection, present := collections[CollectionName]; present {
		return collection, nil
	}
	if !sc.HasCollection(CollectionName) {
		if err := sc.CreateCollection(CollectionName); err != nil {
			return nil, err
		}
	}

	//TODO: find out why: Columns: map[string]string{} - needed
	collection := &SQLiteCollection{TableName: CollectionName, Connection: sc.Connection, Columns: map[string]string{}}
	collections[CollectionName] = collection

	return collection, nil
}
