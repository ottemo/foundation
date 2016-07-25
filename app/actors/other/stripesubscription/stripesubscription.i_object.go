package stripesubscription

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/stripesubscription"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Get returns object attribute value or nil for the requested Stripe Subscription attribute
func (it *DefaultStripeSubscription) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "visitor_id":
		return it.VisitorID
	case "customer_name":
		return it.CustomerName
	case "customer_email":
		return it.CustomerEmail
	case "billing_address":
		return it.BillingAddress
	case "shipping_address":
		return it.ShippingAddress
	case "stripe_subscription_id":
		return it.StripeSubscriptionID
	case "stripe_customer_id":
		return it.StripeCustomerID
	case "stripe_coupon":
		return it.StripeCoupon
	case "stripe_events":
		return it.StripeEvents
	case "price":
		return it.Price
	case "created_at":
		return it.CreatedAt
	case "updated_at":
		return it.UpdatedAt
	case "description":
		return it.Description
	case "info":
		return it.Info
	case "status":
		return it.Status
	}

	return nil
}

// Set will set attribute value of the Stripe Subscription to object or return an error
func (it *DefaultStripeSubscription) Set(attribute string, value interface{}) error {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)
	case "visitor_id":
		it.VisitorID = utils.InterfaceToString(value)
	case "customer_name", "full_name":
		it.CustomerName = utils.InterfaceToString(value)
	case "customer_email":
		it.CustomerEmail = utils.InterfaceToString(value)
	case "billing_address":
		it.BillingAddress = utils.InterfaceToMap(value)
	case "shipping_address":
		it.ShippingAddress = utils.InterfaceToMap(value)
	case "stripe_subscription_id":
		it.StripeSubscriptionID = utils.InterfaceToString(value)
	case "stripe_customer_id":
		it.StripeCustomerID = utils.InterfaceToString(value)
	case "stripe_coupon":
		it.StripeCoupon = utils.InterfaceToString(value)
	case "stripe_events":
		it.StripeEvents = utils.InterfaceToMap(value)
	case "price":
		it.Price = utils.InterfaceToFloat64(value)
	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)
	case "updated_at":
		it.UpdatedAt = utils.InterfaceToTime(value)
	case "description":
		it.Description = utils.InterfaceToString(value)
	case "info":
		it.Info = utils.InterfaceToMap(value)
	case "status":
		it.Status = utils.InterfaceToString(value)
	}

	return nil
}

// FromHashMap fills Stripe Subscription object attributes from a map[string]interface{}
func (it *DefaultStripeSubscription) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap represents Stripe Subscription object as map[string]interface{}
func (it *DefaultStripeSubscription) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["visitor_id"] = it.VisitorID
	result["customer_name"] = it.CustomerName
	result["customer_email"] = it.CustomerEmail
	result["billing_address"] = it.BillingAddress
	result["shipping_address"] = it.ShippingAddress

	result["stripe_subscription_id"] = it.StripeSubscriptionID
	result["stripe_customer_id"] = it.StripeCustomerID
	result["stripe_coupon"] = it.StripeCoupon
	result["stripe_events"] = it.StripeEvents

	result["price"] = it.Price

	result["created_at"] = it.CreatedAt
	result["updated_at"] = it.UpdatedAt

	result["description"] = it.Description
	result["info"] = it.Info
	result["status"] = it.Status

	return result
}

// GetAttributesInfo returns the Stripe Subscription attributes information in an array
func (it *DefaultStripeSubscription) GetAttributesInfo() []models.StructAttributeInfo {
	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
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
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
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
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "customer_name",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Customer Name",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "customer_email",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "billing_address",
			Type:       db.ConstTypeJSON,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Billing Address",
			Group:      "General",
			Editors:    "visitor_address",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "shipping_address",
			Type:       db.ConstTypeJSON,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Shipping Address",
			Group:      "General",
			Editors:    "visitor_address",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "stripe_subscription_id",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Stripe Subscription ID",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "stripe_customer_id",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Stripe Customer ID",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "stripe_coupon",
			Type:       db.ConstTypeVarchar,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Stripe Coupon",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "stripe_events",
			Type:       db.ConstTypeJSON,
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "Stripe Events",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "price",
			Type:       db.ConstTypeDecimal,
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
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
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
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
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "description",
			Type:       db.ConstTypeText,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Description",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      stripesubscription.ConstModelNameStripeSubscription,
			Collection: ConstCollectionNameStripeSubscription,
			Attribute:  "status",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Status",
			Group:      "General",
			Editors:    "selector",
			Options:    "",
			Default:    "",
		},
	}

	return info
}
