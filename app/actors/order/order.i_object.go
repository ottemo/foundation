package order

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
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

	case "created_at":
		return it.CreatedAt

	case "updated_at":
		return it.UpdatedAt

	case "description":
		return it.Description

	case "payment_info":
		return it.PaymentInfo
	}

	return nil
}

// Set sets attribute to Order object, returns error on problems
func (it *DefaultOrder) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

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

	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)

	case "updated_at":
		it.UpdatedAt = utils.InterfaceToTime(value)

	case "description":
		it.Description = utils.InterfaceToString(value)

	case "payment_info":
		it.PaymentInfo = utils.InterfaceToMap(value)

	default:
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5a0efabf916942e99d788e8374998ad6", "unknown attribute: "+attribute)
	}

	return nil
}

// FromHashMap fills Order attributes with values provided in input map
func (it *DefaultOrder) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			env.ErrorDispatch(err)
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

	result["created_at"] = it.Get("created_at")
	result["updated_at"] = it.Get("updated_at")

	result["description"] = it.Get("description")
	result["payment_info"] = it.Get("payment_info")

	return result
}

// GetAttributesInfo describes attributes of Order model
func (it *DefaultOrder) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
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
			Model:      "Order",
			Collection: "Order",
			Attribute:  "increment_id",
			Type:       "varchar(50)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Increment ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "status",
			Type:       "varchar(50)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Status",
			Group:      "General",
			Editors:    "selector",
			Options:    "pending,canceled,complete",
			Default:    "pending",
		},
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "visitor_id",
			Type:       "id",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Visitor",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: visitor",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "customer_email",
			Type:       "varchar(100)",
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
			Model:      "Order",
			Collection: "Order",
			Attribute:  "customer_name",
			Type:       "varchar(100)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Customer Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "shipping_address",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Shipping Address",
			Group:      "General",
			Editors:    "visitor_address",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "billing_address",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Customer Name",
			Group:      "General",
			Editors:    "visitor_address",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "payment_method",
			Type:       "varchar(100)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Payment Method",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: payments",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "shipping_method",
			Type:       "varchar(100)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Shipping Method",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: shipping",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "subtotal",
			Type:       "decimal(10,2)",
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
			Model:      "Order",
			Collection: "Order",
			Attribute:  "discount",
			Type:       "decimal(10,2)",
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
			Model:      "Order",
			Collection: "Order",
			Attribute:  "tax_amount",
			Type:       "decimal(10,2)",
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
			Model:      "Order",
			Collection: "Order",
			Attribute:  "shipping_amount",
			Type:       "decimal(10,2)",
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
			Model:      "Order",
			Collection: "Order",
			Attribute:  "grand_total",
			Type:       "decimal(10,2)",
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
			Model:      "Order",
			Collection: "Order",
			Attribute:  "created_at",
			Type:       "datetime",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Created At",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "updated_at",
			Type:       "datetime",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Updated At",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      "Order",
			Collection: "Order",
			Attribute:  "description",
			Type:       "text",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Description",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
	}

	return info
}
