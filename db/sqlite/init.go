package sqlite

import (
	"github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
)

// init makes package self-initialization routine
func init() {
	dbEngine = new(DBEngine)
	dbEngine.attributeTypes = make(map[string]map[string]string)

	var _ db.InterfaceDBEngine = dbEngine

	env.RegisterOnConfigIniStart(dbEngine.Startup)
	db.RegisterDBEngine(dbEngine)
}

// Startup is a database engine startup routines
func (it *DBEngine) Startup() error {

	it.attributeTypes = make(map[string]map[string]string)

	// opening connection
	uri := "ottemo.db"

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("db.sqlite3.uri", uri); iniValue != "" {
			uri = iniValue
		}
	}

	if newConnection, err := sqlite3.Open(uri); err == nil {
		it.connection = newConnection
	} else {
		return env.ErrorDispatch(err)
	}

	// timer routine to check connection state and reconnect by perforce
	ticker := time.NewTicker(ConstConnectionValidateInterval)
	go func() {
		for _ = range ticker.C {
			err := connectionExec("select count(*) from sqlite_master")
			if err != nil {
				dbEngine.connectionMutex.Lock()
				newConnection, err := sqlite3.Open(uri)
				dbEngine.connectionMutex.Unlock()

				if err != nil {
					env.ErrorDispatch(err)
				} else {
					it.connection = newConnection
				}
			}
		}
	}()

	// making column info table
	SQL := "CREATE TABLE IF NOT EXISTS " + ConstCollectionNameColumnInfo + ` (
		_id        INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		collection VARCHAR(255),
		column     VARCHAR(255),
		type       VARCHAR(255),
		indexed    NUMERIC)`

	err := it.connection.Exec(SQL)
	if err != nil {
		return sqlError(SQL, err)
	}

	db.OnDatabaseStart()

	return nil
}
