package mongo

import (
	"fmt"
	"sort"

	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"gopkg.in/mgo.v2/bson"
)

// LoadByID loads one record from DB by record _id
func (it *DBCollection) LoadByID(id string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := it.collection.FindId(id).One(&result)
	if err != nil {
		err = fmt.Errorf("unable to load id: %v from mongodb (%v)", id, err)
	}

	return result, err
}

// Load loads records from DB for current collection and filter if it set
func (it *DBCollection) Load() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	err := it.prepareQuery().All(&result)

	return result, env.ErrorDispatch(err)
}

// Iterate applies [iterator] function to each record, stops on return false
func (it *DBCollection) Iterate(iteratorFunc func(record map[string]interface{}) bool) error {
	record := make(map[string]interface{})

	iterator := it.prepareQuery().Iter()
	for iterator.Next(&record) {
		proceed := iteratorFunc(record)

		if !proceed {
			err := iterator.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Count returns count of rows matching current select statement
func (it *DBCollection) Count() (int, error) {
	return it.collection.Find(it.makeSelector()).Count()
}

// Distinct returns distinct values of specified attribute
func (it *DBCollection) Distinct(columnName string) ([]interface{}, error) {
	var result []interface{}

	err := it.prepareQuery().Distinct(columnName, &result)

	return result, env.ErrorDispatch(err)
}

// Save stores record in DB for current collection
func (it *DBCollection) Save(Item map[string]interface{}) (string, error) {

	// id verification/updating
	//-----------------------
	id := bson.NewObjectId().Hex()

	if _id, present := Item["_id"]; present {
		if _id, ok := _id.(string); ok && _id != "" {
			if bson.IsObjectIdHex(_id) {
				id = _id
			}
		}
	}
	Item["_id"] = id

	// sorting by attribute name
	//--------------------------
	bsonDocument := make(bson.D, 0, len(Item))
	keysList := make([]string, 0, len(Item))

	for key := range Item {
		keysList = append(keysList, key)
	}
	sort.Strings(keysList)

	for _, key := range keysList {
		convertedValue := it.convertValueToType(it.GetColumnType(key), Item[key])
		bsonDocument = append(bsonDocument, bson.DocElem{Name: key, Value: convertedValue})
	}

	// saving document to DB
	//----------------------
	changeInfo, err := it.collection.UpsertId(id, bsonDocument)

	if changeInfo != nil && changeInfo.UpsertedId != nil {
		//id = changeInfo.UpsertedId
	}

	return id, env.ErrorDispatch(err)
}

// Delete removes records that matches current select statement from DB, returns amount of affected rows
func (it *DBCollection) Delete() (int, error) {
	changeInfo, err := it.collection.RemoveAll(it.makeSelector())

	return changeInfo.Removed, env.ErrorDispatch(err)
}

// DeleteByID removes record from DB by is's id
func (it *DBCollection) DeleteByID(id string) error {
	return it.collection.RemoveId(id)
}

// SetupFilterGroup setups filter group params for collection
func (it *DBCollection) SetupFilterGroup(groupName string, orSequence bool, parentGroup string) error {
	if _, present := it.FilterGroups[parentGroup]; !present && parentGroup != "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "d576838d-dda0-4986-bae6-be2e8520e3a0", "invalid parent group")
	}

	filterGroup := it.getFilterGroup(groupName)
	filterGroup.OrSequence = orSequence
	filterGroup.ParentGroup = parentGroup

	return nil
}

// RemoveFilterGroup removes filter group for collection
func (it *DBCollection) RemoveFilterGroup(GroupName string) error {
	if _, present := it.FilterGroups[GroupName]; !present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "6324cd11-6217-42a1-95e6-b53aa5f33ce3", "invalid group name")
	}
	delete(it.FilterGroups, GroupName)
	return nil
}

