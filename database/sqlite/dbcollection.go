package sqlite

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	sqlite3 "code.google.com/p/go-sqlite/go1/sqlite3"
)

// SQLiteCollection holds the information related to collections in SQLite.
type SQLiteCollection struct {
	Connection *sqlite3.Conn
	TableName  string
	Columns    map[string]string

	Filters []string
}

func sqlError(SQL string, err error) error {
	return errors.New("SQL \"" + SQL + "\" error: " + err.Error())
}

// GetDBType returns a string holding the type information.  Valid types
// include: integer, float, text, char, struct, data, decimal and money.
func GetDBType(ColumnType string) (string, error) {
	ColumnType = strings.ToLower(ColumnType)
	switch {
	case ColumnType == "int" || ColumnType == "integer":
		return "INTEGER", nil
	case ColumnType == "real" || ColumnType == "float":
		return "REAL", nil
	case ColumnType == "string" || ColumnType == "text" || strings.Contains(ColumnType, "char"):
		return "TEXT", nil
	case ColumnType == "blob" || ColumnType == "struct" || ColumnType == "data":
		return "BLOB", nil
	case strings.Contains(ColumnType, "numeric") || strings.Contains(ColumnType, "decimal") || ColumnType == "money":
		return "NUMERIC", nil
	}

	return "?", errors.New("Unknown type '" + ColumnType + "'")
}

// LoadByID returns a map of IDs
func (sc *SQLiteCollection) LoadByID(id string) (map[string]interface{}, error) {
	row := make(sqlite3.RowMap)

	SQL := "SELECT * FROM " + sc.TableName + " WHERE _id = " + id
	if s, err := sc.Connection.Query(SQL); err != nil {
		return nil, err
	} else {
		if err := s.Scan(row); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}(row), nil
}

// Load returns a map based on the given filter.
func (sc *SQLiteCollection) Load() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, 10)
	var err error = nil

	row := make(sqlite3.RowMap)

	sqlLoadFilter := strings.Join(sc.Filters, " AND ")
	if sqlLoadFilter != "" {
		sqlLoadFilter = " WHERE " + sqlLoadFilter
	}

	SQL := "SELECT * FROM " + sc.TableName + sqlLoadFilter
	if s, err := sc.Connection.Query(SQL); err == nil {
		for ; err == nil; err = s.Next() {

			if err := s.Scan(row); err == nil {
				result = append(result, row)
			} else {
				return result, err
			}
		}
	} else {
		err = sqlError(SQL, err)
	}

	return result, err
}

func (sc *SQLiteCollection) Save(Item map[string]interface{}) (string, error) {

	// _id in SQLite supposed to be auto-incremented int but for MongoDB it forced to be string
	// collection interface also forced us to use string but we still want it ti be int in DB
	// to make that we need to convert it before save from  string to int or nil
	// and after save get auto-incremented id as convert to string
	if Item["_id"] != nil {
		if intValue, err := strconv.ParseInt(Item["_id"].(string), 10, 64); err == nil {
			Item["_id"] = intValue
		} else {
			Item["_id"] = nil
		}
	}

	// SQL generation
	columns := make([]string, 0, len(Item))
	args := make([]string, 0, len(Item))
	values := make([]interface{}, 0, len(Item))

	for k, v := range Item {
		columns = append(columns, "\""+k+"\"")
		args = append(args, "$_"+k)
		values = append(values, v)
	}

	SQL := "INSERT OR REPLACE INTO  " + sc.TableName +
		" (" + strings.Join(columns, ",") + ") VALUES " +
		" (" + strings.Join(args, ",") + ")"
	if err := sc.Connection.Exec(SQL, values...); err == nil {

		// auto-incremented _id back to string
		newIdInt := sc.Connection.LastInsertId()
		newIdString := strconv.FormatInt(newIdInt, 10)
		Item["_id"] = newIdString

		return newIdString, nil
	} else {
		return "", sqlError(SQL, err)
	}

	return "", nil
}

func (sc *SQLiteCollection) Delete() (int, error) {
	sqlDeleteFilter := strings.Join(sc.Filters, " AND ")
	if sqlDeleteFilter != "" {
		sqlDeleteFilter = " WHERE " + sqlDeleteFilter
	}

	SQL := "DELETE FROM " + sc.TableName + sqlDeleteFilter
	err := sc.Connection.Exec(SQL)
	affected := sc.Connection.RowsAffected()

	return affected, err
}

func (sc *SQLiteCollection) DeleteByID(id string) error {
	SQL := "DELETE FROM " + sc.TableName + " WHERE _id = " + id
	return sc.Connection.Exec(SQL)
}

