package sqlite

import (
	"strings"

	sqlite3 "github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"time"
)

// exec routines
func connectionExecWLastInsertID(SQL string, args ...interface{}) (int64, error) {
	dbEngine.connectionMutex.Lock()
	defer dbEngine.connectionMutex.Unlock()

	err := dbEngine.connection.Exec(SQL, args...)
	if err != nil {
		return dbEngine.connection.LastInsertId(), err
	}

	return 0, err
}

// exec routines
func connectionExecWAffected(SQL string, args ...interface{}) (int, error) {
	dbEngine.connectionMutex.Lock()
	defer dbEngine.connectionMutex.Unlock()

	err := dbEngine.connection.Exec(SQL, args...)
	if err != nil {
		return dbEngine.connection.RowsAffected(), err
	}

	return 0, err
}

// exec routines
func connectionExec(SQL string, args ...interface{}) error {
	dbEngine.connectionMutex.Lock()
	defer dbEngine.connectionMutex.Unlock()

	return dbEngine.connection.Exec(SQL, args...)
}

// query routines
func connectionQuery(SQL string) (*sqlite3.Stmt, error) {
	dbEngine.connectionMutex.Lock()
	return dbEngine.connection.Query(SQL)
}

// close sqlite3 statement routine
func closeStatement(statement *sqlite3.Stmt) {
	if statement != nil {
		statement.Close()
	}
	dbEngine.connectionMutex.Unlock()
}

// formats SQL query error for output to log
func sqlError(SQL string, err error) error {
	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "261ce31db907443ab7dce51c7dba6b52", "SQL \""+SQL+"\" error: "+err.Error())
}

// returns string that represents value for SQL query
func convertValueForSQL(value interface{}) string {

	switch value.(type) {
	case *DBCollection:
		return value.(*DBCollection).getSelectSQL()

	case bool:
		if value.(bool) {
			return "1"
		}
		return "0"

	case string:
		result := value.(string)
		result = strings.Replace(result, "'", "''", -1)
		result = strings.Replace(result, "\\", "\\\\", -1)
		result = "'" + result + "'"

		return result

	case int, int32, int64:
		return utils.InterfaceToString(value)

	case map[string]interface{}, map[string]string:
		return convertValueForSQL(utils.EncodeToJSONString(value))

	case []string, []int, []int64, []int32, []float64, []bool:
		return convertValueForSQL(utils.InterfaceToArray(value))

	case time.Time:
		return convertValueForSQL(value.(time.Time).Unix())

	case []interface{}:
		result := ""
		for _, item := range value.([]interface{}) {
			if result != "" {
				result += ", "
			}
			result += utils.InterfaceToString(item)
		}
		return convertValueForSQL(result)
	}

	return convertValueForSQL(utils.InterfaceToString(value))
}

// GetDBType returns type used inside sqlite for given general name
func GetDBType(ColumnType string) (string, error) {
	ColumnType = strings.ToLower(ColumnType)
	switch {
	case strings.HasPrefix(ColumnType, "[]"):
		return "TEXT", nil
	case ColumnType == db.ConstDBBasetypeID:
		if ConstUseUUIDids {
			return "TEXT", nil
		}
		return "INTEGER", nil
	case ColumnType == "int" || ColumnType == "integer":
		return "INTEGER", nil
	case ColumnType == "real" || ColumnType == "float":
		return "REAL", nil
	case ColumnType == "string" || ColumnType == "text" || ColumnType == "json" || strings.Contains(ColumnType, "char"):
		return "TEXT", nil
	case ColumnType == "blob" || ColumnType == "struct" || ColumnType == "data":
		return "BLOB", nil
	case strings.Contains(ColumnType, "numeric") || strings.Contains(ColumnType, "decimal") || ColumnType == "money":
		return "NUMERIC", nil
	case strings.Contains(ColumnType, "date") || strings.Contains(ColumnType, "time"):
		return "NUMERIC", nil
	case ColumnType == "bool" || ColumnType == "boolean":
		return "NUMERIC", nil
	}

	return "?", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3bc554afad7d442688c430f91c1cb151", "Unknown type '"+ColumnType+"'")
}