// AddGroupFilter adds selection filter to specific filter group (all filter groups will be joined before db query)
func (it *DBCollection) AddGroupFilter(GroupName string, ColumnName string, Operator string, Value interface{}) error {
	err := it.updateFilterGroup(GroupName, ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	return nil
}

// AddStaticFilter adds selection filter that will not be cleared by ClearFilters() function
func (it *DBCollection) AddStaticFilter(ColumnName string, Operator string, Value interface{}) error {
	err := it.updateFilterGroup(ConstFilterGroupStatic, ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	return nil
}

// AddFilter adds selection filter to current collection object
func (it *DBCollection) AddFilter(ColumnName string, Operator string, Value interface{}) error {
	err := it.updateFilterGroup(ConstFilterGroupDefault, ColumnName, Operator, Value)
	if err != nil {
		return err
	}
	return nil
}

// ClearFilters removes all filters that were set for current collection
func (it *DBCollection) ClearFilters() error {
	for filterGroup := range it.FilterGroups {
		if filterGroup != ConstFilterGroupStatic {
			delete(it.FilterGroups, filterGroup)
		}
	}
	return nil
}

// AddSort adds sorting for current collection
func (it *DBCollection) AddSort(ColumnName string, Desc bool) error {
	if Desc {
		it.Sort = append(it.Sort, "-"+ColumnName)
	} else {
		it.Sort = append(it.Sort, ColumnName)
	}
	return nil
}

// ClearSort removes any sorting that was set for current collection
func (it *DBCollection) ClearSort() error {
	it.Sort = make([]string, 0)
	return nil
}

// SetLimit results pagination
func (it *DBCollection) SetLimit(Offset int, Limit int) error {
	it.Limit = Limit
	it.Offset = Offset

	return nil
}

// SetResultColumns limits column selection for Load() and LoadByID()function
func (it *DBCollection) SetResultColumns(columns ...string) error {
	for _, columnName := range columns {
		it.ResultAttributes = []string{}

		if !it.HasColumn(columnName) {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f26abaeb-7c9a-48e8-9d3a-332b4299f9c7", "there is no column "+columnName+" found")
		}

		it.ResultAttributes = append(it.ResultAttributes, columnName)
	}
	return nil
}

// ListColumns returns attributes available for current collection
func (it *DBCollection) ListColumns() map[string]string {

	result := map[string]string{}

	infoCollection := it.database.C(ConstCollectionNameColumnInfo)
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

	// updating cached attribute types information
	if _, present := attributeTypes[it.Name]; !present {
		attributeTypes[it.Name] = make(map[string]string)
	}

	attributeTypesMutex.Lock()
	for attributeName, attributeType := range result {
		attributeTypes[it.Name][attributeName] = attributeType
	}
	attributeTypesMutex.Unlock()

	return result
}

// GetColumnType returns SQL like type of attribute in current collection, or if not present ""
func (it *DBCollection) GetColumnType(ColumnName string) string {
	// _id - has static type
	if ColumnName == "_id" {
		return db.ConstTypeID
	}

	// looking in cache first
	attributeType, present := attributeTypes[it.Name][ColumnName]
	if !present {
		// updating cache, and looking again
		it.ListColumns()
		attributeType, present = attributeTypes[it.Name][ColumnName]
	}

	return attributeType
}

// HasColumn checks attribute presence in current collection
func (it *DBCollection) HasColumn(ColumnName string) bool {
	// _id - always present
	if ColumnName == "_id" {
		return true
	}

	// looking in cache first
	attributeTypesMutex.Lock()
	_, present := attributeTypes[it.Name][ColumnName]
	if !present {
		// updating cache, and looking again
		it.ListColumns()
		_, present = attributeTypes[it.Name][ColumnName]
	}
	attributeTypesMutex.Unlock()

	return present
}

// AddColumn adds new attribute to current collection
func (it *DBCollection) AddColumn(ColumnName string, ColumnType string, indexed bool) error {

	infoCollection := it.database.C(ConstCollectionNameColumnInfo)

	selector := map[string]interface{}{"collection": it.Name, "column": ColumnName}
	data := map[string]interface{}{"collection": it.Name, "column": ColumnName, "type": ColumnType, "indexed": indexed}

	_, err := infoCollection.Upsert(selector, data)

	return env.ErrorDispatch(err)
}

// RemoveColumn removes attribute from current collection
//   - for MoongoDB it means update "collection_column_info" collection
//   and, update all objects of current collection to exclude attribute
func (it *DBCollection) RemoveColumn(ColumnName string) error {

	infoCollection := it.database.C(ConstCollectionNameColumnInfo)
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

	it.ListColumns()

	return nil
}
