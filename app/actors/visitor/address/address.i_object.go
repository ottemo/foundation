package address

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/visitor"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
	"strings"
)

func (it *DefaultVisitorAddress) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "visitor_id", "visitorId":
		return it.visitor_id
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

func (it *DefaultVisitorAddress) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)

	case "visitor_id", "visitorId":
		it.visitor_id = utils.InterfaceToString(value)

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

func (it *DefaultVisitorAddress) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

func (it *DefaultVisitorAddress) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["visitor_id"] = it.visitor_id

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

func (it *DefaultVisitorAddress) GetAttributesInfo() []models.T_AttributeInfo {
	info := []models.T_AttributeInfo{
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "_id",
			Type:       "id",
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "visitor_id",
			Type:       "id",
			Label:      "Visitor ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "address_line1",
			Type:       "varchar(255)",
			Label:      "Address Line 1",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "address_line2",
			Type:       "varchar(255)",
			Label:      "Address Line 2",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "first_name",
			Type:       "varchar(100)",
			Label:      "First Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "last_name",
			Type:       "varchar(100)",
			Label:      "Last Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "company",
			Type:       "varchar(100)",
			Label:      "Company",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "country",
			Type:       "varchar(50)",
			Label:      "Country",
			Group:      "General",
			Editors:    "select",
			Options:    utils.EncodeToJsonString(models.COUNTRIES_LIST),
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "city",
			Type:       "varchar(100)",
			Label:      "City",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "phone",
			Type:       "varchar(100)",
			Label:      "Phone",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "state",
			Type:       "varchar(2)",
			Label:      "State",
			Group:      "General",
			Editors:    "select",
			Options:    utils.EncodeToJsonString(models.STATES_LIST),
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      visitor.MODEL_NAME_VISITOR_ADDRESS,
			Collection: COLLECTION_NAME_VISITOR_ADDRESS,
			Attribute:  "zip_code",
			Type:       "varchar(10)",
			Label:      "Zip",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
	}

	return info
}
