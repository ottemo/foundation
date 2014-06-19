package sqlite

import (
	"strings"

	sqlite3 "code.google.com/p/go-sqlite/go1/sqlite3"
	config "github.com/ottemo/foundation/config"
	database "github.com/ottemo/foundation/database"
)

var collections = map[string]database.DBCollection{}

type SQLite struct {
	Connection *sqlite3.Conn
}

func init() {
	instance := new(SQLite)

	config.RegisterOnConfigIniStart(instance.Startup)
	database.RegisterDBEngine(instance)
}

func (it *SQLite) Startup() error {

	var uri string = "ottemo.db"

	if iniConfig := config.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("db.sqlite3.uri"); iniValue != "" {
			uri = iniValue
		}
	}

	if newConnection, err := sqlite3.Open(uri); err == nil {
		it.Connection = newConnection
	} else {
		return err
	}

	database.OnDatabaseStart()

	return nil
}

func (it *SQLite) GetName() string { return "Sqlite3" }

func (it *SQLite) HasCollection(CollectionName string) bool {
	CollectionName = strings.ToLower(CollectionName)

	SQL := "SELECT name FROM sqlite_master WHERE type='table' AND name='" + CollectionName + "'"
	if _, err := it.Connection.Query(SQL); err == nil {
		return true
	} else {
		return false
	}
}

func (it *SQLite) CreateCollection(CollectionName string) error {
	CollectionName = strings.ToLower(CollectionName)

	SQL := "CREATE TABLE " + CollectionName + "(_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL)"
	if err := it.Connection.Exec(SQL); err == nil {
		return nil
	} else {
		return err
	}
}

func (it *SQLite) GetCollection(CollectionName string) (database.I_DBCollection, error) {
	CollectionName = strings.ToLower(CollectionName)

	if collection, present := collections[CollectionName]; present {
		return collection, nil
	} else {
		if !it.HasCollection(CollectionName) {
			if err := it.CreateCollection(CollectionName); err != nil {
				return nil, err
			}
		}

		//TODO: find out why: Columns: map[string]string{} - needed
		collection := &SQLiteCollection{TableName: CollectionName, Connection: it.Connection, Columns: map[string]string{}}
		collections[CollectionName] = collection

		return collection, nil
	}
}
