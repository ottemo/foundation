package subscription

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models/subscription"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Get returns object attribute value or nil for the requested Subscription attribute
func (it *DefaultSubscription) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.id
	case "visitor_id":
		return it.VisitorID
	case "items":
		return it.items
	case "order_id":
		return it.OrderID
	case "customer_email":
		return it.CustomerEmail
	case "customer_name", "full_name":
		return it.CustomerName
	case "shipping_address":
		return it.GetShippingAddress()
	case "billing_address":
		return it.GetBillingAddress()
	case "shipping_method":
		return it.ShippingMethodCode
	case "shipping_rate":
		return it.ShippingRate
	case "payment_instrument":
		return it.PaymentInstrument
	case "status":
		return it.GetStatus()
	case "period":
		return it.GetPeriod()
	case "last_submit":
		return it.LastSubmit
	case "action_date":
		return it.GetActionDate()
	case "created_at":
		return it.CreatedAt
	case "updated_at":
		return it.UpdatedAt
	case "info":
		return it.Info
	}

	return nil
}

// Set will set attribute value of the Subscription to object or return an error
func (it *DefaultSubscription) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		it.id = utils.InterfaceToString(value)

	case "visitor_id":
		it.VisitorID = utils.InterfaceToString(value)

	case "order_id":
		it.OrderID = utils.InterfaceToString(value)

	case "items":
		arrayValue := utils.InterfaceToArray(value)

		it.items = make([]subscription.StructSubscriptionItem, 0)

		for _, arrayItem := range arrayValue {

			if subscriptionItem, ok := arrayItem.(subscription.StructSubscriptionItem); ok {
				it.items = append(it.items, subscriptionItem)
				continue
			}

			mapValue := utils.InterfaceToMap(arrayItem)
			if utils.StrKeysInMap(mapValue, "product_id", "qty", "options") {
				subscriptionItem := subscription.StructSubscriptionItem{
					Name:      utils.InterfaceToString(mapValue["name"]),
					ProductID: utils.InterfaceToString(mapValue["product_id"]),
					Qty:       utils.InterfaceToInt(mapValue["qty"]),
					Options:   utils.InterfaceToMap(mapValue["options"]),
					Sku:       utils.InterfaceToString(mapValue["sku"]),
					Price:     utils.InterfaceToFloat64(mapValue["price"]),
				}

				if subscriptionItem.Qty > 0 && subscriptionItem.ProductID != "" {
					it.items = append(it.items, subscriptionItem)
				}

				continue
			}

			if utils.StrKeysInMap(mapValue, "ProductID", "Qty", "Options") {
				subscriptionItem := subscription.StructSubscriptionItem{
					Name:      utils.InterfaceToString(mapValue["Name"]),
					ProductID: utils.InterfaceToString(mapValue["ProductID"]),
					Qty:       utils.InterfaceToInt(mapValue["Qty"]),
					Options:   utils.InterfaceToMap(mapValue["Options"]),
					Sku:       utils.InterfaceToString(mapValue["Sku"]),
					Price:     utils.InterfaceToFloat64(mapValue["Price"]),
				}

				if subscriptionItem.Qty > 0 && subscriptionItem.ProductID != "" {
					it.items = append(it.items, subscriptionItem)
				}

				continue
			}

			if utils.StrKeysInMap(mapValue, "productid", "qty", "options") {
				subscriptionItem := subscription.StructSubscriptionItem{
					Name:      utils.InterfaceToString(mapValue["name"]),
					ProductID: utils.InterfaceToString(mapValue["productid"]),
					Qty:       utils.InterfaceToInt(mapValue["qty"]),
					Options:   utils.InterfaceToMap(mapValue["options"]),
					Sku:       utils.InterfaceToString(mapValue["sku"]),
					Price:     utils.InterfaceToFloat64(mapValue["price"]),
				}

				if subscriptionItem.Qty > 0 && subscriptionItem.ProductID != "" {
					it.items = append(it.items, subscriptionItem)
				}
			}
		}

	case "customer_email":
		it.CustomerEmail = utils.InterfaceToString(value)

	case "customer_name", "full_name":
		it.CustomerName = utils.InterfaceToString(value)

	case "shipping_address":
		shippingAddress := utils.InterfaceToMap(value)
		if len(shippingAddress) > 0 {
			it.ShippingAddress = shippingAddress
		}

	case "billing_address":
		billingAddress := utils.InterfaceToMap(value)
		if len(billingAddress) > 0 {
			it.BillingAddress = billingAddress
		}

	case "shipping_method":
		shippingMethodCode := utils.InterfaceToString(value)
		if checkout.GetShippingMethodByCode(shippingMethodCode) != nil {
			it.ShippingMethodCode = shippingMethodCode
		}

	case "shipping_rate":
		mapValue := utils.InterfaceToMap(value)
		if utils.StrKeysInMap(mapValue, "Name", "Code", "Price") {
			it.ShippingRate.Name = utils.InterfaceToString(mapValue["Name"])
			it.ShippingRate.Code = utils.InterfaceToString(mapValue["Code"])
			it.ShippingRate.Price = utils.InterfaceToFloat64(mapValue["Price"])
		} else if utils.StrKeysInMap(mapValue, "name", "code", "price") {
			it.ShippingRate.Name = utils.InterfaceToString(mapValue["name"])
			it.ShippingRate.Code = utils.InterfaceToString(mapValue["code"])
			it.ShippingRate.Price = utils.InterfaceToFloat64(mapValue["price"])
		}

	case "payment_instrument":
		it.PaymentInstrument = utils.InterfaceToMap(value)

	case "status":
		it.SetStatus(utils.InterfaceToString(value))

	case "period":
		it.SetPeriod(utils.InterfaceToInt(value))

	case "last_submit":
		it.LastSubmit = utils.InterfaceToTime(value)

	case "action_date":
		it.SetActionDate(utils.InterfaceToTime(value))

	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)

	case "updated_at":
		it.UpdatedAt = utils.InterfaceToTime(value)
	case "info":
		it.Info = utils.InterfaceToMap(value)
	}

	return nil
}

