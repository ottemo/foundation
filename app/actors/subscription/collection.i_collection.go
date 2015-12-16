package subscription

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/subscription"
)

// List enumerates items of Subscription model type in a Subscription collection
func (it *DefaultSubscriptionCollection) List() ([]models.StructListItem, error) {
	var result []models.StructListItem

	dbRecords, err := it.listCollection.Load()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	for _, dbRecordData := range dbRecords {
		subscriptionModel, err := subscription.GetSubscriptionModel()
		if err != nil {
			return result, env.ErrorDispatch(err)
		}
		subscriptionModel.FromHashMap(dbRecordData)

		// retrieving minimal data needed for list
		resultItem := new(models.StructListItem)

		resultItem.ID = subscriptionModel.GetID()
		resultItem.Name = subscriptionModel.GetEmail()
		resultItem.Image = ""
		resultItem.Desc = subscriptionModel.GetName()

		// if extra attributes were required
		if len(it.listExtraAtributes) > 0 {
			resultItem.Extra = make(map[string]interface{})

			for _, attributeName := range it.listExtraAtributes {
				resultItem.Extra[attributeName] = subscriptionModel.Get(attributeName)
			}
		}

		result = append(result, *resultItem)
	}

	return result, nil
}

// ListAddExtraAttribute provides the ability to add additional attributes if the attribute does not already exist
func (it *DefaultSubscriptionCollection) ListAddExtraAttribute(attribute string) error {

	if utils.IsAmongStr(attribute, "_id", "id", "visitor_id", "order_id", "cart_id", "email", "name", "address", "status", "action", "last_submit", "created_at", "created_at") {
		if !utils.IsInListStr(attribute, it.listExtraAtributes) {
			it.listExtraAtributes = append(it.listExtraAtributes, attribute)
		} else {
			return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e7550100-c4dd-4889-a770-2c000c4547c5", "attribute already in list")
		}
	} else {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "bbbd220b-b926-4777-8917-570e5e350e82", "not allowed attribute")
	}

	return nil
}

// ListFilterAdd provides the ability to add a selection filter to List() function
func (it *DefaultSubscriptionCollection) ListFilterAdd(Attribute string, Operator string, Value interface{}) error {
	it.listCollection.AddFilter(Attribute, Operator, Value.(string))
	return nil
}

// ListFilterReset clears the presets made by ListFilterAdd() and ListAddExtraAttribute() functions
func (it *DefaultSubscriptionCollection) ListFilterReset() error {
	it.listCollection.ClearFilters()
	return nil
}

// ListLimit sets the pagination when provided offset and limit values
func (it *DefaultSubscriptionCollection) ListLimit(offset int, limit int) error {
	return it.listCollection.SetLimit(offset, limit)
}
