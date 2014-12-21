package order

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/order"
)

// Get returns attribute of OrderItem or nil
func (it *DefaultOrderItem) Get(attribute string) interface{} {

	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id

	case "idx":
		return it.idx

	case "order":
		orderInstance, err := order.LoadOrderByID(it.OrderID)
		if err == nil {
			return orderInstance
		}
		return nil

	case "order_id":
		return it.OrderID

	case "product_id":
		return it.ProductID

	case "qty":
		return it.Qty

	case "name":
		return it.Name

	case "sku":
		return it.Sku

	case "short_description":
		return it.ShortDescription

	case "options":
		return it.Options

	case "price":
		return it.Price

	case "weight":
		return it.Weight
	}

	return nil
}

// Set sets attribute to OrderItem object, returns error on problems
func (it *DefaultOrderItem) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)

	case "idx":
		it.idx = utils.InterfaceToInt(value)

	case "order_id":
		it.OrderID = utils.InterfaceToString(value)

	case "product_id":
		it.ProductID = utils.InterfaceToString(value)

	case "qty":
		it.Qty = utils.InterfaceToInt(value)

	case "name":
		it.Name = utils.InterfaceToString(value)

	case "sku":
		it.Sku = utils.InterfaceToString(value)

	case "short_description":
		it.ShortDescription = utils.InterfaceToString(value)

	case "options":
		it.Options = utils.InterfaceToMap(value)

	case "price":
		it.Price = utils.InterfaceToFloat64(value)

	case "weight":
		it.Weight = utils.InterfaceToFloat64(value)

	default:
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bb97327a389f4a7a91c28ad4e82c6be5", "unknown attribute: "+attribute)
	}

	return nil
}

// FromHashMap fills OrderItem attributes with values provided in input map
func (it *DefaultOrderItem) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap makes map from OrderItem attribute values
func (it *DefaultOrderItem) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.Get("_id")
	result["idx"] = it.Get("idx")

	result["order_id"] = it.Get("order_id")
	result["product_id"] = it.Get("product_id")

	result["qty"] = it.Get("qty")

	result["name"] = it.Get("name")
	result["sku"] = it.Get("sku")
	result["short_description"] = it.Get("short_description")

	result["options"] = it.Get("options")

	result["price"] = it.Get("price")
	result["weight"] = it.Get("weight")

	return result
}

// GetAttributesInfo describes attributes of OrderItem model
func (it *DefaultOrderItem) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "_id",
			Type:       "id",
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "idx",
			Type:       "int",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Increment ID",
			Group:      "General",
			Editors:    "integer",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "order_id",
			Type:       "id",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Order",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "product_id",
			Type:       "id",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Visitor",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "qty",
			Type:       "int",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Qty",
			Group:      "General",
			Editors:    "integer",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "name",
			Type:       "varchar(150)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Name",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "sku",
			Type:       "varchar(100)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Sku",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "short_description",
			Type:       "varchar(255)",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Short Description",
			Group:      "General",
			Editors:    "text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "options",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Options",
			Group:      "General",
			Editors:    "product_options",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "price",
			Type:       "decimal(10,2)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Price",
			Group:      "Prices",
			Editors:    "numeric",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "OrderItem",
			Collection: "OrderItem",
			Attribute:  "weight",
			Type:       "decimal(10,2)",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Weight",
			Group:      "Sizes",
			Editors:    "numeric",
			Options:    "",
			Default:    "",
		},
	}

	return info
}
