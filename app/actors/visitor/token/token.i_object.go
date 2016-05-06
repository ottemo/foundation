package token

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Get will return the requested attribute when provided a string representation of the attribute
func (it *DefaultVisitorCard) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "visitor_id":
		return it.visitorID
	case "token_id":
		return it.tokenID
	case "customer_id":
		return it.customerID
	case "holder":
		return it.Holder
	case "payment":
		return it.Payment
	case "type":
		return it.Type
	case "number":
		return it.Number
	case "expiration_date":
		return it.ExpirationDate
	case "token_updated":
		return it.TokenUpdated
	case "created_at":
		return it.CreatedAt
	}

	return nil
}

// Set will set a Visitor Token attribute and requiring a name and a value
func (it *DefaultVisitorCard) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)

	case "visitor_id":
		it.visitorID = utils.InterfaceToString(value)

	case "token_id":
		it.tokenID = utils.InterfaceToString(value)

	case "customer_id":
		it.customerID = utils.InterfaceToString(value)

	case "holder":
		it.Holder = utils.InterfaceToString(value)

	case "payment":
		it.Payment = utils.InterfaceToString(value)

	case "type":
		it.Type = utils.InterfaceToString(value)

	case "number":
		it.Number = utils.InterfaceToString(value)

	case "expiration_date":
		it.ExpirationDate = utils.InterfaceToString(value)

	case "expire_year":
		it.ExpirationYear = utils.InterfaceToInt(value)

	case "expire_month":
		it.ExpirationMonth = utils.InterfaceToInt(value)

	case "token_updated":
		it.TokenUpdated = utils.InterfaceToTime(value)

	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)
	}
	return nil
}

// FromHashMap will take a map[string]interface and apply the attribute values to the Visitor Token
func (it *DefaultVisitorCard) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap will return a set of Visitor Token attributes in a map[string]interface
func (it *DefaultVisitorCard) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id
	result["visitor_id"] = it.visitorID

	result["holder"] = it.Holder
	result["payment"] = it.Payment
	result["type"] = it.Type
	result["number"] = it.Number
	result["expiration_date"] = it.ExpirationDate
	result["created_at"] = it.CreatedAt
	result["token_updated"] = it.TokenUpdated

	result["token_id"] = it.tokenID
	result["customer_id"] = it.customerID

	return result
}

// GetAttributesInfo will return a set of Vistor Token attributes in []models.StructAttributeInfo
func (it *DefaultVisitorCard) GetAttributesInfo() []models.StructAttributeInfo {
	info := []models.StructAttributeInfo{}

	return info
}