func (sc *SQLiteCollection) AddFilter(ColumnName string, Operator string, Value string) error {
	if sc.HasColumn(ColumnName) {
		Operator = strings.ToUpper(Operator)
		if Operator == "" || Operator == "=" || Operator == "<>" || Operator == ">" || Operator == "<" || Operator == "LIKE" {
			sc.Filters = append(sc.Filters, ColumnName+" "+Operator+" "+Value)
		} else {
			return errors.New("unknown operator '" + Operator + "' supposed  '', '=', '>', '<', '<>', 'LIKE' " + ColumnName + "'")
		}
	} else {
		return errors.New("can't find collumn '" + ColumnName + "'")
	}

	return nil
}

func (sc *SQLiteCollection) ClearFilters() error {
	sc.Filters = make([]string, 0)
	return nil
}

// Collection columns stuff
//--------------------------

func (sc *SQLiteCollection) RefreshColumns() {
	SQL := "PRAGMA table_info(" + sc.TableName + ")"

	row := make(sqlite3.RowMap)
	for stmt, err := sc.Connection.Query(SQL); err == nil; err = stmt.Next() {
		stmt.Scan(row)

		key := row["name"].(string)
		value := row["type"].(string)
		sc.Columns[key] = value
	}
}

func (sc *SQLiteCollection) ListColumns() map[string]string {
	sc.RefreshColumns()
	return sc.Columns
}

func (sc *SQLiteCollection) HasColumn(ColumnName string) bool {
	if _, present := sc.Columns[ColumnName]; present {
		return true
	} else {
		sc.RefreshColumns()
		_, present := sc.Columns[ColumnName]
		return present
	}
}

func (sc *SQLiteCollection) AddColumn(ColumnName string, ColumnType string, indexed bool) error {

	// TODO: there probably need column name check to be only lowercase, exclude some chars, etc.

	if sc.HasColumn(ColumnName) {
		return errors.New("column '" + ColumnName + "' already exists for '" + sc.TableName + "' collection")
	}

	if ColumnType, err := GetDBType(ColumnType); err == nil {

		SQL := "ALTER TABLE " + sc.TableName + " ADD COLUMN \"" + ColumnName + "\" " + ColumnType
		if err := sc.Connection.Exec(SQL); err == nil {
			return nil
		} else {
			return sqlError(SQL, err)
		}

	} else {
		return err
	}

}

func (sc *SQLiteCollection) RemoveColumn(ColumnName string) error {

	if !sc.HasColumn(ColumnName) {
		return errors.New("column '" + ColumnName + "' not exists in '" + sc.TableName + "' collection")
	}

	sc.Connection.Begin()
	defer sc.Connection.Commit()

	var SQL string

	SQL = "SELECT sql FROM sqlite_master WHERE tbl_name='" + sc.TableName + "' AND type='table'"
	if stmt, err := sc.Connection.Query(SQL); err == nil {

		var tableCreateSQL string = ""

		if err := stmt.Scan(&tableCreateSQL); err == nil {

			tableColumnsWTypes := ""
			tableColumnsWoTypes := ""

			re := regexp.MustCompile("CREATE TABLE .*\\((.*)\\)")
			if regexMatch := re.FindStringSubmatch(tableCreateSQL); len(regexMatch) >= 2 {
				tableColumnsList := strings.Split(regexMatch[1], ", ")

				for _, tableColumn := range tableColumnsList {
					tableColumn = strings.Trim(tableColumn, "\n\t ")
					if !strings.HasPrefix(tableColumn, ColumnName) {
						if tableColumnsWTypes != "" {
							tableColumnsWTypes += ", "
							tableColumnsWoTypes += ", "
						}
						tableColumnsWTypes += "\"" + tableColumn + "\""
						tableColumnsWoTypes += "\"" + tableColumn[0:strings.Index(tableColumn, " ")] + "\""
					}

				}
			} else {
				return errors.New("can't find table create columns in '" + tableCreateSQL + "', found [" + strings.Join(regexMatch, ", ") + "]")
			}

			SQL = "CREATE TABLE " + sc.TableName + "_removecolumn (" + tableColumnsWTypes + ") "
			if err := sc.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			SQL = "INSERT INTO " + sc.TableName + "_removecolumn (" + tableColumnsWoTypes + ") SELECT " + tableColumnsWoTypes + " FROM " + sc.TableName
			if err := sc.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			// SQL = "DROP TABLE " + sc.TableName
			SQL = "ALTER TABLE " + sc.TableName + " RENAME TO " + sc.TableName + "_fordelete"
			if err := sc.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			SQL = "ALTER TABLE " + sc.TableName + "_removecolumn RENAME TO " + sc.TableName
			if err := sc.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			sc.Connection.Commit()

			// TODO: Fix this issue with table lock on DROP table

			// SQL = "DROP TABLE " + sc.TableName + "_fordelete"
			// if err := sc.Connection.Exec(SQL); err != nil {
			// 	return sqlError(SQL, err)
			// }

		} else {
			return err
		}

	} else {
		return sqlError(SQL, err)
	}

	return nil
}
