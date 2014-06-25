package defaultVisitor

func (it *DefaultVisitor) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return err
		}
	}

	return nil
}

func (it *DefaultVisitor) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["email"] = it.Email
	result["first_name"] = it.FirstName
	result["last_name"] = it.LastName

	result["billing_address"] = nil
	result["shipping_address"] = nil

	//result["billing_address_id"] = it.BillingAddressId
	//result["shipping_address_id"] = it.ShippingAddressId

	if it.BillingAddress != nil {
		result["billing_address"] = it.BillingAddress.ToHashMap()
	}

	if it.ShippingAddress != nil {
		result["shipping_address"] = it.ShippingAddress.ToHashMap()
	}

	return result
}
