package sqlite

import (
	"errors"
	"time"

	"github.com/mxk/go-sqlite/sqlite3"

	"github.com/ottemo/foundation/env"
)

// ------------------------------------------------------------------------------------
// InterfaceDBConnector implementation (package "github.com/ottemo/foundation/db/interfaces")
// ------------------------------------------------------------------------------------

// GetConnectionParams returns configured DB connection params
func (it *DBEngine) GetConnectionParams() interface{} {
	var connectionParams = connectionParamsType{
		uri: "ottemo.db",
	}

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("db.sqlite3.uri", connectionParams.uri); iniValue != "" {
			connectionParams.uri = iniValue
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

	newConnection, err := sqlite3.Open(connectionParams.uri)
	if err == nil {
		// Ping
		err = newConnection.Exec("select count(*) from sqlite_master")
	}

	if err == nil {
		it.connection = newConnection
	}

	return err
}

// AfterConnect makes initialization of DB engine
func (it *DBEngine) AfterConnect(srcConnectionParams interface{}) error {
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

	return nil
}

// Ping checks connection alive
func (it *DBEngine) Ping() error {
	// This method doesn't provide 100% correct solution
	return connectionExec("select count(*) from sqlite_master")
}

// GetValidationInterval returns delay between Ping
func (it *DBEngine) GetValidationInterval() time.Duration {
	return ConstConnectionValidateInterval
}

// Reconnect tries to reconnect to DB
func (it *DBEngine) Reconnect(srcConnectionParams interface{}) error {
	connectionParams, ok := srcConnectionParams.(connectionParamsType)
	if !ok {
		return errors.New("Wrong connection parameters type.")
	}

	dbEngine.connectionMutex.Lock()
	newConnection, err := sqlite3.Open(connectionParams.uri)
	dbEngine.connectionMutex.Unlock()

	if err == nil {
		it.connection = newConnection
	}

	return err
}

// IsConnected returns connection status
func (it *DBEngine) IsConnected() bool {
	return it.isConnected
}

// SetConnected sets connection status
func (it *DBEngine) SetConnected(connected bool) {
	it.isConnected = connected
}

// GetEngineName returns DBEngine name (InterfaceDBConnector)
func (it *DBEngine) GetEngineName() string {
	return it.GetName()
}

// LogConnection outputs message to log
func (it *DBEngine) LogConnection(message string) {
	env.Log("sqlite.log", "DEBUG", message)
}
