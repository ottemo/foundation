package sqlite

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	sqlite3 "github.com/mxk/go-sqlite/sqlite3"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// makes SQL filter string based on ColumnName, Operator and Value parameters or returns nil
//   - internal usage function for AddFilter and AddStaticFilter routines
func (it *DBCollection) makeSQLFilterString(ColumnName string, Operator string, Value interface{}) (string, error) {
	if !it.HasColumn(ColumnName) {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "51a0ae66-a5fe-4db3-9c5f-6b55f196f714", "can't find column '"+ColumnName+"'")
	}

	Operator = strings.ToUpper(Operator)
	allowedOperators := []string{"=", "!=", "<>", ">", ">=", "<", "<=", "LIKE", "IN"}

	if !utils.IsInListStr(Operator, allowedOperators) {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "793c0ec0-aa84-46cf-9305-6245d9198d45", "unknown operator '"+Operator+"' for column '"+ColumnName+"', allowed: '"+strings.Join(allowedOperators, "', ")+"'")
	}

	columnType := it.GetColumnType(ColumnName)

	// array column - special case
	if strings.HasPrefix(columnType, "[]") {
		value := strings.Trim(convertValueForSQL(Value), "'")
		template := "(', ' || `" + ColumnName + "` || ',') LIKE '%, $value,%'"

		var resultItems []string
		for _, arrayItem := range strings.Split(value, ", ") {
			item := utils.InterfaceToString(arrayItem)
			resultItems = append(resultItems, strings.Replace(template, "$value", item, 1))
		}

		if len(resultItems) == 1 {
			return resultItems[0], nil
		}
		return strings.Join(resultItems, " OR "), nil
	}

	// regular columns - default case
	switch Operator {
	case "LIKE":
		if typedValue, ok := Value.(string); ok {
			if !strings.Contains(typedValue, "%") {
				Value = "'%" + typedValue + "%'"
			} else {
				newValue := strings.Trim(Value.(string), "'")
				newValue = strings.Trim(newValue, "\"")
				Value = "'" + newValue + "'"
			}
		} else {
			Value = "''"
		}

	case "IN":
		if typedValue, ok := Value.(*DBCollection); ok {
			Value = "(" + typedValue.getSelectSQL() + ")"
		} else {
			newValue := "("
			for _, arrayItem := range utils.InterfaceToArray(Value) {
				newValue += convertValueForSQL(arrayItem) + ", "
			}
			newValue = strings.TrimRight(newValue, ", ") + ")"
			Value = newValue
		}

	default:
		Value = convertValueForSQL(Value)
	}
	return "`" + ColumnName + "` " + Operator + " " + utils.InterfaceToString(Value), nil
}

// returns SQL select statement for current collection
func (it *DBCollection) getSelectSQL() string {
	SQL := "SELECT " + it.getSQLResultColumns() + " FROM " + it.Name + it.getSQLFilters() + it.getSQLOrder() + it.Limit
	return SQL
}

// un-serialize object values
func (it *DBCollection) modifyResultRow(row sqlite3.RowMap) sqlite3.RowMap {

	for columnName, columnValue := range row {
		columnType, present := dbEngine.attributeTypes[it.Name][columnName]
		if !present {
			columnType = ""
		}

		if columnName != "_id" && columnType != "" {
			row[columnName] = db.ConvertTypeFromDbToGo(columnValue, columnType)
		}
	}

	if _, present := row["_id"]; present {
		row["_id"] = utils.InterfaceToString(row["_id"])
	}

	return row
}

// joins result columns in string
func (it *DBCollection) getSQLResultColumns() string {
	sqlColumns := "`" + strings.Join(it.ResultColumns, "`, `") + "`"
	if sqlColumns == "``" {
		sqlColumns = "*"
	}

	return sqlColumns
}

// joins order olumns in one string with preceding keyword
func (it *DBCollection) getSQLOrder() string {
	sqlOrder := strings.Join(it.Order, ", ")
	if sqlOrder != "" {
		sqlOrder = " ORDER BY " + sqlOrder
	}

	return sqlOrder
}

// collects all filters in a single string (for internal usage)
func (it *DBCollection) getSQLFilters() string {

	var collectSubfilters func(string) []string

	collectSubfilters = func(parentGroupName string) []string {
		var result []string

		for filterGroupName, filterGroup := range it.FilterGroups {
			if filterGroup.ParentGroup == parentGroupName {
				joinOperator := " AND "
				if filterGroup.OrSequence {
					joinOperator = " OR "
				}
				subFilters := collectSubfilters(filterGroupName)
				subFilters = append(subFilters, filterGroup.FilterValues...)
				filterValue := strings.Join(subFilters, joinOperator)
				if len(subFilters) > 1 {
					filterValue = "(" + filterValue + ")"
				}
				result = append(result, filterValue)
			}
		}

		return result
	}

	sqlFilters := strings.Join(collectSubfilters(""), " AND ")
	if sqlFilters != "" {
		sqlFilters = " WHERE " + sqlFilters
	}

	return sqlFilters
}

// returns filter group, creates new one if not exists
func (it *DBCollection) getFilterGroup(groupName string) *StructDBFilterGroup {
	filterGroup, present := it.FilterGroups[groupName]
	if !present {
		filterGroup = &StructDBFilterGroup{Name: groupName, FilterValues: make([]string, 0)}
		it.FilterGroups[groupName] = filterGroup
	}
	return filterGroup
}

// adds filter(combination of [column, operator, value]) in named filter group
func (it *DBCollection) updateFilterGroup(groupName string, columnName string, operator string, value interface{}) error {

	/*if !it.HasColumn(columnName) {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e9e9c8c2-39bd-48b5-9fd6-5929b9bf30f5", "not existing column " + columnName)
	}*/

	newValue, err := it.makeSQLFilterString(columnName, operator, value)
	if err != nil {
		return err
	}

	filterGroup := it.getFilterGroup(groupName)
	filterGroup.FilterValues = append(filterGroup.FilterValues, newValue)

	return nil
}

// generates new UUID for _id column
func (it *DBCollection) makeUUID(id string) string {

	if len(id) != 24 {
		timeStamp := strconv.FormatInt(time.Now().Unix(), 16)

		randomBytes := make([]byte, 8)
		if _, err := rand.Reader.Read(randomBytes); err != nil {
			_ = env.ErrorDispatch(err)
		}

		randomHex := make([]byte, 16)
		hex.Encode(randomHex, randomBytes)

		id = timeStamp + string(randomHex)
	}

	return id
}
