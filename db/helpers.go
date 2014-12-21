package db

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// GetCollection returns database collection or error otherwise
func GetCollection(CollectionName string) (InterfaceDBCollection, error) {
	dbEngine := GetDBEngine()
	if dbEngine == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7f379d36eb934addb1eedf2c3a35a590", "Can't get DBEngine")
	}

	return dbEngine.GetCollection(CollectionName)
}

// ConvertTypeFromDbToGo returns object that represents GO side value for given valueType
func ConvertTypeFromDbToGo(value interface{}, valueType string) interface{} {
	switch {
	case strings.HasPrefix(valueType, "[]"):
		var result []interface{}
		if value == nil {
			return result
		}

		valueString := utils.InterfaceToString(value)
		if valueString == "" {
			return result
		}

		array := strings.Split(valueString, ", ")
		arrayType := strings.TrimPrefix(valueType, "[]")

		for _, arrayValue := range array {
			result = append(result, ConvertTypeFromDbToGo(arrayValue, arrayType))
		}
		return result

	case strings.HasPrefix(valueType, ConstDBBasetypeBoolean):
		return utils.InterfaceToBool(value)

	case strings.HasPrefix(valueType, ConstDBBasetypeInteger):
		return utils.InterfaceToInt(value)

	case strings.HasPrefix(valueType, ConstDBBasetypeDecimal),
		strings.HasPrefix(valueType, ConstDBBasetypeFloat),
		strings.HasPrefix(valueType, ConstDBBasetypeMoney):

		return utils.InterfaceToFloat64(value)

	case valueType == ConstDBBasetypeDatetime:
		return utils.InterfaceToTime(value)

	case valueType == ConstDBBasetypeJSON:
		result, _ := utils.DecodeJSONToStringKeyMap(value)
		return result

	case strings.HasPrefix(valueType, ConstDBBasetypeVarchar), valueType == ConstDBBasetypeText, valueType == ConstDBBasetypeID:
		return utils.InterfaceToString(value)

	}

	return value
}
