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
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7f379d36-eb93-4add-b1ee-df2c3a35a590", "Can't get DBEngine")
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

		stringValue := utils.InterfaceToString(value)
		if stringValue == "" {
			return result
		}

		jsonArray, err := utils.DecodeJSONToArray(stringValue)
		if err == nil {
			return jsonArray
		}

		array := strings.Split(stringValue, ", ")
		arrayType := strings.TrimPrefix(valueType, "[]")

		for _, arrayValue := range array {
			arrayValue = strings.Replace(arrayValue, "#2C;", ",", -1)
			result = append(result, ConvertTypeFromDbToGo(arrayValue, arrayType))
		}
		return result

	case strings.HasPrefix(valueType, ConstTypeBoolean):
		return utils.InterfaceToBool(value)

	case strings.HasPrefix(valueType, ConstTypeInteger):
		return utils.InterfaceToInt(value)

	case strings.HasPrefix(valueType, ConstTypeDecimal),
		strings.HasPrefix(valueType, ConstTypeFloat),
		strings.HasPrefix(valueType, ConstTypeMoney):

		return utils.InterfaceToFloat64(value)

	case valueType == ConstTypeDatetime:
		return utils.InterfaceToTime(value)

	case valueType == ConstTypeJSON:
		result, _ := utils.DecodeJSONToStringKeyMap(value)
		return result

	case strings.HasPrefix(valueType, ConstTypeVarchar), valueType == ConstTypeText, valueType == ConstTypeID:
		return utils.InterfaceToString(value)

	}

	return value
}

// TypeParse shortcut for utils.DataTypeParse
func TypeParse(typeName string) utils.DataType {
	return utils.DataTypeParse(typeName)
}

// TypeWPrecisionAndScale shortcut for utils.DataTypeWPrecisionAndScale
func TypeWPrecisionAndScale(dataType string, precision int, scale int) string {
	return utils.DataTypeWPrecisionAndScale(dataType, precision, scale)
}

// TypeWPrecision shortcut for utils.DataTypeWPrecision
func TypeWPrecision(dataType string, precision int) string {
	return utils.DataTypeWPrecision(dataType, precision)
}

// TypeArrayOf shortcut for utils.DataTypeArrayOf
func TypeArrayOf(dataType string) string {
	return utils.DataTypeArrayOf(dataType)
}

// TypeIsArray shortcut for utils.DataTypeIsArray
func TypeIsArray(dataType string) bool {
	return utils.DataTypeIsArray(dataType)
}

// TypeIsString shortcut for utils.DataTypeIsString
func TypeIsString(dataType string) bool {
	return utils.DataTypeIsString(dataType)
}

// TypeIsFloat shortcut for utils.DataTypeIsFloat
func TypeIsFloat(dataType string) bool {
	return utils.DataTypeIsFloat(dataType)
}
