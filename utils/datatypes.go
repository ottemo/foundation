package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// StaticTypeRegexp is a regular expression used to parse datatype
	StaticTypeRegexp = regexp.MustCompile(`^\s*(\[\])?(\w+)\s*(?:\(\s*(\d+)?\s*(?:\)|,\s*(\d+)\s*\)))$`)
)

// set of known data types
const (
	ConstDataTypeID       = "id"
	ConstDataTypeBoolean  = "bool"
	ConstDataTypeVarchar  = "varchar"
	ConstDataTypeText     = "text"
	ConstDataTypeInteger  = "int"
	ConstDataTypeDecimal  = "decimal"
	ConstDataTypeMoney    = "money"
	ConstDataTypeFloat    = "float"
	ConstDataTypeDatetime = "datetime"
	ConstDataTypeJSON     = "json"
)

// DataType represents data type details
type DataType struct {
	Name      string
	Precision int
	Scale     int
	IsArray   bool
	IsKnown   bool
}

// String makes a string value representation of current type
func (it *DataType) String() string {
	var result string

	if it.IsArray {
		result += "[]"
	}

	result += it.Name

	if it.Precision > 0 {
		result += "(" + strconv.Itoa(it.Precision)

		if it.Scale > 0 {
			result += ", " + strconv.Itoa(it.Scale)
		}

		result += ")"
	}

	return result
}

// DataTypeIsFloat returns true if dataType representation for GO language is float64 type
func DataTypeIsFloat(dataType string) bool {
	parsedType := DataTypeParse(dataType)
	if IsAmongStr(parsedType.Name, ConstDataTypeFloat, ConstDataTypeDecimal, ConstDataTypeMoney) {
		return true
	}

	return false
}

// DataTypeIsString returns true if dataType representation for GO language is string type
func DataTypeIsString(dataType string) bool {
	parsedType := DataTypeParse(dataType)
	if IsAmongStr(parsedType.Name, ConstDataTypeVarchar, ConstDataTypeText, ConstDataTypeJSON) {
		return true
	}

	return false
}

// DataTypeIsArray returns true if dataType is kind of array
func DataTypeIsArray(dataType string) bool {
	parsedType := DataTypeParse(dataType)
	return parsedType.IsArray
}

// DataTypeArrayOf adds array modifier to given dataType, returns "" for unknown types
func DataTypeArrayOf(dataType string) string {
	if !IsAmongStr(dataType, ConstDataTypeID, ConstDataTypeBoolean, ConstDataTypeVarchar, ConstDataTypeText, ConstDataTypeInteger,
		ConstDataTypeDecimal, ConstDataTypeMoney, ConstDataTypeFloat, ConstDataTypeDatetime, ConstDataTypeJSON) {

		return ""
	}

	return "[]" + dataType
}

// DataTypeWPrecision adds precision modifier to given dataType, returns "" for unknown types
func DataTypeWPrecision(dataType string, precision int) string {
	if !IsAmongStr(dataType, ConstDataTypeID, ConstDataTypeBoolean, ConstDataTypeVarchar, ConstDataTypeText, ConstDataTypeInteger,
		ConstDataTypeDecimal, ConstDataTypeMoney, ConstDataTypeFloat, ConstDataTypeDatetime, ConstDataTypeJSON) {

		return ""
	}

	return dataType + "(" + strconv.Itoa(precision) + ")"
}

// DataTypeWPrecisionAndScale adds precision and scale modifier to given dataType, returns "" for unknown types
func DataTypeWPrecisionAndScale(dataType string, precision int, scale int) string {
	if !IsAmongStr(dataType, ConstDataTypeID, ConstDataTypeBoolean, ConstDataTypeVarchar, ConstDataTypeText, ConstDataTypeInteger,
		ConstDataTypeDecimal, ConstDataTypeMoney, ConstDataTypeFloat, ConstDataTypeDatetime, ConstDataTypeJSON) {

		return ""
	}

	return dataType + "(" + strconv.Itoa(precision) + "," + strconv.Itoa(scale) + ")"
}

