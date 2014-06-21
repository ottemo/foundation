package mongodb

import (
	"errors"
	"strings"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// MongoDBCollection holds the database and collection data related to a
// MongoDB instance.
type MongoDBCollection struct {
	database   *mgo.Database
	collection *mgo.Collection
	Name       string

	Selector map[string]interface{}
}

// ColumnInfoCollection is a constant for meta information.
const (
	ColumnInfoCollection = "collection_column_info"
)

func sqlError(SQL string, err error) error {
	return errors.New("SQL \"" + SQL + "\" error: " + err.Error())
}

func getMongoOperator(Operator string) (string, error) {
	Operator = strings.ToLower(Operator)

	switch Operator {
	case "=":
		return "", nil
	case ">":
		return "gt;", nil
	case "<":
		return "lt;", nil
	case "like":
		return "like", nil
	}

	return "?", errors.New("Unknown operator '" + Operator + "'")
}

// LoadByID returns a single ID from the database.
func (mc *MongoDBCollection) LoadByID(id string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := mc.collection.FindId(id).One(&result)

	return result, err
}

// Load returns all found collections that match given selector string.
func (mc *MongoDBCollection) Load() ([]map[string]interface{}, error) {
	// result := make([]map[string]interface{}, 0)
	var result []map[string]interface{}

	err := mc.collection.Find(mc.Selector).All(&result)

	return result, err
}

// Save will create the new value in MongoDB.
func (mc *MongoDBCollection) Save(Item map[string]interface{}) (string, error) {

	id := bson.NewObjectId().Hex()

	if _id, present := Item["_id"]; present {
		if _id, ok := _id.(string); ok && _id != "" {
			if bson.IsObjectIdHex(_id) {
				id = _id
			}
		}
	}

	Item["_id"] = id

	changeInfo, err := mc.collection.UpsertId(id, Item)

	if changeInfo != nil && changeInfo.UpsertedId != nil {
		//id = changeInfo.UpsertedId
	}

	return id, err
}

// Delete will remove the given entity matching the selector.
func (mc *MongoDBCollection) Delete() (int, error) {
	changeInfo, err := mc.collection.RemoveAll(mc.Selector)

	return changeInfo.Removed, err
}

// DeleteByID will remove the entity with the given ID.
func (mc *MongoDBCollection) DeleteByID(id string) error {

	return mc.collection.RemoveId(id)
}

// AddFilter will create an additional filter column.
func (mc *MongoDBCollection) AddFilter(ColumnName string, Operator string, Value string) error {

	Operator, err := getMongoOperator(Operator)
	if err != nil {
		return err
	}

	var filterValue interface{} = Value
	if Operator != "" {
		filterValue = map[string]interface{}{Operator: Value}
	} else {
		filterValue = Value
	}

	mc.Selector[ColumnName] = filterValue

	return nil
}

// ClearFilters removes the filters for collection.
func (mc *MongoDBCollection) ClearFilters() error {
	mc.Selector = make(map[string]interface{})
	return nil
}

// ListColumns will return a map of columns.
func (mc *MongoDBCollection) ListColumns() map[string]string {

	result := map[string]string{}

	infoCollection := mc.database.C(ColumnInfoCollection)
	selector := map[string]string{"collection": mc.Name}
	iter := infoCollection.Find(selector).Iter()

	row := map[string]string{}
	for iter.Next(&row) {
		colName, okColumn := row["column"]
		colType, okType := row["type"]

		if okColumn && okType {
			result[colName] = colType
		}
	}

	return result
}

// HasColumn will return a boolean true if the column exists.
func (mc *MongoDBCollection) HasColumn(ColumnName string) bool {

	infoCollection := mc.database.C(ColumnInfoCollection)
	selector := map[string]interface{}{"collection": mc.Name, "column": ColumnName}
	count, _ := infoCollection.Find(selector).Count()

	return count > 0
}

// AddColumn will append a column for the given collection.
func (mc *MongoDBCollection) AddColumn(ColumnName string, ColumnType string, indexed bool) error {

	infoCollection := mc.database.C(ColumnInfoCollection)

	selector := map[string]interface{}{"collection": mc.Name, "column": ColumnName}
	data := map[string]interface{}{"collection": mc.Name, "column": ColumnName, "type": ColumnType, "indexed": indexed}

	_, err := infoCollection.Upsert(selector, data)

	return err
}

// RemoveColumn will remove the column of the collection.
func (mc *MongoDBCollection) RemoveColumn(ColumnName string) error {

	infoCollection := mc.database.C(ColumnInfoCollection)
	removeSelector := map[string]string{"collection": mc.Name, "column": ColumnName}

	err := infoCollection.Remove(removeSelector)
	if err != nil {
		return err
	}

	updateSelector := map[string]interface{}{ColumnName: map[string]interface{}{"$exists": true}}
	data := map[string]interface{}{"$unset": map[string]interface{}{ColumnName: ""}}

	_, err = mc.collection.UpdateAll(updateSelector, data)

	if err != nil {
		return err
	}

	return nil
}
