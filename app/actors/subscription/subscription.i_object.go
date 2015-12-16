package subscription

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Get returns object attribute value or nil for the requested Subscription attribute
func (it *DefaultSubscription) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "visitor_id":
		return it.VisitorID
	case "cart_id":
		return it.CartID
	case "order_id":
		return it.OrderID
	case "email":
		return it.Email
	case "name", "full_name":
		return it.Name
	case "address":
		return it.Address
	case "status":
		return it.Status
	case "state":
		return it.State
	case "action":
		return it.Action
	case "period":
		return it.Period
	case "last_submit":
		return it.LastSubmit
	case "action_date":
		return it.ActionDate
	case "created_at":
		return it.CreatedAt
	case "updated_at":
		return it.UpdatedAt
	}

	return nil
}

// Set will set attribute value of the Subscription to object or return an error
func (it *DefaultSubscription) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)
	case "visitor_id":
		it.VisitorID = utils.InterfaceToString(value)
	case "cart_id":
		it.CartID = utils.InterfaceToString(value)
	case "order_id":
		it.OrderID = utils.InterfaceToString(value)
	case "email":
		it.Email = utils.InterfaceToString(value)
	case "name", "full_name":
		it.Name = utils.InterfaceToString(value)
	case "address":
		it.Address = utils.InterfaceToMap(value)
	case "status":
		it.Status = utils.InterfaceToString(value)
	case "state":
		it.State = utils.InterfaceToString(value)
	case "action":
		it.Action = utils.InterfaceToString(value)
	case "period":
		it.Period = utils.InterfaceToInt(value)
	case "last_submit":
		it.LastSubmit = utils.InterfaceToTime(value)
	case "action_date":
		it.ActionDate = utils.InterfaceToTime(value)
	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)
	case "updated_at":
		it.UpdatedAt = utils.InterfaceToTime(value)
	}

	return nil
}

// FromHashMap fills Subscription object attributes from a map[string]interface{}
func (it *DefaultSubscription) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap represents Subscription object as map[string]interface{}
func (it *DefaultSubscription) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["visitor_id"] = it.VisitorID
	result["cart_id"] = it.CartID
	result["order_id"] = it.OrderID

	result["email"] = it.Email
	result["name"] = it.Name

	result["status"] = it.Status
	result["state"] = it.State
	result["action"] = it.Action

	result["period"] = it.Period
	result["address"] = it.Address

	result["action_date"] = it.ActionDate
	result["last_submit"] = it.LastSubmit
	result["updated_at"] = it.UpdatedAt
	result["created_at"] = it.CreatedAt

	return result
}

// GetAttributesInfo returns the Subscription attributes information in an array
func (it *DefaultSubscription) GetAttributesInfo() []models.StructAttributeInfo {
	return make([]models.StructAttributeInfo, 0)
}
