package utils

import (
	"encoding/json"
	"errors"
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
func checkToJSON(value interface{}, count int) interface{} {
	if _, err := json.Marshal(value); err == nil {
		return value
	}
	// limiting for execution from the same function
	count++
	if count >= 25 {
		return InterfaceToString(value)
	}

	// this switch should use this function in case we can't convert object to json string using native method
	var result interface{}
	switch typedValue := value.(type) {
	case map[string]interface{}:
		for key, partValue := range typedValue {
			typedValue[key] = checkToJSON(partValue, count)
		}
		result = typedValue
		break

	case map[interface{}]interface{}:
		convertedMap := make(map[string]interface{})
		for key, partValue := range typedValue {
			convertedMap[InterfaceToString(key)] = checkToJSON(partValue, count)
		}
		result = convertedMap
		break

	case []interface{}:
		for key, partValue := range typedValue {
			typedValue[key] = checkToJSON(partValue, count)
		}
		result = typedValue
		break

	default:
		return InterfaceToString(value)
	}

	_, err := json.Marshal(result)
	if err == nil {
		return result
	}
	return InterfaceToString(value)
}
