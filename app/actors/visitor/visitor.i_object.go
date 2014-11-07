package visitor

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Get returns the requested default attribute value for the Visitor or nil
func (it *DefaultVisitor) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "email":
		return it.Email
	case "fname", "first_name":
		return it.FirstName
	case "lname", "last_name":
		return it.LastName
	case "billing_address_id":
		return it.BillingAddress.GetId()
	case "shipping_address_id":
		return it.ShippingAddress.GetId()
	case "billing_address":
		return it.BillingAddress
	case "shipping_address":
		return it.ShippingAddress
	case "validate":
		return it.ValidateKey
	case "facebook_id":
		return it.FacebookId
	case "google_id":
		return it.GoogleId
	case "birthday":
		return it.Birthday
	case "is_admin":
		return it.Admin
	case "created_at":
		return it.IsAdmin
	}

	return it.CustomAttributes.Get(attribute)
}

// Set takes a string and will set the provided attribute value on the Visitor or returns error
func (it *DefaultVisitor) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)
	case "email", "e_mail", "e-mail":
		it.Email = utils.InterfaceToString(value)
	case "fname", "first_name":
		it.FirstName = utils.InterfaceToString(value)
	case "lname", "last_name":
		it.LastName = utils.InterfaceToString(value)
	case "password", "passwd":
		it.SetPassword(utils.InterfaceToString(value))
	case "validate":
		it.ValidateKey = utils.InterfaceToString(value)
	case "facebook_id":
		it.FacebookId = utils.InterfaceToString(value)
	case "google_id":
		it.GoogleId = utils.InterfaceToString(value)
	case "birthday":
		it.Birthday = utils.InterfaceToTime(value)
	case "is_admin":
		it.Admin = utils.InterfaceToBool(value)
	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)

	// only address id was specified - trying to load it
	case "billing_address_id", "shipping_address_id":
		value := utils.InterfaceToString(value)

		var address visitor.I_VisitorAddress
		var err error

		if value != "" {
			address, err = visitor.LoadVisitorAddressById(value)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}

		if address == nil || address.GetId() != "" {

			if attribute == "billing_address_id" {
				it.BillingAddress = address
			} else {
				it.ShippingAddress = address
			}

		}

	// address detailed information was specified
	case "billing_address", "shipping_address":
		switch typedValue := value.(type) {

		// we have already have structure
		case visitor.I_VisitorAddress:
			if attribute == "billing_address" {
				it.BillingAddress = typedValue
			} else {
				it.ShippingAddress = typedValue
			}

		// we have sub-map, supposedly I_VisitorAddress capable
		case map[string]interface{}:
			var addressModel visitor.I_VisitorAddress
			var err error

			if len(typedValue) != 0 {
				addressModel, err = visitor.GetVisitorAddressModel()
				if err != nil {
					return env.ErrorDispatch(err)
				}

				err = addressModel.FromHashMap(typedValue)
				if err != nil {
					return env.ErrorDispatch(err)
				}
			}

			if attribute == "billing_address" {
				it.BillingAddress = addressModel
			} else {
				it.ShippingAddress = addressModel
			}
		default:
			return env.ErrorNew("unsupported billing or shipping address value")
		}

	default:
		err := it.CustomAttributes.Set(attribute, value)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// FromHashMap will populate the Visitor object attributes from the given map[string]interface{}
func (it *DefaultVisitor) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap will return the Visitor represented as a map[string]interface{}
func (it *DefaultVisitor) ToHashMap() map[string]interface{} {

	result := it.CustomAttributes.ToHashMap()

	result["_id"] = it.id

	result["email"] = it.Email
	result["first_name"] = it.FirstName
	result["last_name"] = it.LastName

	result["is_admin"] = it.Admin
	result["birthday"] = it.Birthday
	result["created_at"] = it.CreatedAt

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

// GetAttributesInfo will return the current Visitor object attributes and related information
func (it *DefaultVisitor) GetAttributesInfo() []models.T_AttributeInfo {

	result := []models.T_AttributeInfo{
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "_id",
			Type:       "id",
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "email",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "E-mail",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "first_name",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "First Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "last_name",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Last Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "password",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Password",
			Group:      "Password",
			Editors:    "password",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "billing_address_id",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Billing Address",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model=VisitorAddress",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "shipping_address_id",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Shipping Address",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model:VisitorAddress",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "birthday",
			Type:       "datetime",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Birthday",
			Group:      "General",
			Editors:    "datetime",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "created_at",
			Type:       "datetime",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Created at",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR,
			Collection: COLLECTION_NAME_VISITOR,
			Attribute:  "is_admin",
			Type:       "bool",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Is admin",
			Group:      "General",
			Editors:    "boolean",
			Options:    "",
			Default:    "false",
		},
	}

	customAttributesInfo := it.CustomAttributes.GetAttributesInfo()
	for _, customAttribute := range customAttributesInfo {
		result = append(result, customAttribute)
	}

	return result
}
