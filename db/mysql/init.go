package mysql

import (
	"database/sql"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/utils"
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
	uri := "/"
	dbName := "ottemo"

	poolConnections := 10
	maxConnections := 0

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("db.mysql.uri", uri); iniValue != "" {
			uri = iniValue
		}

		if iniValue := iniConfig.GetValue("db.mysql.db", dbName); iniValue != "" {
			dbName = iniValue
		}

		if iniValue := iniConfig.GetValue("db.mysql.maxConnections", ""); iniValue != "" {
			maxConnections = utils.InterfaceToInt(iniValue)
		}

		if iniValue := iniConfig.GetValue("db.mysql.poolConnections", ""); iniValue != "" {
			poolConnections = utils.InterfaceToInt(iniValue)
		}
	}

	if newConnection, err := sql.Open("mysql", uri); err == nil {
		it.connection = newConnection
	} else {
		return env.ErrorDispatch(err)
	}

	if (poolConnections > 0) {
		it.connection.SetMaxIdleConns(poolConnections)
	}

	if (maxConnections > 0) {
		it.connection.SetMaxOpenConns(maxConnections)
	}

	// making sure DB selected otherwise trying to obtain DB
	rows, err := it.connection.Query("SELECT DATABASE()")
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if !rows.Next() || rows.Scan(dbName) != nil || dbName == "" {
		if _, err := it.connection.Exec("USE " + dbName); err != nil {
			if _, err = it.connection.Exec("CREATE DATABASE " + dbName); err != nil {
				return env.ErrorDispatch(err)
			}
			if _, err = it.connection.Exec("USE " + dbName); err != nil {
				return env.ErrorDispatch(err)
			}
		}
	}

	// making column info table
	SQL := "CREATE TABLE IF NOT EXISTS `" + ConstCollectionNameColumnInfo + "` ( " +
		"`_id`        INTEGER NOT NULL AUTO_INCREMENT," +
		"`collection` VARCHAR(255)," +
		"`column`     VARCHAR(255)," +
		"`type`       VARCHAR(255)," +
		"`indexed`    TINYINT(1)," +
		"PRIMARY KEY(`_id`) )"

	_, err = it.connection.Exec(SQL)
	if err != nil {
		return sqlError(SQL, err)
	}

	db.OnDatabaseStart()

	return nil
}
