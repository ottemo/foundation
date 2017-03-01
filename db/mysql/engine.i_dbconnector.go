package mysql

import (
	"errors"
	"database/sql"
	"time"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// ------------------------------------------------------------------------------------
// InterfaceDBConnector implementation (package "github.com/ottemo/foundation/db/interfaces")
// ------------------------------------------------------------------------------------

// GetConnectionParams returns configured DB connection params
func (it *DBEngine) GetConnectionParams() interface{} {
	var connectionParams = connectionParamsType{
		uri: "/",
		dbName: "ottemo",
		poolConnections: 10,
		maxConnections: 0,
	}

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("db.mysql.uri", connectionParams.uri); iniValue != "" {
			connectionParams.uri = iniValue
		}

		if iniValue := iniConfig.GetValue("db.mysql.db", connectionParams.dbName); iniValue != "" {
			connectionParams.dbName = iniValue
		}

		if iniValue := iniConfig.GetValue("db.mysql.maxConnections", ""); iniValue != "" {
			connectionParams.maxConnections = utils.InterfaceToInt(iniValue)
		}

		if iniValue := iniConfig.GetValue("db.mysql.poolConnections", ""); iniValue != "" {
			connectionParams.poolConnections = utils.InterfaceToInt(iniValue)
		}
	}

	return connectionParams
}

// Connect establishes DB connection
func (it *DBEngine) Connect(srcConnectionParams interface{}) error {
	connectionParams, ok := srcConnectionParams.(connectionParamsType)
	if !ok {
		return errors.New("Wrong connection parameters type.")
	}

	if newConnection, err := sql.Open("mysql", connectionParams.uri); err == nil {
		it.connection = newConnection
	} else {
		return err
	}

	if connectionParams.poolConnections > 0 {
		it.connection.SetMaxIdleConns(connectionParams.poolConnections)
	}

	if connectionParams.maxConnections > 0 {
		it.connection.SetMaxOpenConns(connectionParams.maxConnections)
	}

	// making sure DB selected otherwise trying to obtain DB
	rows, err := it.connection.Query("SELECT DATABASE()")
	if err != nil {
		return err
	}

	if !rows.Next() || rows.Scan(connectionParams.dbName) != nil || connectionParams.dbName == "" {
		if _, err := it.connection.Exec("USE " + connectionParams.dbName); err != nil {
			if _, err = it.connection.Exec("CREATE DATABASE " + connectionParams.dbName); err != nil {
				return err
			}
			if _, err = it.connection.Exec("USE " + connectionParams.dbName); err != nil {
				return err
			}
		}
	}

	return nil
}

// AfterConnect makes initialization of DB engine
func (it *DBEngine) AfterConnect(srcConnectionParams interface{}) error {
	SQL := "CREATE TABLE IF NOT EXISTS `" + ConstCollectionNameColumnInfo + "` ( " +
		"`_id`        INTEGER NOT NULL AUTO_INCREMENT," +
		"`collection` VARCHAR(255)," +
		"`column`     VARCHAR(255)," +
		"`type`       VARCHAR(255)," +
		"`indexed`    TINYINT(1)," +
		"PRIMARY KEY(`_id`) )"

	_, err := it.connection.Exec(SQL)
	if err != nil {
		return sqlError(SQL, err)
	}

	return nil
}

// Reconnect tries to reconnect to DB
func (it *DBEngine) Reconnect(connectionParams interface{}) error {
	// Call Connect to "select database".
	// Simple it.Ping is not enough - it restores connection, but doesn't "select database".
	return it.Connect(connectionParams)
}

// IsConnected returns connection status
func (it *DBEngine) IsConnected() bool {
	return it.isConnected
}

// SetConnected sets connection status
func (it *DBEngine) SetConnected(connected bool) {
	it.isConnected = connected
}

// Ping checks connection alive
func (it *DBEngine) Ping() error {
	// Better way to check connection is alive and "database is selected"
	rows, err := it.connection.Query("SELECT * from " + ConstCollectionNameColumnInfo);

	defer func(rows *sql.Rows){
		if rows != nil {
			err := rows.Close()
			if err != nil {
				it.LogConnection(err.Error())
			}
		}
	}(rows)

	return err
}

// GetValidationInterval returns delay between Ping
func (it *DBEngine) GetValidationInterval() time.Duration {
	return ConstConnectionValidateInterval
}

// GetEngineName returns DBEngine name (InterfaceDBConnector)
func (it *DBEngine) GetEngineName() string {
	return it.GetName()
}

// LogConnection outputs message to log
func (it *DBEngine) LogConnection(message string) {
	env.Log("mysql.log", "DEBUG", message)
}
