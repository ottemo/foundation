package mysql

import (
	"strconv"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
)

// GetName returns current DB engine name
func (it *DBEngine) GetName() string {
	return "MySQL"
}

// HasCollection checks if collection(table) already exists
func (it *DBEngine) HasCollection(collectionName string) bool {
	// collectionName = strings.ToLower(collectionName)

	SQL := "SHOW TABLES LIKE '" + collectionName + "'"

	rows, err := connectionQuery(SQL)
	defer closeCursor(rows)

	if err == nil && rows.Next() {
		return true
	}

	return false
}

// CreateCollection creates cllection(table) by it's name
func (it *DBEngine) CreateCollection(collectionName string) error {
	// collectionName = strings.ToLower(collectionName)

	SQL := "CREATE TABLE " + collectionName + " (_id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY  )"
	if ConstUseUUIDids {
		SQL = "CREATE TABLE " + collectionName + " (_id CHAR(24) NOT NULL PRIMARY KEY)"
	}

	err := connectionExec(SQL)
	if err == nil {
		return nil
	}

	return env.ErrorDispatch(err)
}

// GetCollection returns collection(table) by name or creates new one
func (it *DBEngine) GetCollection(collectionName string) (db.InterfaceDBCollection, error) {
	if !ConstSQLNameValidator.MatchString(collectionName) {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fb53c672-31bb-4f4f-8ff7-e20bcdd0fcc4", "not valid collection name for DB engine")
	}

	if !it.HasCollection(collectionName) {
		if err := it.CreateCollection(collectionName); err != nil {
			return nil, env.ErrorDispatch(err)
		}
	}

	collection := &DBCollection{
		Name:          collectionName,
		FilterGroups:  make(map[string]*StructDBFilterGroup),
		Order:         make([]string, 0),
		ResultColumns: make([]string, 0),
	}

	if _, present := it.attributeTypes[collectionName]; !present {
		collection.ListColumns()
	}

	return collection, nil
}

// RawQuery returns collection(table) by name or creates new one
func (it *DBEngine) RawQuery(query string) (map[string]interface{}, error) {

	result := make([]map[string]interface{}, 0, 10)

	rows, err := connectionQuery(query)
	defer closeCursor(rows)

	if err == nil {
		return nil, env.ErrorDispatch(err)
	}

	for ok := rows.Next(); ok == false; ok = rows.Next() {
		if row, err := getRowAsStringMap(rows); err == nil {

			if ConstUseUUIDids {
				if _, present := row["_id"]; present {
					row["_id"] = strconv.FormatInt(row["_id"].(int64), 10)
				}
			}

			result = append(result, row)
		} else {
			return result[0], nil
		}
	}

	return result[0], nil
}
