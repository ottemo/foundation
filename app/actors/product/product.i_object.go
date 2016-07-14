package product

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Get returns an object attribute value or nil
func (it *DefaultProduct) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "enable", "enabled":
		return it.Enabled
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
	case "qty":
		return it.GetQty()
	case "options":
		return it.GetOptions()
	case "inventory":
		return it.GetInventory()
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
	case "enable", "enabled":
		it.Enabled = utils.InterfaceToBool(value)
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
	case "qty":
		it.Qty = utils.InterfaceToInt(value)
	case "options":
		it.Options = utils.InterfaceToMap(value)
	case "inventory":
		inventory := utils.InterfaceToArray(value)
		for _, options := range inventory {
			it.Inventory = append(it.Inventory, utils.InterfaceToMap(options))
		}
	case "related_pids":
		it.RelatedProductIds = make([]string, 0)

		switch typedValue := value.(type) {
		case []product.InterfaceProduct:
			for _, productItem := range typedValue {
				it.RelatedProductIds = append(it.RelatedProductIds, productItem.GetID())
			}

		case []interface{}:

			for _, listItem := range typedValue {
				var relatedProductIDs []string

				currentProductID := it.GetID()
				if relatedProductID, ok := listItem.(string); ok &&
					relatedProductID != "" &&
					currentProductID != relatedProductID {

					relatedProductIDs = append(relatedProductIDs, relatedProductID)
				}

				// checking related products existence
				dbCollection, err := db.GetCollection(ConstCollectionNameProduct)
				if err != nil {
					return env.ErrorDispatch(err)
				}
				err = dbCollection.AddFilter("_id", "in", relatedProductIDs)
				if err != nil {
					return env.ErrorDispatch(err)
				}
				err = dbCollection.SetResultColumns("_id")
				if err != nil {
					return env.ErrorDispatch(err)
				}
				records, err := dbCollection.Load()
				if err != nil {
					return env.ErrorDispatch(err)
				}

				// adding only exist products to model
				for _, record := range records {
					productID := utils.InterfaceToString(record["_id"])
					it.RelatedProductIds = append(it.RelatedProductIds, productID)
				}
			}

		default:
			if value != nil {
				return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3c402ecc-7c7d-49ab-879e-16af5f4661ed", "unsupported 'related_pids' attribute value")
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
			env.ErrorDispatch(err)
		}
	}
	return nil
}

// ToHashMap returns a map[string]interface{}
func (it *DefaultProduct) ToHashMap() map[string]interface{} {
	result := it.CustomAttributes.ToHashMap()

	result["_id"] = it.id

	result["enabled"] = it.Enabled

	result["sku"] = it.Sku
	result["name"] = it.Name

	result["short_description"] = it.ShortDescription
	result["description"] = it.Description

	result["default_image"] = it.DefaultImage

	result["price"] = it.Price
	result["weight"] = it.Weight

	result["options"] = it.GetOptions()

	result["related_pids"] = it.Get("related_pids")

	if product.GetRegisteredStock() != nil {
		result["qty"] = it.Get("qty")
	}

	return result
}

// GetAttributesInfo returns the requested object attributes
func (it *DefaultProduct) GetAttributesInfo() []models.StructAttributeInfo {
	result := []models.StructAttributeInfo{
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "enabled",
			Type:       db.ConstTypeBoolean,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Enabled",
			Group:      "General",
			Editors:    "boolean",
			Options:    "",
			Default:    "",
		},
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "sku",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "SKU",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: "sku",
		},
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "name",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "short_description",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Short Description",
			Group:      "General",
			Editors:    "multiline_text",
			Options:    "",
			Default:    "",
		},
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "description",
			Type:       db.ConstTypeText,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Description",
			Group:      "General",
			Editors:    "multiline_text",
			Options:    "",
			Default:    "",
		},
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "default_image",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "DefaultImage",
			Group:      "General",
			Editors:    "image_selector",
			Options:    "",
			Default:    "",
		},
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "price",
			Type:       db.ConstTypeMoney,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Price",
			Group:      "General",
			Editors:    "price",
			Options:    "",
			Default:    "",
			Validators: "price",
		},
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "weight",
			Type:       db.ConstTypeDecimal,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Weight",
			Group:      "General",
			Editors:    "numeric",
			Options:    "",
			Default:    "",
			Validators: "numeric positive",
		},
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "options",
			Type:       db.ConstTypeJSON,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Options",
			Group:      "Options",
			Editors:    "product_options",
			Options:    "",
			Default:    "",
		},
		{
			Model:      product.ConstModelNameProduct,
			Collection: ConstCollectionNameProduct,
			Attribute:  "related_pids",
			Type:       db.TypeArrayOf(db.ConstTypeInteger),
			IsRequired: false,
			IsStatic:   true,
			Label:      "Related Products",
			Group:      "General",
			Editors:    "product_selector",
			Options:    "",
			Default:    "",
		},
	}

	if product.GetRegisteredStock() != nil {
		result = append(result,
			models.StructAttributeInfo{
				Model:      product.ConstModelNameProduct,
				Collection: ConstCollectionNameProduct,
				Attribute:  "qty",
				Type:       db.ConstTypeInteger,
				IsRequired: true,
				IsStatic:   true,
				Label:      "Qty",
				Group:      "General",
				Editors:    "numeric",
				Options:    "",
				Default:    "0",
				Validators: "numeric positive",
			})
	}

	customAttributesInfo := it.CustomAttributes.GetAttributesInfo()
	for _, customAttribute := range customAttributesInfo {
		result = append(result, customAttribute)
	}

	return result
}
