package default_product

import ("strconv")

func (it *DefaultProductModel) FromHashMap(input map[string]interface{}) error {

	if value, ok := input["_id"]; ok {
		if value, ok := value.(int64); ok {
			it.id = strconv.FormatInt(value, 10)
		}
	}
	if value, ok := input["sku"]; ok {
		if value, ok := value.(string); ok {
			it.Sku = value
		}
	}
	if value, ok := input["name"]; ok {
		if value, ok := value.(string); ok {
			it.Name = value
		}
	}

	it.CustomAttributes.FromHashMap(input)

	return nil
}

func (it *DefaultProductModel) ToHashMap() map[string]interface{} {
	result := it.CustomAttributes.ToHashMap()

	result["_id"] = it.id
	result["sku"] = it.Sku
	result["name"] = it.Name

	return result
}
