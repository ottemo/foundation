package sqlite

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	sqlite3 "code.google.com/p/go-sqlite/go1/sqlite3"
)

type SQLiteCollection struct {
	Connection *sqlite3.Conn
	TableName  string
	Columns    map[string]string

	Filters []string
}

func sqlError(SQL string, err error) error {
	return errors.New("SQL \"" + SQL + "\" error: " + err.Error())
}

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

func (coll *SQLiteCollection) LoadById(id string) (map[string]interface{}, error) {
	row := make(sqlite3.RowMap)

	SQL := "SELECT * FROM " + coll.TableName + " WHERE _id = " + id
	if s, err := coll.Connection.Query(SQL); err == nil {
		if err := s.Scan(row); err == nil {
			return map[string]interface{}(row), nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (coll *SQLiteCollection) Load() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, 10)
	var err error = nil

	row := make(sqlite3.RowMap)

	sqlLoadFilter := strings.Join(it.Filters, " AND ")
	if sqlLoadFilter != "" {
		sqlLoadFilter = " WHERE " + sqlLoadFilter
	}

	SQL := "SELECT * FROM " + coll.TableName + sqlLoadFilter
	if s, err := coll.Connection.Query(SQL); err == nil {
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

func (coll *SQLiteCollection) Save(Item map[string]interface{}) (string, error) {

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

	SQL := "INSERT OR REPLACE INTO  " + coll.TableName +
		" (" + strings.Join(columns, ",") + ") VALUES " +
		" (" + strings.Join(args, ",") + ")"
	if err := coll.Connection.Exec(SQL, values...); err == nil {

		// auto-incremented _id back to string
		newIdInt := coll.Connection.LastInsertId()
		newIdString := strconv.FormatInt(newIdInt, 10)
		Item["_id"] = newIdString

		return newIdString, nil
	} else {
		return "", sqlError(SQL, err)
	}

	return "", nil
}

func (coll *SQLiteCollection) Delete() (int, error) {
	sqlDeleteFilter := strings.Join(coll.Filters, " AND ")
	if sqlDeleteFilter != "" {
		sqlDeleteFilter = " WHERE " + sqlDeleteFilter
	}

	SQL := "DELETE FROM " + coll.TableName + sqlDeleteFilter
	err := coll.Connection.Exec(SQL)
	affected := coll.Connection.RowsAffected()

	return affected, err
}

func (coll *SQLiteCollection) DeleteById(id string) error {
	SQL := "DELETE FROM " + coll.TableName + " WHERE _id = " + id
	return coll.Connection.Exec(SQL)
}

func (coll *SQLiteCollection) AddFilter(ColumnName string, Operator string, Value string) error {
	if coll.HasColumn(ColumnName) {
		Operator = strings.ToUpper(Operator)
		if Operator == "" || Operator == "=" || Operator == "<>" || Operator == ">" || Operator == "<" || Operator == "LIKE" {
			coll.Filters = append(coll.Filters, ColumnName+" "+Operator+" "+Value)
		} else {
			return errors.New("unknown operator '" + Operator + "' supposed  '', '=', '>', '<', '<>', 'LIKE' " + ColumnName + "'")
		}
	} else {
		return errors.New("can't find collumn '" + ColumnName + "'")
	}

	return nil
}

func (coll *SQLiteCollection) ClearFilters() error {
	coll.Filters = make([]string, 0)
	return nil
}

// Collection columns stuff
//--------------------------

func (coll *SQLiteCollection) RefreshColumns() {
	SQL := "PRAGMA table_info(" + coll.TableName + ")"

	row := make(sqlite3.RowMap)
	for stmt, err := coll.Connection.Query(SQL); err == nil; err = stmt.Next() {
		stmt.Scan(row)

		key := row["name"].(string)
		value := row["type"].(string)
		coll.Columns[key] = value
	}
}

func (coll *SQLiteCollection) ListColumns() map[string]string {
	coll.RefreshColumns()
	return coll.Columns
}

func (coll *SQLiteCollection) HasColumn(ColumnName string) bool {
	if _, present := coll.Columns[ColumnName]; present {
		return true
	} else {
		coll.RefreshColumns()
		_, present := coll.Columns[ColumnName]
		return present
	}
}

func (coll *SQLiteCollection) AddColumn(ColumnName string, ColumnType string, indexed bool) error {

	// TODO: there probably need column name check to be only lowercase, exclude some chars, etc.

	if coll.HasColumn(ColumnName) {
		return errors.New("column '" + ColumnName + "' already exists for '" + coll.TableName + "' collection")
	}

	if ColumnType, err := GetDBType(ColumnType); err == nil {

		SQL := "ALTER TABLE " + coll.TableName + " ADD COLUMN \"" + ColumnName + "\" " + ColumnType
		if err := coll.Connection.Exec(SQL); err == nil {
			return nil
		} else {
			return sqlError(SQL, err)
		}

	} else {
		return err
	}

}

func (coll *SQLiteCollection) RemoveColumn(ColumnName string) error {

	if !coll.HasColumn(ColumnName) {
		return errors.New("column '" + ColumnName + "' not exists in '" + coll.TableName + "' collection")
	}

	coll.Connection.Begin()
	defer coll.Connection.Commit()

	var SQL string

	SQL = "SELECT sql FROM sqlite_master WHERE tbl_name='" + coll.TableName + "' AND type='table'"
	if stmt, err := coll.Connection.Query(SQL); err == nil {

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

			SQL = "CREATE TABLE " + coll.TableName + "_removecolumn (" + tableColumnsWTypes + ") "
			if err := coll.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			SQL = "INSERT INTO " + coll.TableName + "_removecolumn (" + tableColumnsWoTypes + ") SELECT " + tableColumnsWoTypes + " FROM " + coll.TableName
			if err := coll.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			// SQL = "DROP TABLE " + coll.TableName
			SQL = "ALTER TABLE " + coll.TableName + " RENAME TO " + coll.TableName + "_fordelete"
			if err := coll.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			SQL = "ALTER TABLE " + coll.TableName + "_removecolumn RENAME TO " + coll.TableName
			if err := coll.Connection.Exec(SQL); err != nil {
				return sqlError(SQL, err)
			}

			coll.Connection.Commit()

			// TODO: Fix this issue with table lock on DROP table

			// SQL = "DROP TABLE " + coll.TableName + "_fordelete"
			// if err := coll.Connection.Exec(SQL); err != nil {
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
