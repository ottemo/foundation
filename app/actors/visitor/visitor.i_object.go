package visitor

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Get returns object attribute value or nil for the requested Visitor attribute
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
		if it.BillingAddress != nil {
			return it.BillingAddress.GetID()
		}
		return nil
	case "shipping_address_id":
		if it.ShippingAddress != nil {
			return it.ShippingAddress.GetID()
		}
		return nil
	case "billing_address":
		return it.BillingAddress
	case "shipping_address":
		return it.ShippingAddress
	case "validate":
		return it.ValidateKey
	case "facebook_id":
		return it.FacebookID
	case "google_id":
		return it.GoogleID
	case "is_admin":
		return it.Admin
	case "created_at":
		return it.IsAdmin
	}

	return it.CustomAttributes.Get(attribute)
}

// Set will set attribute value of the Visitor to object or return an error
func (it *DefaultVisitor) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)
	case "email", "e_mail", "e-mail":
		it.Email = strings.ToLower(utils.InterfaceToString(value))
	case "fname", "first_name":
		it.FirstName = utils.InterfaceToString(value)
	case "lname", "last_name":
		it.LastName = utils.InterfaceToString(value)
	case "password", "passwd":
		it.SetPassword(utils.InterfaceToString(value))
	case "validate":
		it.ValidateKey = utils.InterfaceToString(value)
	case "facebook_id":
		it.FacebookID = utils.InterfaceToString(value)
	case "google_id":
		it.GoogleID = utils.InterfaceToString(value)
	case "is_admin":
		it.Admin = utils.InterfaceToBool(value)
	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)

	// only address id was specified - trying to load it
	case "billing_address_id", "shipping_address_id":
		value := utils.InterfaceToString(value)

		var address visitor.InterfaceVisitorAddress
		var err error

		if value != "" {
			address, err = visitor.LoadVisitorAddressByID(value)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}

		if address == nil || address.GetID() != "" {

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
		case visitor.InterfaceVisitorAddress:
			if attribute == "billing_address" {
				it.BillingAddress = typedValue
			} else {
				it.ShippingAddress = typedValue
			}

		// we have sub-map, supposedly InterfaceVisitorAddress capable
		case map[string]interface{}:
			var addressModel visitor.InterfaceVisitorAddress
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
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "efa9cd9c-2d9a-4637-ac59-b4856d2e623e", "unsupported billing or shipping address value")
		}

	default:
		err := it.CustomAttributes.Set(attribute, value)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// FromHashMap fills Visitor object attributes from a map[string]interface{}
func (it *DefaultVisitor) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap represents Visitor object as map[string]interface{}
func (it *DefaultVisitor) ToHashMap() map[string]interface{} {

	result := it.CustomAttributes.ToHashMap()

	result["_id"] = it.id

	result["email"] = it.Email
	result["first_name"] = it.FirstName
	result["last_name"] = it.LastName

	result["is_admin"] = it.Admin
	result["created_at"] = it.CreatedAt

	result["billing_address"] = nil
	result["shipping_address"] = nil

	//result["billing_address_id"] = it.BillingAddressID
	//result["shipping_address_id"] = it.ShippingAddressID

	if it.BillingAddress != nil {
		result["billing_address"] = it.BillingAddress.ToHashMap()
	}

	if it.ShippingAddress != nil {
		result["shipping_address"] = it.ShippingAddress.ToHashMap()
	}

	return result
}

// GetAttributesInfo returns the Visitor attributes information in an array
func (it *DefaultVisitor) GetAttributesInfo() []models.StructAttributeInfo {

	result := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitor,
			Collection: ConstCollectionNameVisitor,
			Attribute:  "_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitor,
			Collection: ConstCollectionNameVisitor,
			Attribute:  "email",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "E-mail",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: "email",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitor,
			Collection: ConstCollectionNameVisitor,
			Attribute:  "first_name",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "First Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitor,
			Collection: ConstCollectionNameVisitor,
			Attribute:  "last_name",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Last Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitor,
			Collection: ConstCollectionNameVisitor,
			Attribute:  "password",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Password",
			Group:      "Password",
			Editors:    "password",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitor,
			Collection: ConstCollectionNameVisitor,
			Attribute:  "billing_address_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Billing Address",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model=VisitorAddress",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitor,
			Collection: ConstCollectionNameVisitor,
			Attribute:  "shipping_address_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Shipping Address",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model:VisitorAddress",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitor,
			Collection: ConstCollectionNameVisitor,
			Attribute:  "created_at",
			Type:       db.ConstTypeDatetime,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Created at",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitor,
			Collection: ConstCollectionNameVisitor,
			Attribute:  "is_admin",
			Type:       db.ConstTypeBoolean,
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
