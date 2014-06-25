package address

func (it *DefaultVisitorAddress) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return err
		}
	}

	return nil
}

func (it *DefaultVisitorAddress) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["street"] = it.Street
	result["city"] = it.City
	result["state"] = it.State
	result["phone"] = it.Phone
	result["zip_code"] = it.ZipCode

	return result
}
