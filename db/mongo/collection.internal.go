package mongo

import (
	"sort"
	"strings"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// converts value from GO representation to DB before usage in queries
func (it *DBCollection) convertValueToType(columnType string, value interface{}) interface{} {

	if utils.DataTypeIsArray(columnType) {
		value = utils.InterfaceToArray(value)
	}

	switch typedValue := value.(type) {
	case []interface{}:
		columnType = utils.DataTypeArrayBaseType(columnType)
		for idx, listValue := range typedValue {
			typedValue[idx] = it.convertValueToType(columnType, listValue)
		}
	default:
		switch {
		case columnType == "id" || columnType == "string" || columnType == "text" || strings.Contains(columnType, "char"):
			return utils.InterfaceToString(value)
		case columnType == "int" || columnType == "integer":
			return utils.InterfaceToInt(value)
		case columnType == "real" || columnType == "float":
			return utils.InterfaceToFloat64(value)
		case strings.Contains(columnType, "numeric") || strings.Contains(columnType, "decimal") || columnType == "money":
			return utils.InterfaceToFloat64(value)
		case strings.Contains(columnType, "time") || strings.Contains(columnType, "date"):
			return utils.InterfaceToTime(value)
		case columnType == "bool" || columnType == "boolean":
			return utils.InterfaceToBool(value)
		}
	}

	return value
}

// converts known SQL filter operator to mongoDB one, also modifies value if needed
func (it *DBCollection) getSelectorValue(columnName string, operator string, value interface{}) (interface{}, error) {
	operator = strings.ToLower(operator)

	columnType := it.GetColumnType(columnName)

	switch operator {
	case "=":
		return it.convertValueToType(columnType, value), nil
	case "!=", "<>":
		return bson.D{bson.DocElem{Name: "$ne", Value: it.convertValueToType(columnType, value)}}, nil
	case ">":
		return bson.D{bson.DocElem{Name: "$gt", Value: it.convertValueToType(columnType, value)}}, nil
	case ">=":
		return bson.D{bson.DocElem{Name: "$gte", Value: it.convertValueToType(columnType, value)}}, nil
	case "<":
		return bson.D{bson.DocElem{Name: "$lt", Value: it.convertValueToType(columnType, value)}}, nil
	case "<=":
		return bson.D{bson.DocElem{Name: "$lte", Value: it.convertValueToType(columnType, value)}}, nil
	case "like":
		stringValue := utils.InterfaceToString(value)
		stringValue = strings.Replace(stringValue, "%", ".*", -1)
		stringValue = utils.EscapeRegexSpecials(stringValue)
		return bson.D{bson.DocElem{Name: "$regex", Value: stringValue}, bson.DocElem{Name: "$options", Value: "i"}}, nil

	case "in", "nin":
		newOperator := "$" + operator

		switch typedValue := value.(type) {
		case *DBCollection:
			refValue := new(bson.Raw)

			if len(typedValue.ResultAttributes) != 1 {
				typedValue.ResultAttributes = []string{"_id"}
			}

			if it.subcollections == nil {
				it.subcollections = make([]*DBCollection, 0)
			}

			if it.subresults == nil {
				it.subresults = make([]*bson.Raw, 0)
			}

			it.subcollections = append(it.subcollections, typedValue)
			it.subresults = append(it.subresults, refValue)

			return bson.D{bson.DocElem{Name: newOperator, Value: refValue}}, nil
		default:
			if !utils.DataTypeIsArray(columnType) {
				value = it.convertValueToType(utils.DataTypeArrayOf(columnType), value)
			}
			return bson.D{bson.DocElem{Name: newOperator, Value: value}}, nil
		}
	}

	return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ab9d771a-eb9e-4318-904e-ad13171dbe84", "Unknown operator '"+operator+"'")
}

// returns filter group, creates new one if not exists
func (it *DBCollection) getFilterGroup(groupName string) *StructDBFilterGroup {
	filterGroup, present := it.FilterGroups[groupName]
	if !present {
		filterGroup = &StructDBFilterGroup{Name: groupName, FilterValues: make([]bson.D, 0)}
		it.FilterGroups[groupName] = filterGroup
	}
	return filterGroup
}

// adds filter(combination of [column, operator, value]) in named filter group
func (it *DBCollection) updateFilterGroup(groupName string, columnName string, operator string, value interface{}) error {

	/*if !it.HasColumn(columnName) {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c7f88f11-a136-459d-b7fd-c2773a095f71", "not existing column " + columnName)
	}*/

	// converting operator and value for mongoDB usage
	//-------------------------------------------------
	newValue, err := it.getSelectorValue(columnName, operator, value)
	if err != nil {
		return err
	}

	// adding filter with converted operator/value to filter group
	//------------------------------------------------------------
	newFilter := bson.D{bson.DocElem{Name: columnName, Value: newValue}}

	filterGroup := it.getFilterGroup(groupName)
	filterGroup.FilterValues = append(filterGroup.FilterValues, newFilter)

	return nil
}

// joins filters groups in one selector
func (it *DBCollection) makeSelector() bson.D {

	// making sorted array of filter groups
	//-------------------------------------
	sortedFilterGroupsNames := make([]string, len(it.FilterGroups))
	idx := 0
	for groupName := range it.FilterGroups {
		sortedFilterGroupsNames[idx] = groupName
		idx++
	}
	sort.Strings(sortedFilterGroupsNames)

	// making recursive groups injects, based on Parent field
	//-------------------------------------------------------
	topLevelGroup := &StructDBFilterGroup{Name: "", FilterValues: make([]bson.D, 0)}
	var groupsStack []*StructDBFilterGroup
	currentGroup := topLevelGroup

	for {

		childFound := false
		// loop over sorted filter group names
		for idx, filterGroupName := range sortedFilterGroupsNames {
			if filterGroupName == "" {
				continue
			}

			iterationFilterGroup := it.FilterGroups[filterGroupName]

			// looking for child groups, making stack on them
			//-----------------------------------------------
			if iterationFilterGroup.ParentGroup == currentGroup.Name {
				groupsStack = append(groupsStack, currentGroup)
				currentGroup = iterationFilterGroup

				// excluding group filter from our list
				sortedFilterGroupsNames[idx] = ""

				childFound = true
				break
			}
		}

		// no child found for currentGroup, collapsing stack for one level
		//----------------------------------------------------------------
		if childFound == false {

			// making document from StructDBFilterGroup before pop stack
			joinOperator := "$and"
			if currentGroup.OrSequence {
				joinOperator = "$or"
			}
			bsonDoc := bson.D{bson.DocElem{Name: joinOperator, Value: currentGroup.FilterValues}}

			// popping stack - moving level down for one level, if possible
			lastIndex := len(groupsStack) - 1
			if lastIndex >= 0 {
				currentGroup = groupsStack[lastIndex]
				groupsStack = groupsStack[0:lastIndex]
			} else {
				break
			}

			// appending top level child to parent
			currentGroup.FilterValues = append(currentGroup.FilterValues, bsonDoc)
		}
	}

	if len(topLevelGroup.FilterValues) > 0 {
		return bson.D{bson.DocElem{Name: "$and", Value: topLevelGroup.FilterValues}}
	}

	return bson.D{}
}

// returns bson.Query struct with applied Sort, Offset, Limit parameters, and executed subqueries
func (it *DBCollection) prepareQuery() *mgo.Query {
	selector := it.makeSelector()
	if ConstMongoDebug {
		env.Log("mongo.log", env.ConstLogPrefixDebug, it.Name+": "+BsonDToString(selector))
	}
	query := it.collection.Find(selector)

	if len(it.Sort) > 0 {
		query.Sort(it.Sort...)
	}

	if it.Offset > 0 {
		query = query.Skip(it.Offset)
	}
	if it.Limit > 0 {
		query = query.Limit(it.Limit)
	}

	for idx, subCollection := range it.subcollections {
		subCollection.prepareQuery().Distinct(subCollection.ResultAttributes[0], it.subresults[idx])
	}

	return query
}
