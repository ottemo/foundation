package mongodb

import (
	"errors"
	"strings"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoDBCollection struct {
	database   *mgo.Database
	collection *mgo.Collection
	Name       string

	Selector map[string]interface{}
}

const (
	COLUMN_INFO_COLLECTION = "collection_column_info"
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

func (it *MongoDBCollection) LoadById(id string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := it.collection.FindId(id).One(&result)

	return result, err
}

func (it *MongoDBCollection) Load() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	err := it.collection.Find(it.Selector).All(&result)

	return result, err
}

func (it *MongoDBCollection) Save(Item map[string]interface{}) (string, error) {

	id := bson.NewObjectId().Hex()

	if _id, present := Item["_id"]; present {
		if _id, ok := _id.(string); ok && _id != "" {
			if bson.IsObjectIdHex(_id) {
				id = _id
			}
		}
	}

	Item["_id"] = id

	changeInfo, err := it.collection.UpsertId(id, Item)

	if changeInfo != nil && changeInfo.UpsertedId != nil {
		//id = changeInfo.UpsertedId
	}

	return id, err
}

func (it *MongoDBCollection) Delete() (int, error) {
	changeInfo, err := it.collection.RemoveAll(it.Selector)

	return changeInfo.Removed, err
}

func (it *MongoDBCollection) DeleteById(id string) error {

	return it.collection.RemoveId(id)
}

func (it *MongoDBCollection) AddFilter(ColumnName string, Operator string, Value string) error {

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

	it.Selector[ColumnName] = filterValue

	return nil
}

func (it *MongoDBCollection) ClearFilters() error {
	it.Selector = make(map[string]interface{})
	return nil
}

// Collection columns stuff
//--------------------------
func (it *MongoDBCollection) ListColumns() map[string]string {

	result := map[string]string{}

	infoCollection := it.database.C(COLUMN_INFO_COLLECTION)
	selector := map[string]string{"collection": it.Name}
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

func (it *MongoDBCollection) HasColumn(ColumnName string) bool {

	infoCollection := it.database.C(COLUMN_INFO_COLLECTION)
	selector := map[string]interface{}{"collection": it.Name, "column": ColumnName}
	count, _ := infoCollection.Find(selector).Count()

	return count > 0
}

func (it *MongoDBCollection) AddColumn(ColumnName string, ColumnType string, indexed bool) error {

	infoCollection := it.database.C(COLUMN_INFO_COLLECTION)

	selector := map[string]interface{}{"collection": it.Name, "column": ColumnName}
	data := map[string]interface{}{"collection": it.Name, "column": ColumnName, "type": ColumnType, "indexed": indexed}

	_, err := infoCollection.Upsert(selector, data)

	return err
}

func (it *MongoDBCollection) RemoveColumn(ColumnName string) error {

	infoCollection := it.database.C(COLUMN_INFO_COLLECTION)
	removeSelector := map[string]string{"collection": it.Name, "column": ColumnName}

	err := infoCollection.Remove(removeSelector)
	if err != nil {
		return err
	}

	updateSelector := map[string]interface{}{ColumnName: map[string]interface{}{"$exists": true}}
	data := map[string]interface{}{"$unset": map[string]interface{}{ColumnName: ""}}

	_, err = it.collection.UpdateAll(updateSelector, data)

	if err != nil {
		return err
	}

	return nil
}
