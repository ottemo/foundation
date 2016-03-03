package order

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Get returns attribute of Order or nil
func (it *DefaultOrder) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id

	case "increment_id":
		return it.IncrementID

	case "status":
		return it.Status

	case "visitor_id":
		return it.VisitorID

	case "session_id":
		return it.SessionID

	case "cart_id":
		return it.CartID

	case "shipping_address":
		return it.ShippingAddress

	case "billing_address":
		return it.BillingAddress

	case "customer_email":
		return it.CustomerEmail

	case "customer_name":
		return it.CustomerName

	case "payment_method":
		return it.PaymentMethod

	case "shipping_method":
		return it.ShippingMethod

	case "subtotal":
		return it.Subtotal

	case "discount":
		return it.Discount

	case "tax_amount":
		return it.TaxAmount

	case "shipping_amount":
		return it.ShippingAmount

	case "grand_total":
		return it.GrandTotal

	case "taxes":
		return it.Taxes

	case "discounts":
		return it.Discounts

	case "created_at":
		return it.CreatedAt

	case "updated_at":
		return it.UpdatedAt

	case "description":
		return it.Description

	case "payment_info":
		return it.PaymentInfo

	case "custom_info":
		return it.CustomInfo

	case "shipping_info":
		return it.ShippingInfo

	case "notes":
		return it.Notes

	}

	return nil
}

// Set sets attribute to Order object, returns error on problems
func (it *DefaultOrder) Set(attribute string, value interface{}) error {

	attribute = strings.ToLower(attribute)

	// attributes can have index like { "note[1]": "my note" }
	attributeIdx := ""
	idx1 := strings.LastIndex(attribute, "[")
	if idx1 >= 0 {
		idx2 := strings.LastIndex(attribute, "]")
		if idx1 < idx2 {
			attributeIdx = attribute[idx1+1 : idx2]
			attribute = attribute[0:idx1]
		}
	}

	switch attribute {
	case "_id", "id":
		it.SetID(utils.InterfaceToString(value))

	case "increment_id":
		it.IncrementID = utils.InterfaceToString(value)

	case "status":
		if "" == it.Status {
			it.Status = utils.InterfaceToString(value)
		} else {
			it.SetStatus(utils.InterfaceToString(value))
		}

	case "visitor_id":
		it.VisitorID = utils.InterfaceToString(value)

	case "session_id":
		it.SessionID = utils.InterfaceToString(value)

	case "cart_id":
		it.CartID = utils.InterfaceToString(value)

	case "customer_email":
		it.CustomerEmail = utils.InterfaceToString(value)

	case "customer_name":
		it.CustomerName = utils.InterfaceToString(value)

	case "billing_address":
		it.BillingAddress = utils.InterfaceToMap(value)

	case "shipping_address":
		it.ShippingAddress = utils.InterfaceToMap(value)

	case "payment_method":
		it.PaymentMethod = utils.InterfaceToString(value)

	case "shipping_method":
		it.ShippingMethod = utils.InterfaceToString(value)

	case "subtotal":
		it.Subtotal = utils.InterfaceToFloat64(value)

	case "discount":
		it.Discount = utils.InterfaceToFloat64(value)

	case "tax_amount":
		it.TaxAmount = utils.InterfaceToFloat64(value)

	case "shipping_amount":
		it.ShippingAmount = utils.InterfaceToFloat64(value)

	case "grand_total":
		it.GrandTotal = utils.InterfaceToFloat64(value)

	case "taxes":
		it.Taxes = make([]order.StructTaxRate, 0)

		arrayValue := utils.InterfaceToArray(value)
		for _, arrayItem := range arrayValue {
			mapValue := utils.InterfaceToMap(arrayItem)
			taxRate := order.StructTaxRate{}

			// Coming from the db
			if utils.StrKeysInMap(mapValue, "name", "code", "amount") {
				taxRate = order.StructTaxRate{
					Name:   utils.InterfaceToString(mapValue["name"]),
					Code:   utils.InterfaceToString(mapValue["code"]),
					Amount: utils.InterfaceToFloat64(mapValue["amount"]),
				}
			}

			// Coming from a struct
			if utils.StrKeysInMap(mapValue, "Name", "Code", "Amount") {
				taxRate = order.StructTaxRate{
					Name:   utils.InterfaceToString(mapValue["Name"]),
					Code:   utils.InterfaceToString(mapValue["Code"]),
					Amount: utils.InterfaceToFloat64(mapValue["Amount"]),
				}
			}

			// if we have data then append
			if taxRate.Name != "" || taxRate.Code != "" || taxRate.Amount != 0 {
				it.Taxes = append(it.Taxes, taxRate)
			}
		}

	case "discounts":
		it.Discounts = make([]order.StructDiscount, 0)

		arrayValue := utils.InterfaceToArray(value)
		for _, arrayItem := range arrayValue {
			mapValue := utils.InterfaceToMap(arrayItem)
			discount := order.StructDiscount{}

			// Coming from the db
			if utils.StrKeysInMap(mapValue, "name", "code", "amount") {
				discount = order.StructDiscount{
					Name:   utils.InterfaceToString(mapValue["name"]),
					Code:   utils.InterfaceToString(mapValue["code"]),
					Amount: utils.InterfaceToFloat64(mapValue["amount"]),
				}
			}

			// Coming from a struct
			if utils.StrKeysInMap(mapValue, "Name", "Code", "Amount") {
				discount = order.StructDiscount{
					Name:   utils.InterfaceToString(mapValue["Name"]),
					Code:   utils.InterfaceToString(mapValue["Code"]),
					Amount: utils.InterfaceToFloat64(mapValue["Amount"]),
				}
			}

			// If we have any data then append
			if discount.Name != "" || discount.Code != "" || discount.Amount != 0 {
				it.Discounts = append(it.Discounts, discount)
			}

		}

	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)

	case "updated_at":
		it.UpdatedAt = utils.InterfaceToTime(value)

	case "description":
		it.Description = utils.InterfaceToString(value)

	case "payment_info":
		it.PaymentInfo = utils.InterfaceToMap(value)

	case "custom_info":
		it.CustomInfo = utils.InterfaceToMap(value)

	case "shipping_info":
		it.ShippingInfo = utils.InterfaceToMap(value)

	case "notes":
		it.Notes = utils.InterfaceToStringArray(value)

	case "note":
		if stringValue := utils.InterfaceToString(value); value != "" {
			if attributeIdx != "" {
				noteIdx := utils.InterfaceToInt(attributeIdx)
				if len(it.Notes) > noteIdx {
					it.Notes[noteIdx] = stringValue
				}
			} else {
				it.Notes = append(it.Notes, stringValue)
			}
		}

	default:
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5a0efabf-9169-42e9-9d78-8e8374998ad6", "unknown attribute: "+attribute)
	}

	return nil
}

