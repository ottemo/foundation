package utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

// EncodeToJSONString encodes inputData to JSON string if it's possible
func EncodeToJSONString(inputData interface{}) string {
	if result, err := json.Marshal(inputData); err == nil {
		return string(result)
	}

	result, _ := json.Marshal(checkToJSON(inputData, 0))
	return string(result)
}

// DecodeJSONToArray decodes json string to []interface{} if it's possible
func DecodeJSONToArray(jsonData interface{}) ([]interface{}, error) {
	var result []interface{}

	var err error
	switch value := jsonData.(type) {
	case string:
		err = json.Unmarshal([]byte(value), &result)
	case []byte:
		err = json.Unmarshal(value, &result)
	default:
		err = errors.New("unsupported json data")
	}

	return result, err
}

// DecodeJSONToStringKeyMap decodes json string to map[string]interface{} if it's possible
func DecodeJSONToStringKeyMap(jsonData interface{}) (map[string]interface{}, error) {

	result := make(map[string]interface{})

	var err error

	switch value := jsonData.(type) {
	case string:
		err = json.Unmarshal([]byte(value), &result)
	case []byte:
		err = json.Unmarshal(value, &result)
	default:
		err = errors.New("unsupported json data")
	}

	return result, err
}

// checkToJSON internal function to convert data to JSON, some data may by not present after it
func checkToJSON(value interface{}, depth int) interface{} {
	if _, err := json.Marshal(value); err == nil {
		return value
	}

	// prevent from infinite loop of this function
	if depth >= 25 {
		return "*" + fmt.Sprint(value)
	}
	depth++

	// this switch should use this function in case we can't convert object to json string using native method
	var result interface{}
	switch typedValue := value.(type) {
	case map[string]interface{}:
		for key, partValue := range typedValue {
			typedValue[key] = checkToJSON(partValue, depth)
		}
		result = typedValue
		break

	case map[interface{}]interface{}:
		convertedMap := make(map[string]interface{})
		for key, partValue := range typedValue {
			convertedMap[InterfaceToString(key)] = checkToJSON(partValue, depth)
		}
		result = convertedMap
		break

	case []interface{}:
		for key, partValue := range typedValue {
			typedValue[key] = checkToJSON(partValue, depth)
		}
		result = typedValue
		break

	default:
		return fmt.Sprint(value)
	}

	if _, err := json.Marshal(result); err != nil {
		return fmt.Sprint(value)
	}
	return result
}
