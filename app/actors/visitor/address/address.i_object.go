package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

// Get will return the requested attribute when provided a string representation of the attribute
func (it *DefaultVisitorAddress) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "visitor_id", "visitorID":
		return it.visitorID
	case "fname", "first_name":
		return it.FirstName
	case "lname", "last_name":
		return it.LastName
	case "address_line1":
		return it.AddressLine1
	case "address_line2":
		return it.AddressLine2
	case "company":
		return it.Company
	case "country":
		return it.Country
	case "city":
		return it.City
	case "state":
		return it.State
	case "phone":
		return it.Phone
	case "zip", "zip_code":
		return it.ZipCode
	}

	return nil
}

// Set will set a Visitor Address attribute and requiring a name and a value
func (it *DefaultVisitorAddress) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)

	case "visitor_id", "visitorID":
		it.visitorID = utils.InterfaceToString(value)

	case "fname", "first_name":
		it.FirstName = utils.InterfaceToString(value)

	case "lname", "last_name":
		it.LastName = utils.InterfaceToString(value)

	case "line1", "address_line1":
		it.AddressLine1 = utils.InterfaceToString(value)

	case "line2", "address_line2":
		it.AddressLine2 = utils.InterfaceToString(value)

	case "company":
		it.Company = utils.InterfaceToString(value)

	case "country":
		it.Country = utils.InterfaceToString(value)

	case "city":
		it.City = utils.InterfaceToString(value)

	case "state":
		it.State = utils.InterfaceToString(value)

	case "phone":
		it.Phone = utils.InterfaceToString(value)

	case "zip", "zip_code":
		it.ZipCode = utils.InterfaceToString(value)
	}
	return nil
}

// FromHashMap will take a map[string]interface and apply the attribute values to the Visitor Address
func (it *DefaultVisitorAddress) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.LogError(err)
		}
	}

	return nil
}

// ToHashMap will return a set of Visitor Address attributes in a map[string]interface
func (it *DefaultVisitorAddress) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["visitor_id"] = it.visitorID

	result["first_name"] = it.FirstName
	result["last_name"] = it.LastName

	result["company"] = it.Company

	result["address_line1"] = it.AddressLine1
	result["address_line2"] = it.AddressLine2

	result["country"] = it.Country
	result["city"] = it.City
	result["state"] = it.State

	result["phone"] = it.Phone
	result["zip_code"] = it.ZipCode

	return result
}

// GetAttributesInfo will return a set of Vistor Address attributes in []models.StructAttributeInfo
func (it *DefaultVisitorAddress) GetAttributesInfo() []models.StructAttributeInfo {
	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "_id",
			Type:       db.ConstTypeID,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "visitor_id",
			Type:       db.ConstTypeID,
			Label:      "Visitor ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "address_line1",
			Type:       db.ConstTypeVarchar,
			Label:      "Address Line 1",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "address_line2",
			Type:       db.ConstTypeVarchar,
			Label:      "Address Line 2",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "first_name",
			Type:       db.ConstTypeVarchar,
			Label:      "First Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "last_name",
			Type:       db.ConstTypeVarchar,
			Label:      "Last Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "company",
			Type:       db.ConstTypeVarchar,
			Label:      "Company",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "country",
			Type:       db.ConstTypeVarchar,
			Label:      "Country",
			Group:      "General",
			Editors:    "select",
			Options:    utils.EncodeToJSONString(models.ConstCountriesList),
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "city",
			Type:       db.ConstTypeVarchar,
			Label:      "City",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "phone",
			Type:       db.ConstTypeVarchar,
			Label:      "Phone",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: `regexp(/^(\+?\d{1,3})?[- ]?\(?(\d{3})\)?[- ]?((?:\d{3})-?(?:\d{2})-?(?:\d{2}))$/)`,
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "state",
			Type:       db.ConstTypeVarchar,
			Label:      "State",
			Group:      "General",
			Editors:    "select",
			Options:    utils.EncodeToJSONString(models.ConstStatesList),
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      visitor.ConstModelNameVisitorAddress,
			Collection: ConstCollectionNameVisitorAddress,
			Attribute:  "zip_code",
			Type:       db.ConstTypeVarchar,
			Label:      "Zip",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: `regexp(/^\d{5}(?:[-\s]\d{4})?$/)`,
		},
	}

	return info
}
