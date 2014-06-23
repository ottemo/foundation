package default_address

import (
	"strings"

	"github.com/ottemo/foundation/models"
)

func (dva *DefaultVisitorAddress) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return dva.id
	case "street":
		return dva.Street
	case "city":
		return dva.City
	case "state":
		return dva.State
	case "phone":
		return dva.Phone
	case "zip", "zip_code":
		return dva.ZipCode
	}

	return nil
}

func (dva *DefaultVisitorAddress) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		dva.id = value.(string)
	case "street":
		dva.Street = value.(string)
	case "city":
		dva.City = value.(string)
	case "state":
		dva.State = value.(string)
	case "phone":
		dva.Phone = value.(string)
	case "zip", "zip_code":
		dva.ZipCode = value.(string)
	}
	return nil
}

func (dva *DefaultVisitorAddress) GetAttributesInfo() []models.AttributeInfo {
	info := []models.AttributeInfo{
		models.AttributeInfo{
			Model:      "VisitorAddress",
			Collection: "visitor_address",
			Attribute:  "_id",
			Type:       "text",
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.AttributeInfo{
			Model:      "VisitorAddress",
			Collection: "visitor_address",
			Attribute:  "street",
			Type:       "text",
			Label:      "Street",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.AttributeInfo{
			Model:      "VisitorAddress",
			Collection: "visitor_address",
			Attribute:  "city",
			Type:       "text",
			Label:      "City",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.AttributeInfo{
			Model:      "VisitorAddress",
			Collection: "visitor_address",
			Attribute:  "phone",
			Type:       "text",
			Label:      "Phone",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.AttributeInfo{
			Model:      "VisitorAddress",
			Collection: "visitor_address",
			Attribute:  "zip_code",
			Type:       "text",
			Label:      "Zip",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
	}

	return info
}
