package defaultproduct

import (
	"strconv"
)

func (dpm *DefaultProductModel) FromHashMap(input map[string]interface{}) error {

	if value, ok := input["_id"]; ok {
		if value, ok := value.(int64); ok {
			dpm.id = strconv.FormatInt(value, 10)
		}
	}
	if value, ok := input["sku"]; ok {
		if value, ok := value.(string); ok {
			dpm.Sku = value
		}
	}
	if value, ok := input["name"]; ok {
		if value, ok := value.(string); ok {
			dpm.Name = value
		}
	}

	dpm.CustomAttribute.FromHashMap(input)

	return nil
}

func (dpm *DefaultProductModel) ToHashMap() map[string]interface{} {
	result := dpm.CustomAttribute.ToHashMap()

	result["_id"] = dpm.id
	result["sku"] = dpm.Sku
	result["name"] = dpm.Name

	return result
}
