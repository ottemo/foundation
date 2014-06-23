package defaultproduct

import (
	"strings"

	"github.com/ottemo/foundation/models"
)

func (it *DefaultProductModel) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "sku":
		return it.Sku
	case "name":
		return it.Name
	default:
		return it.CustomAttributes.Get(attribute)
	}

	return nil
}

func (it *DefaultProductModel) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		it.id = value.(string)
	case "sku":
		it.Sku = value.(string)
	case "name":
		it.Name = value.(string)
	default:
		if err := it.CustomAttributes.Set(attribute, value); err != nil {
			return err
		}

	}

	return nil
}

func (it *DefaultProductModel) GetAttributesInfo() []models.AttributeInfo {
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

	dynamicInfo := it.CustomAttributes.GetAttributesInfo()

	return append(dynamicInfo, staticInfo...)
}
