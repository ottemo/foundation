package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"strings"
)

// Get returns an object attribute value or nil
func (it *DefaultProduct) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "sku":
		return it.Sku
	case "name":
		return it.Name
	case "short_description":
		return it.ShortDescription
	case "description":
		return it.Description
	case "default_image", "defaultimage":
		return it.DefaultImage
	case "price":
		return it.Price
	case "weight":
		return it.Weight
	case "options":
		return it.Options
	case "related_pids":
		return it.GetRelatedProductIds()
	}

	return it.CustomAttributes.Get(attribute)
}

// Set will apply the given attribute value to the product or return an error
func (it *DefaultProduct) Set(attribute string, value interface{}) error {
	lowerCaseAttribute := strings.ToLower(attribute)

	switch lowerCaseAttribute {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)
	case "sku":
		it.Sku = utils.InterfaceToString(value)
	case "name":
		it.Name = utils.InterfaceToString(value)
	case "short_description":
		it.ShortDescription = utils.InterfaceToString(value)
	case "description":
		it.Description = utils.InterfaceToString(value)
	case "default_image", "defaultimage":
		it.DefaultImage = utils.InterfaceToString(value)
	case "price":
		it.Price = utils.InterfaceToFloat64(value)
	case "weight":
		it.Weight = utils.InterfaceToFloat64(value)
	case "options":
		it.Options = utils.InterfaceToMap(value)
	case "related_pids":
		it.RelatedProductIds = make([]string, 0)

		switch typedValue := value.(type) {
		case []product.I_Product:
			for _, productItem := range typedValue {
				it.RelatedProductIds = append(it.RelatedProductIds, productItem.GetId())
			}

		case []interface{}:
			for _, listItem := range typedValue {
				if productID, ok := listItem.(string); ok && productID != "" && it.id == productID {
					// checking product existance
					productModel, err := product.LoadProductById(productID)
					if err != nil {
						return env.ErrorDispatch(err)
					}

					it.RelatedProductIds = append(it.RelatedProductIds, productModel.GetId())
				}
			}

		default:
			if value != nil {
				return env.ErrorNew("unsupported 'related_pids' attribute value")
			}
		}

	default:
		err := it.CustomAttributes.Set(attribute, value)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// FromHashMap will populate object attributes from map[string]interface{}
func (it *DefaultProduct) FromHashMap(input map[string]interface{}) error {
	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}
	return nil
}

// ToHashMap will return a map[string]interface{}
func (it *DefaultProduct) ToHashMap() map[string]interface{} {
	result := it.CustomAttributes.ToHashMap()

	result["_id"] = it.id
	result["sku"] = it.Sku
	result["name"] = it.Name

	result["short_description"] = it.ShortDescription
	result["description"] = it.Description

	result["default_image"] = it.DefaultImage

	result["price"] = it.Price
	result["weight"] = it.Weight
	result["options"] = it.Options

	result["related_pids"] = it.Get("related_pids")

	return result
}

// GetAttributesInfo will return the requested object attributes
func (it *DefaultProduct) GetAttributesInfo() []models.T_AttributeInfo {
	result := []models.T_AttributeInfo{
		models.T_AttributeInfo{
			Model:      product.MODEL_NAME_PRODUCT,
			Collection: COLLECTION_NAME_PRODUCT,
			Attribute:  "_id",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      product.MODEL_NAME_PRODUCT,
			Collection: COLLECTION_NAME_PRODUCT,
			Attribute:  "sku",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "SKU",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      product.MODEL_NAME_PRODUCT,
			Collection: COLLECTION_NAME_PRODUCT,
			Attribute:  "name",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      product.MODEL_NAME_PRODUCT,
			Collection: COLLECTION_NAME_PRODUCT,
			Attribute:  "short_description",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Short Description",
			Group:      "General",
			Editors:    "multiline_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      product.MODEL_NAME_PRODUCT,
			Collection: COLLECTION_NAME_PRODUCT,
			Attribute:  "description",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Description",
			Group:      "General",
			Editors:    "multiline_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      product.MODEL_NAME_PRODUCT,
			Collection: COLLECTION_NAME_PRODUCT,
			Attribute:  "default_image",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "DefaultImage",
			Group:      "Pictures",
			Editors:    "image_selector",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      product.MODEL_NAME_PRODUCT,
			Collection: COLLECTION_NAME_PRODUCT,
			Attribute:  "price",
			Type:       "numeric",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Price",
			Group:      "General",
			Editors:    "price",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      product.MODEL_NAME_PRODUCT,
			Collection: COLLECTION_NAME_PRODUCT,
			Attribute:  "weight",
			Type:       "numeric",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Weight",
			Group:      "General",
			Editors:    "numeric",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      product.MODEL_NAME_PRODUCT,
			Collection: COLLECTION_NAME_PRODUCT,
			Attribute:  "options",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Options",
			Group:      "Custom",
			Editors:    "product_options",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      product.MODEL_NAME_PRODUCT,
			Collection: COLLECTION_NAME_PRODUCT,
			Attribute:  "related_pids",
			Type:       "id",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Related Products",
			Group:      "General",
			Editors:    "product_selector",
			Options:    "",
			Default:    "",
		},
	}

	customAttributesInfo := it.CustomAttributes.GetAttributesInfo()
	for _, customAttribute := range customAttributesInfo {
		result = append(result, customAttribute)
	}

	return result
}
