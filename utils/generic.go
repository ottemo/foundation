package utils

import (
	"math"
	"strings"
)

// KeysInMapAndNotBlank checks presence of non blank values for keys in map
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

// GetFirstMapValue returns map key value or nil if not found, will be returned first found key value
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

// IsInArray searches for item in array/slice, returns true if found
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

// IsAmongStr searches for presence of 1-st arg string option among provided options since 2-nd argument
func IsAmongStr(option string, searchOptions ...string) bool {
	for _, listOption := range searchOptions {
		if option == listOption {
			return true
		}
	}
	return false
}

// IsInListStr searches for a string in []string slice
func IsInListStr(searchItem string, searchList []string) bool {
	for _, listItem := range searchList {
		if listItem == searchItem {
			return true
		}
	}
	return false
}

// StrKeysInMap checks presence of string keys in map
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

// Explode returns trimmed []string array of [separators] delimited values
func Explode(value string, separators string) []string {
	var result []string

	splitResult := strings.Split(value, separators)
	for _, arrayValue := range splitResult {
		arrayValue = strings.TrimSpace(arrayValue)
		if arrayValue != "" {
			result = append(result, arrayValue)
		}
	}

	return result
}

// Round rounds value to given precision (roundOn=0.5 usual cases)
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

// RoundPrice normalize price after calculations, so it rounds it to money precision
func RoundPrice(price float64) float64 {
	return Round(price, 0.5, 2)
}

// SplitQuotedStringBy splits string by character(s) unless it in quotes
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

// MatchMapAValuesToMapB compares key values of mapA to same key value of mapB, returns true
// if all keys in mapA present and matches keys in mapB
func MatchMapAValuesToMapB(mapA map[string]interface{}, mapB map[string]interface{}) bool {

	if mapA == nil || mapB == nil {
		if mapA == nil && mapB == nil {
			return true
		}
		return false
	}

	if len(mapA) == 0 {
		return true
	}

	for key, valueA := range mapA {
		if valueB, present := mapB[key]; present {
			switch valueA.(type) {
			case string, int, int32, int64, float32, float64, uint, uint32, uint64, complex64, complex128:
				return valueA == valueB

			case map[string]interface{}:
				typedValueA, okA := valueA.(map[string]interface{})
				typedValueB, okB := valueB.(map[string]interface{})

				if okA && okB {
					return MatchMapAValuesToMapB(typedValueA, typedValueB)
				}
				return false

			default:
				return valueA == valueB
			}
		} else {
			break
		}
	}

	return false
}

// EscapeRegexSpecials returns regular expression special characters escaped value
func EscapeRegexSpecials(value string) string {
	specials := []string{"\\", "-", "[", "]", "/", "{", "}", "(", ")", "*", "+", "?", ".", "^", "$", "|"}

	for _, special := range specials {
		value = strings.Replace(value, special, "\\"+special, -1)
	}

	return value
}
