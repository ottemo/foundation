package utils

import (
	"encoding/json"
	"errors"
)

// EncodeToJSONString encodes inputData to JSON string if it's possible
func EncodeToJSONString(inputData interface{}) string {
	result, _ := json.Marshal(inputData)
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
		err = errors.New("Unsupported json type for conversion to array")
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
		err = errors.New("Unable to create map, unsupported json type")
	}

	return result, err
}

// DecodeJSONToInterface decodes json string to interface{} if it's possible
func DecodeJSONToInterface(jsonData interface{}) (interface{}, error) {

	var result interface{}

	var err error

	switch value := jsonData.(type) {
	case string:
		err = json.Unmarshal([]byte(value), &result)
	case []byte:
		err = json.Unmarshal(value, &result)
	default:
		err = errors.New("Unable to parse json, unsupported json type")
	}

	return result, err
}