// FromHashMap fills Subscription object attributes from a map[string]interface{}
func (it *DefaultSubscription) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap represents Subscription object as map[string]interface{}
func (it *DefaultSubscription) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.id

	result["visitor_id"] = it.VisitorID
	result["order_id"] = it.OrderID

	result["items"] = it.items

	result["customer_email"] = it.CustomerEmail
	result["customer_name"] = it.CustomerName

	result["status"] = it.Status

	result["period"] = it.Period
	result["shipping_address"] = it.ShippingAddress
	result["billing_address"] = it.BillingAddress

	result["shipping_method"] = it.ShippingMethodCode
	result["shipping_rate"] = it.ShippingRate
	result["payment_instrument"] = it.PaymentInstrument

	result["action_date"] = it.ActionDate
	result["last_submit"] = it.LastSubmit
	result["updated_at"] = it.UpdatedAt
	result["created_at"] = it.CreatedAt

	result["info"] = it.Info

	return result
}

// GetAttributesInfo returns the Subscription attributes information in an array
// TODO: finish list of attributes
func (it *DefaultSubscription) GetAttributesInfo() []models.StructAttributeInfo {
	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
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
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
			Attribute:  "status",
			Type:       db.ConstTypeVarchar,
			IsRequired: true,
			IsStatic:   true,
			Label:      "Status",
			Group:      "General",
			Editors:    "selector",
			Options: strings.Join([]string{
				subscription.ConstSubscriptionStatusSuspended,
				subscription.ConstSubscriptionStatusConfirmed,
				subscription.ConstSubscriptionStatusCanceled,
			}, ","),
			Default: subscription.ConstSubscriptionStatusConfirmed,
		},
		models.StructAttributeInfo{
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
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
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
			Attribute:  "order_id",
			Type:       db.ConstTypeID,
			IsRequired: false,
			IsStatic:   true,
			Label:      "Order",
			Group:      "General",
			Editors:    "model_selector",
			Options:    "model: order",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
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
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
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
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
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
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
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
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
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
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
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
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
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
			Model:      subscription.ConstModelNameSubscription,
			Collection: ConstCollectionNameSubscription,
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
	}

	return info
}
