package utils

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

// checks if value is blank (zero value)
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

// checks presence of non blank values for keys in map
//   - first arg must be map
//   - fallowing arguments are map keys you want to check
func KeysInMapAndNotBlank(mapObject interface{}, keys ...interface{}) bool {
	switch typedMapObject := mapObject.(type) {
	case map[string]interface{}:
		for _, key := range keys {
			if key2, ok := key.(string); ok {
				if mapValue, present := typedMapObject[key2]; present {
					if isBlank := CheckIsBlank(mapValue); isBlank {
						return false
					}
				} else {
					return false
				}
			}
		}
	case map[int]interface{}:
		for _, key := range keys {
			if key, ok := key.(int); ok {
				if mapValue, present := typedMapObject[key]; present {
					if isBlank := CheckIsBlank(mapValue); isBlank {
						return false
					}
				} else {
					return false
				}
			}
		}
	}

	return true
}

// returns map key value or nil if not found, will be returned first found key value
func GetFirstMapValue(mapObject interface{}, keys ...string) interface{} {
	switch typedMapObject := mapObject.(type) {
	case map[string]interface{}:
		for _, key := range keys {
			if value, present := typedMapObject[key]; present {
				return value
			}
		}
	case map[string]string:
		for _, key := range keys {
			if value, present := typedMapObject[key]; present {
				return value
			}
		}
	}

	return nil
}

// searches for item in array/slice, returns true if found
func IsInArray(searchValue interface{}, arrayObject interface{}) bool {
	switch typedObject := arrayObject.(type) {
	case []string:
		searchValue, ok := searchValue.(string)
		if !ok {
			return false
		}

		for _, value := range typedObject {
			if value == searchValue {

				return true
			}
		}

	case []interface{}:
		for _, value := range typedObject {
			if value == searchValue {
				return true
			}
		}
	}

	return false
}

// searches for presence of 1-st arg string option among provided options since 2-nd argument
func IsAmongStr(option string, searchOptions ...string) bool {
	for _, listOption := range searchOptions {
		if option == listOption {
			return true
		}
	}
	return false
}

// searches for a string in []string slice
func IsInListStr(searchItem string, searchList []string) bool {
	for _, listItem := range searchList {
		if listItem == searchItem {
			return true
		}
	}
	return false
}

// checks presence of string keys in map
func StrKeysInMap(mapObject interface{}, keys ...string) bool {
	switch typedMap := mapObject.(type) {
	case map[string]interface{}: //, map[string]string:

		for _, key := range keys {
			if _, present := typedMap[key]; !present {
				return false
			}
		}
	}

	return true
}

// returns trimmed []string array of [separators] delimited values
func Explode(value string, separators string) []string {
	result := make([]string, 0)

	splitResult := strings.Split(value, separators)
	for _, arrayValue := range splitResult {
		arrayValue = strings.TrimSpace(arrayValue)
		if arrayValue != "" {
			result = append(result, arrayValue)
		}
	}

	return result
}

// checks if value is MD5
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

// converts interface{} to string
func InterfaceToBool(value interface{}) bool {
	switch typedValue := value.(type) {
	case bool:
		return typedValue
	case string:
		typedValue = strings.TrimSpace(typedValue)
		if boolValue, err := strconv.ParseBool(typedValue); err == nil {
			return boolValue
		} else {
			typedValue = strings.ToLower(typedValue)
			if IsAmongStr(typedValue, "1", "ok", "on", "yes", "y", "x", "+") {
				return true
			}
			if intValue, err := strconv.Atoi(typedValue); err == nil && intValue > 0 {
				return true
			}

			return false
		}
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

// converts interface{} to map[string]interface{}
func InterfaceToArray(value interface{}) []interface{} {
	result := make([]interface{}, 0)

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
		jsonArray, err := DecodeJsonToArray(typedValue)
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

// converts interface{} to map[string]interface{}
func InterfaceToMap(value interface{}) map[string]interface{} {
	switch typedValue := value.(type) {
	case map[string]interface{}:
		return typedValue

	case string:
		value, err := DecodeJsonToStringKeyMap(value)
		if err == nil {
			return value
		}
	}

	return make(map[string]interface{})
}

// converts interface{} to string
func InterfaceToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch value := value.(type) {
	case bool:
		return strconv.FormatBool(value)
	case int:
		return strconv.Itoa(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case float64:
		return strconv.FormatFloat(value, 'f', 6, 64)
	case string:
		return value
	default:
		return EncodeToJsonString(value)
	}
}

// converts interface{} to integer
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

// converts interface{} to float64
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

// checks time for zero value
func IsZeroTime(value time.Time) bool {
	zeroTime := time.Unix(0, 0)
	return value == zeroTime
}

// converts interface{} to time.Time
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
		zoneFormats := []string{"", " MST", " -0700", "Z07:00"}

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

	return time.Unix(0, 0)
}

// convert string to integer
func StringToInteger(value string) (int, error) {
	return strconv.Atoi(value)
}

// convert string to float64
func StringToFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// rounds value to given precision (roundOn=0.5 usual cases)
func Round(val float64, roundOn float64, places int) float64 {

	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)

	var round float64
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	return round / pow
}

// function you should use to normalize price after calculations
// (in future that function should have config options setup)
func RoundPrice(price float64) float64 {
	return Round(price, 0.5, 2)
}

// splits string by character(s) unless it in quotes
func SplitQuotedStringBy(text string, separators ...rune) []string {

	lastQuote := rune(0)
	escapeFlag := false

	operator := func(currentChar rune) bool {

		isSeparatorChar := false
		for _, separator := range separators {
			if currentChar == separator {
				isSeparatorChar = true
				break
			}
		}

		switch {
		case currentChar == '\\':
			escapeFlag = true

		case !escapeFlag && lastQuote == currentChar:
			lastQuote = rune(0)
			return false

		case lastQuote == rune(0) && (currentChar == '"' || currentChar == '\'' || currentChar == '`'):
			lastQuote = currentChar
			return false

		case lastQuote == rune(0) && isSeparatorChar:
			return true
		}

		escapeFlag = false
		return false
	}

	return strings.FieldsFunc(text, operator)
}

// converts string coded value to type
func StringToType(value string, valueType string) (interface{}, error) {
	valueType = strings.ToLower(valueType)

	switch {
	case strings.HasPrefix(valueType, "[]"):
		value = strings.Trim(value, "[] \t\n")
		array := strings.Split(value, ",")
		arrayType := strings.TrimPrefix(valueType, "[]")

		result := make([]interface{}, 0)
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
		return DecodeJsonToStringKeyMap(value)
	}

	return nil, errors.New("unknown value type " + valueType)
}

// converts string to Interface{} which can be float, int, bool, string, or json as map[string]value
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
		result := make([]interface{}, 0)
		trimmedValue = strings.TrimPrefix(trimmedValue, "[")
		trimmedValue = strings.TrimSuffix(trimmedValue, "]")
		for _, value := range SplitQuotedStringBy(trimmedValue, ',') {
			result = append(result, StringToInterface(value))
		}
		return result
	}
	if strings.HasPrefix(trimmedValue, "{") && strings.HasSuffix(trimmedValue, "}") {
		if result, err := DecodeJsonToStringKeyMap(trimmedValue); err == nil {
			return result
		}
	}
	if result := InterfaceToTime(trimmedValue); result != time.Unix(0, 0) {
		return result
	}

	return trimmedValue
}
