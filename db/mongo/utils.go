package mongo

import (
	"gopkg.in/mgo.v2/bson"
	"sort"
	"time"

	"github.com/ottemo/foundation/utils"
)

// ConvertMapToDoc is a recursive function that converts map[string]interface{} to bson.D
// so map keys order is static and alphabetically sorted
func ConvertMapToDoc(inputMap map[string]interface{}) bson.D {
	result := make(bson.D, len(inputMap))

	// making sorted array of map keys
	//--------------------------------
	sortedKeys := make([]string, len(inputMap))
	var idx = 0
	for key := range inputMap {
		sortedKeys[idx] = key
		idx++
	}
	sort.Strings(sortedKeys)

	// converting key values to bson.DocElem
	for _, key := range sortedKeys {
		var docValue interface{} = inputMap[key]
		if mapValue, ok := docValue.(map[string]interface{}); ok {
			docValue = ConvertMapToDoc(mapValue)
		}
		result = append(result, bson.DocElem{Name: key, Value: docValue})
	}

	return result
}

// BsonDToString converts bson.D to readable form, mostly used for debug
func BsonDToString(input bson.D) string {
	result := ""

	result += "{"
	for _, bsonItem := range input {
		result += "'" + bsonItem.Name + "': "

		switch typedValue := bsonItem.Value.(type) {
		case []bson.D:
			result += "["

			addComaFlag := false
			for _, valueItem := range typedValue {
				if addComaFlag {
					result += ", "
				} else {
					addComaFlag = true
				}
				result += BsonDToString(valueItem)
			}

			result += "]"
		case bson.D:
			result += BsonDToString(typedValue)
		case time.Time:
			result += "ISODate(\"" + typedValue.Format(time.RFC3339) + "\")"
		default:
			if bsonItem.Value == nil {
				result += "null"
			}
			result += utils.InterfaceToString(typedValue)
		}
	}
	result += "}"

	return result
}
