package sqlite

import (
	"github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

var (
	dbEngine *SQLite
)

func init() {
	dbEngine = new(SQLite)
	dbEngine.attributeTypes = make(map[string]map[string]string)

	var _ db.I_DBEngine = dbEngine

	env.RegisterOnConfigIniStart(dbEngine.Startup)
	db.RegisterDBEngine(dbEngine)
}

func (it *SQLite) Startup() error {

	it.attributeTypes = make(map[string]map[string]string)

	// opening connection
	var uri string = "ottemo.db"

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

	// making column info table
	SQL := "CREATE TABLE IF NOT EXISTS " + COLLECTION_NAME_COLUMN_INFO + ` (
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