// DataTypeParse tries to parse given string representation of datatype into DataType struct
func DataTypeParse(typeName string) DataType {
	var result DataType

	typeName = strings.TrimSpace(typeName)

	regexpGroups := StaticTypeRegexp.FindStringSubmatch(typeName)

	if len(regexpGroups) != 0 {
		result.Name = typeName
	} else {
		result.IsArray = !(regexpGroups[1] == "")
		result.Name = strings.ToLower(regexpGroups[2])
		result.Precision, _ = StringToInteger(regexpGroups[3])
		result.Scale, _ = StringToInteger(regexpGroups[4])
	}

	if !IsAmongStr(result.Name, ConstDataTypeID, ConstDataTypeBoolean, ConstDataTypeVarchar, ConstDataTypeText, ConstDataTypeInteger,
		ConstDataTypeDecimal, ConstDataTypeMoney, ConstDataTypeFloat, ConstDataTypeDatetime, ConstDataTypeJSON) {

		switch {
		case IsAmongStr(result.Name, "b", "boolean"):
			result.Name = ConstDataTypeBoolean
			result.IsKnown = true

		case IsAmongStr(result.Name, "i", "integer", "single"):
			result.Name = ConstDataTypeInteger
			result.IsKnown = true

		case IsAmongStr(result.Name, "f", "d", "flt", "dbl", "float", "double"):
			result.Name = ConstDataTypeFloat
			result.IsKnown = true

		case IsAmongStr(result.Name, "str", "string"):
			result.Name = ConstDataTypeText
			result.IsKnown = true

		case IsAmongStr(result.Name, "time", "date", "calendar"):
			result.Name = ConstDataTypeDatetime
			result.IsKnown = true

		case IsAmongStr(result.Name, "struct"):
			result.Name = ConstDataTypeJSON
			result.IsKnown = true
		}
	}

	return result
}

// CheckIsBlank checks if value is blank (zero value)
func CheckIsBlank(value interface{}) bool {
	switch typedValue := value.(type) {
	case string:
		return typedValue == ""
	case int, int32, int64:
		return typedValue == 0
	case []string, []interface{}:
		return false
	case map[string]interface{}, map[string]string, map[string]int:
		return false
	case float32, float64:
		return typedValue == 0
	}

	return false
}

// IsMD5 checks if value is MD5
func IsMD5(value string) bool {
	ok := false
	if len(value) == 32 {
		ok = true
		for i := 0; i < 32; i++ {
			c := value[i]
			if !(c == '1' || c == '2' || c == '3' || c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9' || c == '0' ||
				c == 'a' || c == 'b' || c == 'c' || c == 'd' || c == 'e' || c == 'f') {
				ok = false
				break
			}
		}
	}
	return ok
}

// InterfaceToBool converts interface{} to string
func InterfaceToBool(value interface{}) bool {
	switch typedValue := value.(type) {
	case bool:
		return typedValue
	case string:
		typedValue = strings.TrimSpace(typedValue)
		if boolValue, err := strconv.ParseBool(typedValue); err == nil {
			return boolValue
		}

		typedValue = strings.ToLower(typedValue)
		if IsAmongStr(typedValue, "1", "ok", "on", "yes", "y", "x", "+") {
			return true
		}

		if intValue, err := strconv.Atoi(typedValue); err == nil && intValue > 0 {
			return true
		}

		return false
	case int:
		return typedValue > 0
	case int32:
		return typedValue > 0
	case int64:
		return typedValue > 0
	default:
		return false
	}
}

// InterfaceToArray converts interface{} to map[string]interface{}
func InterfaceToArray(value interface{}) []interface{} {
	var result []interface{}

	switch typedValue := value.(type) {
	case []string:
		result = make([]interface{}, len(typedValue))
		for idx, value := range typedValue {
			result[idx] = value
		}
	case []int:
		result = make([]interface{}, len(typedValue))
		for idx, value := range typedValue {
			result[idx] = value
		}
	case []int64:
		result = make([]interface{}, len(typedValue))
		for idx, value := range typedValue {
			result[idx] = value
		}
	case []float64:
		result = make([]interface{}, len(typedValue))
		for idx, value := range typedValue {
			result[idx] = value
		}
	case []interface{}:
		return typedValue

	case string:
		jsonArray, err := DecodeJSONToArray(typedValue)
		if err == nil {
			return jsonArray
		}

		splitValues := strings.Split(typedValue, ",")
		result = make([]interface{}, len(splitValues))
		for idx, value := range splitValues {
			result[idx] = strings.Trim(value, " \t\n")
		}

	default:
		result = append(result, value)
	}

	return result
}

// InterfaceToMap converts interface{} to map[string]interface{}
func InterfaceToMap(value interface{}) map[string]interface{} {
	switch typedValue := value.(type) {
	case map[string]interface{}:
		return typedValue

	case string:
		value, err := DecodeJSONToStringKeyMap(value)
		if err == nil {
			return value
		}
	}

	return make(map[string]interface{})
}

// InterfaceToString converts interface{} to string
func InterfaceToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch value := value.(type) {
	case string:
		return value
	case bool:
		return strconv.FormatBool(value)
	case int:
		return strconv.Itoa(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case float64:
		return strconv.FormatFloat(value, 'f', 6, 64)
	default:
		return EncodeToJSONString(value)
	}
}

// InterfaceToInt converts interface{} to integer
func InterfaceToInt(value interface{}) int {
	if value == nil {
		return 0
	}

	switch typedValue := value.(type) {
	case int:
		return typedValue
	case int32:
		return int(typedValue)
	case int64:
		return int(typedValue)
	case float32:
		return int(typedValue)
	case float64:
		return int(typedValue)
	case string:
		intValue, err := strconv.ParseInt(typedValue, 10, 64)
		if err != nil {
			floatValue, err := strconv.ParseFloat(typedValue, 64)
			if err == nil {
				return int(floatValue)
			}
		}
		return int(intValue)
	default:
		return 0
	}
}

