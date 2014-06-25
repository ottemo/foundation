package defaultVisitor

import (
	"strings"

	"github.com/ottemo/foundation/models"

	"errors"
)

func (it *DefaultVisitor) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
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
	}

	return nil
}

func (it *DefaultVisitor) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.id = value.(string)
	case "email", "e_mail", "e-mail":
		it.Email = value.(string)
	case "fname", "first_name":
		it.FirstName = value.(string)
	case "lname", "last_name":
		it.LastName = value.(string)

	// only address id coming - trying to get it from DB
	case "billing_address_id", "shipping_address_id":
		address := it.getVisitorAddressById(value.(string))
		if address != nil && address.GetId() != "" {

			if attribute == "billing_address_id" {
				it.BillingAddress = address
			} else {
				it.ShippingAddress = address
			}

		} else {
			return errors.New("wrong address id")
		}

	// address with details coming
	case "billing_address", "shipping_address":
		switch value := value.(type) {

		// we have already have structure
		case IVisitorAddress:
			if attribute == "billing_address" {
				it.BillingAddress = value
			} else {
				it.ShippingAddress = value
			}

		// we have sub-map, supposedly IVisitorAddress capable
		case map[string]interface{}:
			model, err := models.GetModel("VisitorAddress")
			if err != nil {
				return err
			}

			if address, ok := model.(IVisitorAddress); ok {
				err := address.FromHashMap(value)
				if err != nil {
					return err
				}

				if attribute == "billing_address" {
					it.BillingAddress = address
				} else {
					it.ShippingAddress = address
				}
			} else {
				errors.New("unsupported visitor addres model " + model.GetImplementationName())
			}

		default:
			return errors.New("unsupported 'billing_address' value")
		}
	}
	return nil
}

func (it *DefaultVisitor) GetAttributesInfo() []models.T_AttributeInfo {

	info := []models.T_AttributeInfo{
		models.T_AttributeInfo{
			Model:      "Visitor",
			Collection: "Visitor",
			Attribute:  "_id",
			Type:       "text",
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Visitor",
			Collection: "Visitor",
			Attribute:  "email",
			Type:       "text",
			Label:      "E-mail",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Visitor",
			Collection: "Visitor",
			Attribute:  "first_name",
			Type:       "text",
			Label:      "First Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Visitor",
			Collection: "Visitor",
			Attribute:  "billing_address",
			Type:       "text",
			Label:      "Billing Address",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model:VisitorAddress",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Visitor",
			Collection: "Visitor",
			Attribute:  "shipping_address",
			Type:       "text",
			Label:      "Shipping Address",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model:VisitorAddress",
			Default:    "",
		},
	}

	return info
}
