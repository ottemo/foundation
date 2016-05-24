package utils

import (
	"math"
	"reflect"
	"regexp"
	"strings"
	"unicode"
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
func Round(value float64, round float64, precision int) float64 {

	negative := false
	if negative = math.Signbit(value); negative {
		value = -value
	}

	precisionPart := math.Pow(10, float64(precision))
	poweredValue := precisionPart * value
	_, roundingPart := math.Modf(poweredValue)

	var roundResult float64
	if roundingPart >= round {
		roundResult = math.Ceil(poweredValue)
	} else {
		roundResult = math.Floor(poweredValue)
	}

	if negative {
		roundResult = -roundResult
	}

	return roundResult / precisionPart
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

	for key, valueA := range mapA {
		if valueB, present := mapB[key]; present {
			switch valueA.(type) {
			case []interface{}:
				typedValueA, okA := valueA.([]interface{})
				typedValueB, okB := valueB.([]interface{})

				if okA && okB {
					for _, itemA := range typedValueA {
						found := false
						for _, itemB := range typedValueB {
							if itemA == itemB {
								found = true
								break
							}
						}
						if !found {
							return false
						}
					}
					return true
				}
				return false

			case map[string]interface{}:
				typedValueA, okA := valueA.(map[string]interface{})
				typedValueB, okB := valueB.(map[string]interface{})

				if okA && okB {
					return MatchMapAValuesToMapB(typedValueA, typedValueB)
				}
				return false

			default:
				if valueA != valueB {
					return false
				}
			}
		} else {
			return false
		}
	}

	return true
}

// EscapeRegexSpecials returns regular expression special characters escaped value
func EscapeRegexSpecials(value string) string {
	specials := []string{"\\", "-", "[", "]", "/", "{", "}", "(", ")", "*", "+", "?", ".", "^", "$", "|"}

	for _, special := range specials {
		value = strings.Replace(value, special, "\\"+special, -1)
	}

	return value
}

// ValidEmailAddress takes an email address as string compares it agasint a regular expression
// - returns true if email address is in a valid format
// - returns false if email address is not in a valid format
func ValidEmailAddress(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)
	return re.MatchString(email)
}

// Clone will create a replica for a given object
func Clone(subject interface{}) interface{} {
	result := subject

	subjectValue := reflect.ValueOf(subject)
	subjectKind := subjectValue.Kind()

	if subjectKind == reflect.Array || subjectKind == reflect.Slice {
		len := subjectValue.Len()
		newValue := reflect.MakeSlice(subjectValue.Type(), len, len)
		for idx := 0; idx < subjectValue.Len(); idx++ {
			value := Clone(subjectValue.Index(idx).Interface())
			newValue.Index(idx).Set(reflect.ValueOf(value))
		}
		result = newValue.Interface()
	} else if subjectKind == reflect.Map {
		newValue := reflect.MakeMap(subjectValue.Type())
		mapKeys := subjectValue.MapKeys()
		for _, key := range mapKeys {
			value := Clone(subjectValue.MapIndex(key).Interface())
			newValue.SetMapIndex(key, reflect.ValueOf(value))
		}
		result = newValue.Interface()
	}

	return result
}

// Convert string to snake_case format
func StrToSnakeCase(str string) string {

	str = strings.TrimSpace(str)

	// split camelCase words by space
	reg := regexp.MustCompile("(\\w+?)([A-Z]{1}[^A-Z][^\\W+])")
	str1 := reg.ReplaceAllString(str, "$1 $2")

	// remove all spacial symbols around "-"
	escapeSpaceInMinus := regexp.MustCompile("(-)\\W")
	str2 := escapeSpaceInMinus.ReplaceAllString(str1, "$1")

	// replace all spacial symbols except "-" on "_"
	nonSymbolReg := regexp.MustCompile("[^\\w^-]+")
	str3 := nonSymbolReg.ReplaceAllString(str2, "_")

	return strings.ToLower(str3)
}

// Convert string to snake case format
func StrToCamelCase(str string) string {

	re := regexp.MustCompile("_")
	str = strings.Title(re.ReplaceAllString(str, " "))

	str = strings.Replace(str, " ", "", -1)
	//first character to lower case
	s := []rune(str)
	s[0] = unicode.ToLower(s[0])
	return string(s)
}