// InterfaceToFloat64 converts interface{} to float64
func InterfaceToFloat64(value interface{}) float64 {
	if value == nil {
		return 0
	}

	switch typedValue := value.(type) {
	case float64:
		return typedValue
	case int64:
		return float64(typedValue)
	case int:
		return float64(typedValue)
	case string:
		floatValue, _ := strconv.ParseFloat(typedValue, 64)
		return floatValue
	default:
		return 0.0
	}
}

// InterfaceToTime converts interface{} to time.Time
func InterfaceToTime(value interface{}) time.Time {
	switch typedValue := value.(type) {
	case int64:
		return time.Unix(value.(int64), 0)
	case time.Time:
		return typedValue
	case string:

		tryFirst := []string{time.RFC3339, time.UnixDate}

		for _, currentFormat := range tryFirst {
			newValue, err := time.Parse(currentFormat, typedValue)
			if err == nil {
				return newValue
			}
		}

		dateFormats := []string{"02/01/2006", "02/01/06", "2006-01-02", "2006-Jan-_2", "_2 Jan 2006", "01.02.2006"}
		timeFormats := []string{"", " 3:04PM", " 15:04", " 15:04:05", "T15:04:05"}
		zoneFormats := []string{"",  " MST-07", " MST", " MST-0700", " -0700", "Z07:00"}

		for _, zoneFormat := range zoneFormats {
			for _, timeFormat := range timeFormats {
				for _, dateFormat := range dateFormats {
					currentFormat := dateFormat + timeFormat + zoneFormat
					newValue, err := time.Parse(currentFormat, typedValue)
					if err == nil {
						return newValue
					}
				}
			}
		}
	}

	return time.Time{} //Unix(0, 0)
}

// IsZeroTime checks time for zero value
func IsZeroTime(value time.Time) bool {
	zeroTime := time.Unix(0, 0)
	return value == zeroTime
}

// StringToInteger converts string to integer
func StringToInteger(value string) (int, error) {
	return strconv.Atoi(value)
}

// StringToFloat converts string to float64
func StringToFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// StringToType converts string coded value to type
func StringToType(value string, valueType string) (interface{}, error) {
	valueType = strings.ToLower(valueType)

	switch {
	case strings.HasPrefix(valueType, "[]"):
		value = strings.Trim(value, "[] \t\n")
		array := strings.Split(value, ",")
		arrayType := strings.TrimPrefix(valueType, "[]")

		var result []interface{}
		for _, arrayValue := range array {
			arrayValue = strings.TrimSpace(arrayValue)

			if arrayValue, err := StringToType(arrayValue, arrayType); err == nil {
				result = append(result, arrayValue)
			} else {
				return nil, err
			}
		}
		return result, nil

	case IsAmongStr(valueType, "b", "bool", "boolean"):
		return InterfaceToBool(value), nil

	case IsAmongStr(valueType, "i", "int", "integer"):
		return InterfaceToInt(value), nil

	case IsAmongStr(valueType, "f", "d", "flt", "dbl", "float", "double", "decimal"):
		return InterfaceToFloat64(value), nil

	case IsAmongStr(valueType, "str", "string"):
		return value, nil

	case IsAmongStr(valueType, "time", "date", "datetime"):
		return InterfaceToTime(value), nil

	case IsAmongStr(valueType, "json"):
		return DecodeJSONToStringKeyMap(value)
	}

	return nil, errors.New("unknown value type " + valueType)
}

// StringToInterface converts string to Interface{} which can be float, int, bool, string, or json as map[string]value
func StringToInterface(value string) interface{} {

	trimmedValue := strings.TrimSpace(value)

	if result, err := strconv.ParseFloat(trimmedValue, 64); err == nil {
		return result
	}
	if result, err := strconv.Atoi(trimmedValue); err == nil {
		return result
	}
	if result, err := strconv.ParseBool(trimmedValue); err == nil {
		return result
	}
	if strings.HasPrefix(trimmedValue, "[") && strings.HasSuffix(trimmedValue, "]") {
		var result []interface{}
		trimmedValue = strings.TrimPrefix(trimmedValue, "[")
		trimmedValue = strings.TrimSuffix(trimmedValue, "]")
		for _, value := range SplitQuotedStringBy(trimmedValue, ',') {
			result = append(result, StringToInterface(value))
		}
		return result
	}
	if strings.HasPrefix(trimmedValue, "{") && strings.HasSuffix(trimmedValue, "}") {
		if result, err := DecodeJSONToStringKeyMap(trimmedValue); err == nil {
			return result
		}
	}
	if result := InterfaceToTime(trimmedValue); result != time.Unix(0, 0) {
		return result
	}

	return trimmedValue
}
