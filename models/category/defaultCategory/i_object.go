package defaultCategory

import (
	"errors"
	"strings"

	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/product"
)

func (it *DefaultCategory) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "name":
		return it.Name
	case "products":
		result := make([]map[string]interface{}, len(it.Products))
		for _, categoryProduct := range it.Products {
			result = append(result, categoryProduct.ToHashMap())
		}
		return result
	}

	return nil
}

func (it *DefaultCategory) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.id = value.(string)

	case "name":
		it.Name = value.(string)

	case "parent":
		value, ok := value.(ICategory)
		if !ok {
			errors.New("unsupported 'parent' value")
		}

		it.Parent = value

	case "products":
		switch value := value.(type) {

		// we have prepared []IProduct
		case []product.IProduct:
			it.Products = value

		// we have sub-maps array, supposedly []IProduct capable
		case []map[string]interface{}:
			for _, value := range value {
				model, err := models.GetModel("Product")
				if err != nil {
					return err
				}

				if categoryProduct, ok := model.(product.IProduct); ok {
					err := categoryProduct.FromHashMap(value)
					if err != nil {
						return err
					}

					it.Products = append(it.Products, categoryProduct)
				} else {
					errors.New("unsupported product model " + model.GetImplementationName())
				}
			}
		default:
			return errors.New("unsupported 'products' value")
		}
	}
	return nil
}

func (it *DefaultCategory) GetAttributesInfo() []models.T_AttributeInfo {

	info := []models.T_AttributeInfo{
		models.T_AttributeInfo{
			Model:      "Category",
			Collection: "Category",
			Attribute:  "_id",
			Type:       "id",
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Category",
			Collection: "Category",
			Attribute:  "name",
			Type:       "text",
			Label:      "Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Category",
			Collection: "Category",
			Attribute:  "parent",
			Type:       "id",
			Label:      "Parent",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: category",
			Default:    "",
		},
		models.T_AttributeInfo{
			Model:      "Category",
			Collection: "Category",
			Attribute:  "products",
			Type:       "id",
			Label:      "Products",
			Group:      "General",
			Editors:    "array_model_selector",
			Options:    "model: product",
			Default:    "",
		},
	}

	return info
}
