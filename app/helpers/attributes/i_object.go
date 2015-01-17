package attributes

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// Get returns object attribute value or nil
func (it *CustomAttributes) Get(attribute string) interface{} {
	return it.values[attribute]
}

// Set sets attribute value to object or returns error
func (it *CustomAttributes) Set(attribute string, value interface{}) error {
	if _, present := it.attributes[attribute]; present {
		it.values[attribute] = value
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "154b03ed-0e75-416d-890b-8775fcd74063", "attribute '"+attribute+"' invalid")
	}

	return nil
}

// GetAttributesInfo represents object as map[string]interface{}
func (it *CustomAttributes) GetAttributesInfo() []models.StructAttributeInfo {
	var info []models.StructAttributeInfo
	for _, attribute := range it.attributes {
		info = append(info, attribute)
	}
	return info
}

// FromHashMap represents object as map[string]interface{}
func (it *CustomAttributes) FromHashMap(input map[string]interface{}) error {
	it.values = input
	return nil
}

// ToHashMap fills object attributes from map[string]interface{}
func (it *CustomAttributes) ToHashMap() map[string]interface{} {
	return it.values
}
