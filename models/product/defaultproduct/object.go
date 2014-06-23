package defaultproduct

import (
	"strings"

	"github.com/ottemo/foundation/models"
)

func (dpm *DefaultProductModel) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return dpm.id
	case "sku":
		return dpm.Sku
	case "name":
		return dpm.Name
	default:
		return dpm.CustomAttributes.Get(attribute)
	}

	return nil
}

func (dpm *DefaultProductModel) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		dpm.id = value.(string)
	case "sku":
		dpm.Sku = value.(string)
	case "name":
		dpm.Name = value.(string)
	default:
		if err := dpm.CustomAttributes.Set(attribute, value); err != nil {
			return err
		}

	}

	return nil
}

func (dpm *DefaultProductModel) GetAttributesInfo() []models.AttributeInfo {
	staticInfo := []models.AttributeInfo{
		models.AttributeInfo{
			Model:      "Product",
			Collection: "product",
			Attribute:  "_id",
			Type:       "text",
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.AttributeInfo{
			Model:      "Product",
			Collection: "product",
			Attribute:  "sku",
			Type:       "text",
			Label:      "SKU",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.AttributeInfo{
			Model:      "Product",
			Collection: "product",
			Attribute:  "Name",
			Type:       "text",
			Label:      "Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
	}

	dynamicInfo := dpm.CustomAttributes.GetAttributesInfo()

	return append(dynamicInfo, staticInfo...)
}
