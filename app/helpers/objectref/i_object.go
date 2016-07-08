package objectref

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// Get returns attribute value for current object or nil if no such attribute
func (it *DBObjectRef) Get(attribute string) interface{} {
	var result interface{}

	if it.currData != nil {
		result, _ = it.currData[attribute]
	}

	return result
}

// Set sets attribute value for current object
func (it *DBObjectRef) Set(attribute string, value interface{}) error {
	if it.currData == nil {
		it.currData = make(map[string]interface{})
	}

	it.currData[attribute] = value

	return nil
}

// FromHashMap fills attributes values based on provided map
func (it *DBObjectRef) FromHashMap(input map[string]interface{}) error {

	if it.currData == nil {
		it.currData = make(map[string]interface{})
	}

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap returns attribute values that current object holds
func (it *DBObjectRef) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	if it.currData != nil {
		for attribute, value := range it.currData {
			result[attribute] = value
		}
	}

	return result
}

// GetAttributesInfo returns stub information about current object attributes
//   - if you using this helper you should rewrite this function in your class
func (it *DBObjectRef) GetAttributesInfo() []models.StructAttributeInfo {

	result := []models.StructAttributeInfo{}

	if it.currData != nil {
		for attribute := range it.currData {
			result = append(result,
				models.StructAttributeInfo{
					Model:      "",
					Collection: "",
					Attribute:  attribute,
					Type:       "",
					IsRequired: false,
					IsStatic:   true,
					Label:      attribute,
					Group:      "General",
					Editors:    "not_editable",
					Options:    "",
					Default:    "",
				})
		}
	}

	return result
}