// FromHashMap fills Order attributes with values provided in input map
func (it *DefaultOrder) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.LogError(err)
		}
	}

	return nil
}

// ToHashMap makes map from Order attribute values
func (it *DefaultOrder) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["increment_id"] = it.Get("increment_id")
	result["status"] = it.Get("status")

	result["visitor_id"] = it.Get("visitor_id")
	result["session_id"] = it.Get("session_id")
	result["cart_id"] = it.Get("cart_id")

	result["customer_email"] = it.Get("customer_email")
	result["customer_name"] = it.Get("customer_name")

	result["shipping_address"] = it.Get("shipping_address")
	result["billing_address"] = it.Get("billing_address")

	result["payment_method"] = it.Get("payment_method")
	result["shipping_method"] = it.Get("shipping_method")

	result["subtotal"] = it.Get("subtotal")
	result["discount"] = it.Get("discount")
	result["tax_amount"] = it.Get("tax_amount")
	result["shipping_amount"] = it.Get("shipping_amount")
	result["grand_total"] = it.Get("grand_total")

	result["taxes"] = it.Get("taxes")
	result["discounts"] = it.Get("discounts")

	result["created_at"] = it.Get("created_at")
	result["updated_at"] = it.Get("updated_at")

	result["description"] = it.Get("description")
	result["payment_info"] = it.Get("payment_info")
	result["custom_info"] = it.Get("custom_info")
	result["shipping_info"] = it.Get("shipping_info")

	result["notes"] = it.Get("notes")

	return result
}

// GetAttributesInfo describes attributes of Order model
func (it *DefaultOrder) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
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
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "increment_id",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Increment ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "status",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Status",
			Group:      "General",
			Editors:    "selector",
			Options:    "new,pending,processed,declined,complete,cancelled",
			Default:    "new",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "visitor_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Visitor",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: visitor",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "session_id",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Session",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "customer_email",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Customer Email",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
			Validators: "email",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "customer_name",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Customer Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "shipping_address",
			Type:       db.ConstTypeJSON,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Shipping Address",
			Group:      "General",
			Editors:    "visitor_address",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "billing_address",
			Type:       db.ConstTypeJSON,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Customer Name",
			Group:      "General",
			Editors:    "visitor_address",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "payment_method",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Payment Method",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: payments",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "shipping_method",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Shipping Method",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: shipping",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "subtotal",
			Type:       db.ConstTypeDecimal,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Totals",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
			Validators: "numeric positive",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "discount",
			Type:       db.ConstTypeDecimal,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Discount",
			Group:      "Totals",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
			Validators: "numeric positive",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "tax_amount",
			Type:       db.ConstTypeDecimal,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Tax Amount",
			Group:      "Totals",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
			Validators: "numeric positive",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "shipping_amount",
			Type:       db.ConstTypeDecimal,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Shipping Amount",
			Group:      "Totals",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
			Validators: "numeric positive",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "grand_total",
			Type:       db.ConstTypeDecimal,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Grand Total",
			Group:      "Totals",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
			Validators: "numeric positive",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "Taxes",
			Type:       db.ConstTypeJSON,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Taxes",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "Discounts",
			Type:       db.ConstTypeJSON,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Discounts",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "created_at",
			Type:       db.ConstTypeDatetime,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Created At",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "updated_at",
			Type:       db.ConstTypeDatetime,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Updated At",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "description",
			Type:       db.ConstTypeText,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Description",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      order.ConstModelNameOrder,
			Collection: ConstCollectionNameOrder,
			Attribute:  "notes",
			Type:       db.TypeArrayOf(db.ConstTypeText),
			IsRequired: false,
			IsStatic:   true,
			Label:      "Notes",
			Group:      "General",
			Editors:    "string_array",
			Options:    "",
			Default:    "",
		},
	}

	return info
}
